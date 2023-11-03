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

// Create a collection to store the discord key
type Secrets struct {
	DiscordKey string `json:"discordKey"`
}

func main() {
	file, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var secrets Secrets
	// Deserialization of the json file by passing a pointer to the secrets variable, to assign the result
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

	/* Register a callback function with colly's OnHTML method to be executed every time an HTML element matching the "h2" CSS selector is encountered during the scraping process.
	* The callback function takes a pointer to a colly.HTMLElement as its argument, which represents the actual <h2> element found
	 */
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

	// Registers a callback function to handle incoming Discord messages, where `s` represents the current session and `m` represents the created message event.
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

	// Ensures that function call is executed just before the main exits, used for cleanup task.
	defer sess.Close()

	fmt.Println("Online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
