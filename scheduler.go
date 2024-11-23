package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func sendArticle() {
	// Create new scheduler
	c := cron.New() // To use seconds, add: `cron.WithSeconds()` inside New() func

	fmt.Printf("Scheduler is starting...")

	var setupData SetupData
	deserializeData(CONFIG_SOURCE, &setupData)
	timeConverter, err := convertTimeToCron(setupData.HourToSend)
	if err != nil {
		fmt.Println(err)
	}
	cronExpression := fmt.Sprintf("0 %s * * *", timeConverter)

	emptyID := ""

	// Add the task to be executed every day
	c.AddFunc(cronExpression, func() { // To test every x seconds, use: "*/15 * * * * *"
		var setupData SetupData
		deserializeData(CONFIG_SOURCE, &setupData)
		if setupData.Mode == "channel_mode" {
			sendToChannel(setupData.SelectedChannelID, searchArticle(emptyID, &setupData.MediumCategory))
		} else {
			sendDM(setupData.UserID, searchArticle(emptyID, &setupData.MediumCategory))
		}
	})

	c.Start()
	// Prevents the func from stopping
	select {}
}
