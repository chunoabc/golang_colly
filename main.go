package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url, ok := os.LookupEnv("URL")
	if !ok {
		url = "https://www.deviantart.com/wlop"
	}
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) { // 當Visit訪問網頁後，網頁響應(Response)時候執行的事情
		// fmt.Println(string(r.Body)) // 返回的Response物件r.Body 是[]Byte格式，要再轉成字串
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println(err)
	})

	// 抓類別Class 名稱
	c.OnHTML(".f_4No > section > a > div > img", func(e *colly.HTMLElement) {
		// fmt.Println("------------")
		// fmt.Println(e.Index)
		// fmt.Println(e)
		imageHref := e.Attr("srcset")
		url := strings.Split(imageHref, " ")[0]
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		// fmt.Println(fmt.Sprintf("%d : %s", e.Index, resp.Header["Content-Type"][0]))
		// fmt.Println(imageHref)
		// fmt.Println("")

		var ext string
		switch resp.Header["Content-Type"][0] {
		case "image/jpeg":
			ext = "jpg"
		// case "image/svg+xml":
		// 	ext = "svg"
		default:
			ext = ""
		}
		if ext != "" {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			// fmt.Println(e.Index)
			name := fmt.Sprintf("%d.%s", e.Index, ext)
			os.WriteFile(name, body, 0666) //存下圖檔
		}
	})

	c.Visit(url)
}
