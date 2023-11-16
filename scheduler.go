package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func sendArticle() {
	// Create new scheduler
	c := cron.New()

	var setupData SetupData
	sendingTime := deserializeData(CONFIG_SOURCE, setupData.HourToSend)
	cronExpression := fmt.Sprintf("0 %d * * *", sendingTime)

	// Add the task to be executed every day
	c.AddFunc(cronExpression, func() {
		fmt.Println("Send: ", time.Now())
	})

	c.Start()
	// Prevents the func from stopping
	select {}
}
