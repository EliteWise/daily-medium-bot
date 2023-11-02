package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
)

type Secrets struct {
	DiscordKey string `json:"discordKey"`
}

func main() {
	file, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Désérialiser le contenu du fichier dans une structure
	var secrets Secrets
	err = json.Unmarshal(file, &secrets)
	if err != nil {
		log.Fatalf("Failed to unmarshal: %s", err)
	}

	sess, err := discordgo.New("Bot " + secrets.DiscordKey)
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()
	var href = ""

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		if e.Text == "Recommended stories" {
			e.DOM.ParentsUntil("~ a").Each(func(index int, sel *goquery.Selection) {
				if href == "" {
					sel.Find("a").Each(func(index int, aSel *goquery.Selection) {
						if index == 0 {
							relativeURL, _ := aSel.Attr("href")
							long_href := e.Request.AbsoluteURL(relativeURL)
							href = strings.Split(long_href, "?source")[0]
						}
					})
				}

			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err1 := c.Visit("https://medium.com/tag/programming")
	if err1 != nil {
		log.Fatal(err1)
	}

	log.Println("Continue")

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!daily" {
			s.ChannelMessageSend(m.ChannelID, href)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("Online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
