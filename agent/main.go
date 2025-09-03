package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os/exec"
    "time"

    "github.com/chromedp/chromedp"
)

// Local model settings.  Adjust lmBaseURL and model for your setup.
const lmBaseURL = "http://localhost:1234/v1/chat/completions"
const model = "qwen2.5-coder:7b"

// runRandomOrgOnce drives a Chrome session to random.org/integers, clicks
// the "Get Numbers" button, and returns the resulting numbers or a status
// message.  The browser remains open for user inspection after the call.
func runRandomOrgOnce(ctx context.Context) (string, error) {
    cmd := exec.Command(
        "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
        "--remote-debugging-port=9222",
        "--no-first-run", "--no-default-browser-check",
        "--user-data-dir=/tmp/chromedp-profile",
    )
    _ = cmd.Start()

    wsURL, err := waitDevTools("http://localhost:9222/json/version", 10*time.Second)
    if err != nil {
        return "", fmt.Errorf("failed to connect to devtools: %w", err)
    }

    allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(ctx, wsURL)
    defer cancelAlloc()
    tabCtx, cancelTab := chromedp.NewContext(allocCtx)
    defer cancelTab()

    var output string
    err = chromedp.Run(tabCtx,
        chromedp.Navigate("https://www.random.org/integers/"),
        chromedp.WaitVisible(`input[value="Get Numbers"]`),
        chromedp.Click(`input[value="Get Numbers"]`),
        chromedp.Sleep(1500*time.Millisecond),
        chromedp.Text(`#invisible > pre`, &output, chromedp.NodeVisible, chromedp.ByQuery),
    )
    if err != nil {
        return "", err
    }
    if output == "" {
        output = "Clicked Get Numbers; please view the browser for results."
    }
    return output, nil
}

// waitDevTools repeatedly queries the Chrome DevTools version endpoint until
// it returns a WebSocketDebuggerURL or the timeout expires.
func waitDevTools(url string, timeout time.Duration) (string, error) {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        resp, err := http.Get(url)
        if err == nil && resp.StatusCode == http.StatusOK {
            defer resp.Body.Close()
            var data struct {
                WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
            }
            if json.NewDecoder(resp.Body).Decode(&data) == nil && data.WebSocketDebuggerURL != "" {
                return data.WebSocketDebuggerURL, nil
            }
        }
        time.Sleep(300 * time.Millisecond)
    }
    return "", fmt.Errorf("devtools endpoint timeout")
}

// toolDef defines an OpenAI-compatible tool with a function signature.
type toolDef struct {
    Type     string `json:"type"`
    Function struct {
        Name        string                 `json:"name"`
        Description string                 `json:"description"`
        Parameters  map[string]interface{} `json:"parameters"`
    } `json:"function"`
}

// chatResp captures the subset of the OpenAI chat completion response used
// to check for tool calls and return messages.
type chatResp struct {
    Choices []struct {
        Message struct {
            Content   string `json:"content"`
            ToolCalls []struct {
                ID       string `json:"id"`
                Type     string `json:"type"`
                Function struct {
                    Name      string `json:"name"`
                    Arguments string `json:"arguments"`
                } `json:"function"`
            } `json:"tool_calls"`
        } `json:"message"`
    } `json:"choices"`
}

// Text returns the primary message content.
func (r chatResp) Text() string {
    if len(r.Choices) == 0 {
        return ""
    }
    return r.Choices[0].Message.Content
}

// ToolCalls returns the tool calls from the first choice.
func (r chatResp) ToolCalls() []struct {
    ID       string `json:"id"`
    Type     string `json:"type"`
    Function struct {
        Name      string `json:"name"`
        Arguments string `json:"arguments"`
    } `json:"function"`
} {
    if len(r.Choices) == 0 {
        return nil
    }
    return r.Choices[0].Message.ToolCalls
}

// callLM sends a chat completion request to the local LLM and decodes the response.
func callLM(messages []map[string]string, tools []toolDef) chatResp {
    body := map[string]interface{}{
        "model":    model,
        "messages": messages,
        "tools":    tools,
    }
    payload, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", lmBaseURL, bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatalf("LM call failed: %v", err)
    }
    defer res.Body.Close()
    respBody, _ := io.ReadAll(res.Body)
    var out chatResp
    if err := json.Unmarshal(respBody, &out); err != nil {
        log.Fatalf("failed to decode LM response: %v", err)
    }
    return out
}

func main() {
    // Advertise the get_random_number tool to the model.
    tools := []toolDef{
        {
            Type: "function",
            Function: struct {
                Name        string                 `json:"name"`
                Description string                 `json:"description"`
                Parameters  map[string]interface{} `json:"parameters"`
            }{
                Name:        "get_random_number",
                Description: "Open random.org and click Get Numbers to return a random number.",
                Parameters:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
            },
        },
    }

    // Initial conversation.
    messages := []map[string]string{
        {"role": "system", "content": "You may call tools to perform external actions."},
        {"role": "user", "content": "Get me a random number."},
    }

    // Ask the model to decide whether to call the tool.
    resp := callLM(messages, tools)
    for _, call := range resp.ToolCalls() {
        if call.Function.Name == "get_random_number" {
            ctx := context.Background()
            out, err := runRandomOrgOnce(ctx)
            if err != nil {
                out = "ERROR: " + err.Error()
            }
            messages = append(messages, map[string]string{
                "role":        "tool",
                "content":     out,
                "tool_call_id": call.ID,
            })
        }
    }

    // Ask the model for a final answer after tool execution.
    final := callLM(messages, tools)
    fmt.Println("LLM response:", final.Text())
}
