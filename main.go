package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	BASE_URL = `https://www.edmunds.com/used-cars-for-sale/`
)

func main() {
	req, err := http.NewRequest("GET", BASE_URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	SetHeaders(req)
	client := http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.Header.Get("content-encoding") == "gzip" {
		str, err := DecompressGzip(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		vehicles, err := ScrapeVehicleTypes(str)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range vehicles {
			fmt.Println(v)
		}
	}
}

func ScrapeVehicleTypes(htmlStr string) ([]string, error) {
	htmlNode, err := html.Parse(strings.NewReader(htmlStr))	
	if err != nil {
		return nil, err
	}
	var vehicles []string
	var scraper func(node *html.Node)
	scraper = func(node *html.Node) {
		underline := false
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "underline" {
					underline = true
				}
				if attr.Key == "href" && strings.HasPrefix(attr.Val, "/used-") && underline {
					if node.FirstChild != nil {
						vehicles = append(vehicles, attr.Val[6:len(attr.Val)-1])
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			scraper(c)
		}
	}
	scraper(htmlNode)
	return vehicles, nil
}

func DecompressGzip(body io.Reader) (string, error) {
	reader, err := gzip.NewReader(body)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SetHeaders(req *http.Request) {
	req.Header.Set("user-agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36`)
	req.Header.Set("accept-encoding", "gzip")
}

/*
func DecompressBrotli(data []byte) (string, error) {

}
*/
