package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MqttPort int
	MqttHost string
}

var config Config

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadConfig() {
	var err error
	config.MqttPort, err = strconv.Atoi(os.Getenv("MqttPort"))
	config.MqttHost = os.Getenv("MqttPort")
	check(err)
}

func main() {

	fmt.Println(config)

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("tick")
	}
}
