package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/beeep"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	url := "https://www.ticketswap.com/event/rotterdam-rave-indoor-closing-2019/d43593ad-8870-4f81-9b86-92872e8de1a0"
	for {
		getAndNotify(url)
		time.Sleep(3 * time.Second)
	}
}

func getAndNotify(url string) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	} else {
		if res.StatusCode == 200 {
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Print("There was an error")
				log.Fatal(err)
			} else {
				doc.Find("span").Each(func(i int, selection *goquery.Selection) {
					if selection.Text() == "Available" {
						availableCountText, err := selection.Parent().Parent().Children().Find("h2").Html()
						availableCount, err := strconv.ParseInt(availableCountText, 10, 32)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("Available count %d", availableCount)
						if availableCount > 0 {
							err := beeep.Notify("Hurry Up!", "Ticket is on sale!", "icon.png")
							if err != nil {
								panic(err)
							}
						}
					}
				})
			}
		}
	}
}
