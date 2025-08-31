package main

import (
    "context"
    "fmt"
    "log"
    "time"

    // chromedp provides a convenient wrapper around the Chrome DevTools Protocol.
    // It allows us to programmatically drive a real browser session. See
    // https://github.com/chromedp/chromedp for details.
    "github.com/chromedp/chromedp"
)

// main launches a Chrome browser, navigates to Random.org's integer generator
// and clicks the "Get Numbers" button. This serves as a simple demonstration
// of agentic automation in Go. The program relies on the chromedp package,
// which in turn requires a Chrome or Chromium installation on your system.
//
// To run this example:
//   1. Ensure Go is installed on your machine (see README for instructions).
//   2. Initialise a module (`go mod init example.com/randomagent`) and add
//      chromedp as a dependency (`go get github.com/chromedp/chromedp`).
//   3. Build or run the program (`go run main.go`).
// You should see your browser open, fetch random.org and click the button.
func main() {
    // Create a new context for chromedp. The parent context controls the
    // lifetime of the browser process. Using WithTimeout prevents the
    // automation from running indefinitely.
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    // Set a hard timeout for the automation. If the actions do not complete
    // within 60 seconds, the context will be cancelled and the program exits.
    ctx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second)
    defer timeoutCancel()

    var result string
    // Execute a series of browser actions defined as chromedp.Tasks.
    err := chromedp.Run(ctx,
        // Navigate to the Random.org integer generator page.
        chromedp.Navigate("https://www.random.org/integers/"),

        // Wait until the "Get Numbers" submit button is visible. This ensures
        // that the page has finished loading before we proceed.
        chromedp.WaitVisible(`input[value="Get Numbers"]`, chromedp.ByQuery),

        // Click the "Get Numbers" button to generate random integers. We use
        // ByQuery because the selector is a CSS query.
        chromedp.Click(`input[value="Get Numbers"]`, chromedp.ByQuery),

        // Give the page a few seconds to update. Without this sleep, the
        // program may exit before the results are rendered.
        chromedp.Sleep(5*time.Second),

        // Capture the inner HTML of the <body> element. This isn't strictly
        // necessary but can be useful for debugging or extracting the results.
        chromedp.InnerHTML("body", &result, chromedp.ByQuery),
    )
    if err != nil {
        log.Fatalf("automation failed: %v", err)
    }

    fmt.Println("Random.org automation completed successfully.")
    // Uncomment the following line to print the page HTML (handy for debugging).
    // fmt.Println(result)
}
