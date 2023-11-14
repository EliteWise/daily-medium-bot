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
		if hour < 12 {
			dayHours = append(dayHours, (strconv.Itoa(hour) + "am"))
		} else {
			dayHours = append(dayHours, (strconv.Itoa(hour) + "pm"))
		}
	}
	return dayHours
}
