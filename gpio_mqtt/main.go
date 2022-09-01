package main

import (
	"fmt"
	"gpio_mqtt/gpio_handler"
	"log"
	"math/rand"
	"net/http"
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
	http.HandleFunc("/", wuiRootHandler)
	err := http.ListenAndServe("localhost:8099", nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(10 * time.Second)
	}
}

func randStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func wuiRootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "URL.Path = %s\n", "test")
}
