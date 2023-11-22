package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func sendArticle() {
	// Create new scheduler
	c := cron.New(cron.WithSeconds()) // To use seconds, add: `cron.WithSeconds()` inside New() func

	fmt.Printf("Scheduler is starting...")

	/*var setupData SetupData
	sendingTime := deserializeData(CONFIG_SOURCE, setupData.HourToSend)
	cronExpression := fmt.Sprintf("0 %d * * *", sendingTime)*/

	emptyID := ""

	// Add the task to be executed every day
	c.AddFunc("*/15 * * * * *", func() { // To test every x seconds, use: "*/15 * * * * *"
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
