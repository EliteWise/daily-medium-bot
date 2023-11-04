package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
