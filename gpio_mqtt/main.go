package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type data struct {
	Data map[string]string `yaml:"options"`
}

func loadSettings() {
	t := data{}
	// read settings
	settFile, err := os.ReadFile("/data/config.yaml")
	check(err)
	err = yaml.Unmarshal(settFile, &t)
	check(err)
	if foo, ok := t.Data["MqttPort"]; ok {
		fmt.Printf("--- t:\n%v\n\n", foo)
	} else {
		panic("not in config MqttPort")
	}
}

func main() {
	loadSettings()
	fmt.Println(len(os.Args), os.Args)
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("tick")
	}
}
