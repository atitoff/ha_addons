package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("MqttPort:", os.Getenv("MqttPort"))
	fmt.Println(len(os.Args), os.Args)

	content, err := ioutil.ReadFile("/data/options.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("tick")
	}
}
