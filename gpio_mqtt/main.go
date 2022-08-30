package main

import (
	"fmt"
	"gpio_mqtt/gpio_handler"
	"math/rand"
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
	// [--mqtthost=core-mosquitto  	1
	// --MqttClientId=gpio_mqtt		2
	// --mqttuser=addons			3
	// --mqttpass=airoogheng0phai9ke6FeiR1rohnoop1aong1oocoonah2oocim9aizohweij1mi 	4
	// MqttPort=1883 				5
	// --LogLevel                   6  ]
	var err error

	randId := randStr(10)

	if len(os.Args) <= 5 {
		config.MqttPort = 1883
		config.MqttHost = "192.168.1.117"
		config.MqttClientId = "mqtt_gpio" + "_" + randId
		config.MqttUsername = "mqtt"
		config.MqttPassword = "mqtt"
		config.LogLevel = "warning"
	} else {
		config.MqttPort, err = strconv.Atoi(os.Args[5])
		if err != nil {
			config.MqttPort = 1883
		}
		config.MqttHost = os.Args[1]
		config.MqttClientId = os.Args[2] + "_" + randId
		config.MqttUsername = os.Args[3]
		config.MqttPassword = os.Args[4]
		config.LogLevel = os.Args[6]
	}

	fmt.Printf(
		"MqttPort: %d\nMqttHost: %s\nMqttClientId: %s\nMqttUsername: %s\nMqttPassword: %s\nLogLevel: %s\n",
		config.MqttPort, config.MqttHost, config.MqttClientId, config.MqttUsername, config.MqttPassword, config.LogLevel,
	)
}

func main() {
	loadConfig()
	go gpio_handler.Run(config)
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("tick")
	}
}

func randStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
