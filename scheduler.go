package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func sendArticle() {
	// Create new scheduler
	c := cron.New(cron.WithSeconds())

	fmt.Printf("Scheduler is starting...")

	//var setupData SetupData
	//sendingTime := deserializeData(CONFIG_SOURCE, setupData.HourToSend)
	//cronExpression := fmt.Sprintf("0 %d * * *", sendingTime)

	emptyID := ""

	// Add the task to be executed every day
	c.AddFunc("*/15 * * * * *", func() {
		var setupData SetupData
		deserializeData(CONFIG_SOURCE, &setupData)
		if setupData.Mode == "channel_mode" {
			sendToChannel(setupData.SelectedChannelID, searchArticle(emptyID))
		} else {
			sendDM(setupData.UserID, searchArticle(emptyID))
		}
	})

	c.Start()
	// Prevents the func from stopping
	select {}
}
