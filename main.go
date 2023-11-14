package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
)

// Create a collection to store the discord key
type Secrets struct {
	DiscordKey string `json:"discordKey"`
}

type SetupData struct {
	SelectedChannel string `json:"selectedChannel"`
	HourToSend      string `json:"hourToSend"`
	PreviousArticle string `json:"previousArticle"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Button struct {
	Label    string `json:"label"`
	Style    string `json:"style"`
	CustomID string `json:"customID"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"Description"`
	FieldLeft   Field  `json:"field_left"`
	FieldRight  Field  `json:"field_right"`
	ButtonLeft  Button `json:"button_left"`
	ButtonRight Button `json:"button_right"`
}

func main() {
	var secrets Secrets
	deserializeData("secrets.json", &secrets)

	sess, err := discordgo.New("Bot " + secrets.DiscordKey)
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()
	var href = ""

	/* Try to get the first span that contains "hours" keyword, and then get the link above
	 */
	c.OnHTML("span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "hours") {
			a := e.DOM.ParentsUntil("body").Filter("a").First()
			if href_, exists := a.Attr("href"); exists {
				long_href := e.Request.AbsoluteURL(href_)
				href = strings.Split(long_href, "?source")[0]
				var setupData SetupData
				setupData.PreviousArticle = href
				serializeData("setup-data.json", setupData)
			}
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
		} else if m.Content == "!daily setup" {

		}
	})

	// Define the new slash command.
	var (
		dailyCommand = &discordgo.ApplicationCommand{
			Name:        "daily",
			Description: "Responds with a daily article.",
		}
		setupCommand = &discordgo.ApplicationCommand{
			Name:        "setup",
			Description: "Sets up your daily preferences.",
		}
	)

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Online!")

		// Register the command for the guild (server).
		_, errCommand := sess.ApplicationCommandCreate(sess.State.Application.ID, "", dailyCommand)
		if errCommand != nil {
			log.Fatalf("Cannot create slash command: %v", errCommand)
		}

		_, errSetupCommand := sess.ApplicationCommandCreate(sess.State.Application.ID, "", setupCommand)
		if errSetupCommand != nil {
			log.Fatalf("Cannot create slash command: %v", errSetupCommand)
		}
	})

	// Interactions Management like clicks on buttons
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {

			// Use i.MessageComponentData().CustomID to identify the clicked button

			switch i.MessageComponentData().CustomID {
			case "private_message_mode":
				s.ChannelMessageSend(i.ChannelID, "Private message mode selected.")
			case "channel_mode":
				s.ChannelMessageSend(i.ChannelID, "Channel mode selected.")
				mediumCatList := retrieveMediumCategories().MediumCategories
				updateEmbed(s, i, "Menu Item", "Desc Menu Item", mediumCatList)
			case "custom_menu":
				log.Println(i.MessageComponentData().Values[0])
			}

			// Acknowledge the interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
			})
			if err != nil {
				log.Println("Erreur lors de l'envoi de l'acknowledgement:", err)
				return
			}
		}
	})

	// Add the handler for the newly created command.
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if i.ApplicationCommandData().Name == "daily" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: href,
					},
				})
			} else if i.ApplicationCommandData().Name == "setup" {
				components := make([]discordgo.MessageComponent, 0)

				components = append(components, discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "Private Message Mode",
							Style:    discordgo.PrimaryButton,
							CustomID: "private_message_mode",
						},
						&discordgo.Button{
							Label:    "Channel Mode",
							Style:    discordgo.SecondaryButton,
							CustomID: "channel_mode",
						},
					},
				})

				embed := &discordgo.MessageEmbed{
					Title:       "Medium Daily Configuration",
					Description: "Choose your options for your daily articles.",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Private Message Mode",
							Value:  "Send articles inside your DM",
							Inline: true,
						},
						{
							Name:   "Channel Mode",
							Value:  "Send articles inside a custom channel",
							Inline: true,
						},
					},
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds:     []*discordgo.MessageEmbed{embed},
						Components: components,
					},
				})
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Ensures that function call is executed just before the main exits, used for cleanup task.
	defer sess.Close()

	waitForInterrupt()
}
