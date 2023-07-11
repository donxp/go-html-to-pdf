package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	// Register HTTP handler
	http.HandleFunc("/generate", handlePostGenerate)

	// Start HTTP server
	fmt.Println("Starting HTTP server")
	http.ListenAndServe(fmt.Sprintf(":%d", mustGetHttpServerPort()), nil)
}

func htmlToPdf(html string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

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
		panic(err)
	}

	return buf
}

func handlePostGenerate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		bearerToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

		if mustGetApiToken() != bearerToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}
		
		htmlBody := string(body[:])
		generatedPdf := htmlToPdf(htmlBody)
		w.Write(generatedPdf)
	default:
		http.Error(w, "404 not found", http.StatusNotFound)
	}
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

func mustGetHttpServerPort() int {
	var serverPort int64 = 3000
	port, success := os.LookupEnv("HTML_TO_PDF_PORT")
	if success {
		serverPort, _ = strconv.ParseInt(port, 10, 32)
	}

	return int(serverPort)
}

func mustGetApiToken() string {
	token := os.Getenv("API_TOKEN")
	return token
}