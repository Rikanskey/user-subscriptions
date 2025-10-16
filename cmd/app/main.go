package main

import "user-subscriptions/internal/runner"

const configDir = "./config"

func main() {
	runner.Start(configDir)
}
