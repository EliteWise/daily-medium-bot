package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func serializeData(json_file string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = os.WriteFile(json_file, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write: %s", err)
	}
	return nil
}

func deserializeData(json_file string, value interface{}) error {
	file, err := os.ReadFile(json_file)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}
	// Deserialization of the json file by passing a pointer to the secrets variable, to assign the result
	err = json.Unmarshal(file, &value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %s", err)
	}
	return nil
}

func waitForInterrupt() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func updateEmbed(session *discordgo.Session, i *discordgo.InteractionCreate, title string, description string, options []string) {
	currentEmbed := i.Message.Embeds[0]
	embed := &discordgo.MessageEmbed{
		Title:       currentEmbed.Title,
		Description: currentEmbed.Description,
	}

	selectMenuOptions := make([]discordgo.SelectMenuOption, len(options))
	for i, optionString := range options {
		selectMenuOptions[i] = discordgo.SelectMenuOption{
			Label: optionString,
			Value: optionString,
		}
	}

	selectMenu := &discordgo.SelectMenu{
		CustomID:    "custom_menu",
		Placeholder: "Select an option",
		Options:     selectMenuOptions,
	}

	row := &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{selectMenu},
	}

	err := session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{row},
		},
	})

	if err != nil {
		log.Printf("Failed to update Embed: %v", err)
	}

}
