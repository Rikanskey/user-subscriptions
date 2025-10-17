package main

import "user-subscriptions/internal/runner"

const (
	configDir    = "./config"
	migrationDir = "migration/db"
)

func main() {
	runner.Start(configDir, migrationDir)
}
