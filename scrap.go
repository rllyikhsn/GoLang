package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Tracking struct {
	Status struct {
		Code     string `json:"code"`
		Messages string `json:"messages"`
	} `json:"status"`
	Data struct {
		ReceivedBy string `json:"receivedBy"`
		// HistoriesLog
		History []HistoriesLog `json:"histories"`
	} `json:"data"`
}

type HistoriesLog struct {
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	Formatted   struct {
		CreatedAt string `json:"createdAt"`
	} `json:"formatted"`
}

func main() {
	res, err := http.Get("https://gist.githubusercontent.com/nubors/eecf5b8dc838d4e6cc9de9f7b5db236f/raw/d34e1823906d3ab36ccc2e687fcafedf3eacfac9/jne-awb.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// if res.StatusCode == 200 {
	// 	bodyText, err := ioutil.ReadAll(res.Body)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Printf("%s\n", bodyText)
	// }

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tc Tracking
	rows := make([]HistoriesLog, 0)
	row := new(HistoriesLog)
	doc.Find(".main-content table").Each(func(a int, x *goquery.Selection) {
		doc.Find("tbody tr").Children().Each(func(c1 int, sc1 *goquery.Selection) {
			if c1 >= 9 {
				if a == 0 {
					if c1%2 == 1 {
						layoutFormat := "01-02-2006 15:04"
						date, _ := time.Parse(layoutFormat, sc1.Text())
						dy := date.Format("2006-01-02T15:04+07:00")
						dx := date.Format("01 January 2006, 15:04 WIB")
						row.CreatedAt = dy
						row.Formatted.CreatedAt = dx
					} else if c1%2 == 0 {
						row.Description = string(sc1.Text())
						rows = append(rows, *row)
					}
				}
			}
		})

		doc.Find("tbody tr").Children().Each(func(c1 int, sc1 *goquery.Selection) {
			if c1 == 0 {
				if a == 0 {
					tc.Status.Code = sc1.Text()
				}
			}
		})
	})

	tc.Status.Messages = "Delivery tracking detail fetched successfully"
	tc.Data.History = rows
	b, err := json.Marshal(tc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))

}
