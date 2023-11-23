package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func retrieveChannels(session *discordgo.Session, i *discordgo.InteractionCreate) []string {
	channels, _ := session.GuildChannels(i.GuildID)
	channelsNames := []string{}

	for _, channel := range channels {
		channelsNames = append(channelsNames, channel.Name)
	}

	return channelsNames
}

type Categories struct {
	MediumCategories []string `json:"mc"`
}

func retrieveMediumCategories() Categories {

	var categories Categories

	deserializeData("medium-categories.json", &categories)

	return categories
}

func retrieveDayHours() []string {
	var dayHours []string
	for hour := 0; hour < 24; hour++ {
		hourDisplay := hour

		if hour == 0 {
			// Midnight
			hourDisplay = 12
		} else if hour > 12 {
			// Afternoon
			hourDisplay = hour - 12
		}

		suffix := "am"
		if hour >= 12 {
			suffix = "pm"
		}

		dayHours = append(dayHours, strconv.Itoa(hourDisplay)+suffix)
	}
	return dayHours
}
