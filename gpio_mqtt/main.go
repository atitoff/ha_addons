package main

import (
	"fmt"
	"gpio_mqtt/gpio_handler"
	"os"
	"strconv"
	"time"
)

type Config = gpio_handler.Config

var config Config

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadConfig() {
	var err error
	config.MqttPort, err = strconv.Atoi(os.Getenv("MqttPort"))
	if err != nil {
		config.MqttPort = 1883
	}
	config.MqttHost = os.Getenv("MqttHost")
	if config.MqttHost == "" {
		config.MqttHost = "192.168.1.117"
	}
	config.MqttClientId = os.Getenv("MqttClientId")
	if config.MqttClientId == "" {
		config.MqttClientId = "mqtt_gpio"
	}
	config.MqttUsername = os.Getenv("MqttUsername")
	if config.MqttUsername == "" {
		config.MqttUsername = "mqtt"
	}
	config.MqttPassword = os.Getenv("MqttPassword")
	if config.MqttPassword == "" {
		config.MqttPassword = "mqtt"
	}
	config.LogLevel = os.Getenv("LogLevel")
	if config.LogLevel == "" {
		config.LogLevel = "warning"
	}
	fmt.Printf(
		"MqttPort: %d\nMqttHost: %s\nMqttClientId: %s\nMqttUsername: %s\nMqttPassword: %s\nLogLevel: %s\n",
		config.MqttPort, config.MqttHost, config.MqttClientId, config.MqttUsername, config.MqttPassword, config.LogLevel,
	)
}

func main() {
	argsWithProg := os.Args
	fmt.Println(argsWithProg)

	loadConfig()
	go gpio_handler.Run(config)
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("tick")
	}
}
