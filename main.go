package main

import (
	"context"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	// Start chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Register HTTP handler

	// Start HTTP server

	html := "<html><body><h1>Hello World</h1></body></html>"

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		printToPDF(&buf)); err != nil {
		log.Fatal(err)
	}

	os.WriteFile("sample.pdf", buf, 0o644)
}

func printToPDF(res *[]byte) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)

		if err != nil {
			return err
		}
		*res = buf
		return nil
	})
}