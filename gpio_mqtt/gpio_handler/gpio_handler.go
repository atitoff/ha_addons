package gpio_handler

import (
	"encoding/hex"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"
)

type Config struct {
	MqttPort     int
	MqttHost     string
	MqttClientId string
	MqttUsername string
	MqttPassword string
	LogLevel     string
}

type PublishTopic struct {
	Topic   string
	Payload string
}

var sendQueue = make(chan PublishTopic, 100)

var client mqtt.Client

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
	subscribe(client)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func Run(settings Config) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", settings.MqttHost, settings.MqttPort))
	opts.SetClientID(settings.MqttClientId)
	opts.SetUsername(settings.MqttUsername)
	opts.SetPassword(settings.MqttPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	var publishData PublishTopic
	for {
		// получаем данные для передачи и передаем
		publishData = <-sendQueue
		token := client.Publish(publishData.Topic, 0, false, publishData.Payload)
		token.Wait()
	}
}

func ErrorToMqtt(err string) {
	sendQueue <- PublishTopic{"log/error", err}
}

func daliSet(topic string, splitTopic []string) {
	var errStr string = fmt.Sprintf("dali mismatch: %s", topic)
	// gpio/2/dali/set/25/255
	// 25 адрес устройства (00..63)
	// 255 устанавливаемое значение (0..255)
	// gpio/2/dali/set_grp/5/255
	// 5 группа (0..15)

	if splitTopic[3] == "set" {
		addr, err := strconv.Atoi(splitTopic[4])
		if err != nil {
			ErrorToMqtt(errStr)
			return
		}
		value, err := strconv.Atoi(splitTopic[5])
		if err != nil {
			ErrorToMqtt(errStr)
			return
		}
		if (addr > 63 || addr < 0) || (value > 255 || value < 0) {
			ErrorToMqtt(errStr)
			return
		}
		sendQueue <- PublishTopic{
			fmt.Sprintf("gpio_%s/dali", splitTopic[1]),
			string([]byte{uint8(addr << 1), uint8(value)})}
	} else if splitTopic[3] == "set_grp" {
		addr, err := strconv.Atoi(splitTopic[4])
		if err != nil {
			ErrorToMqtt(errStr)
			return
		}
		value, err := strconv.Atoi(splitTopic[5])
		if err != nil {
			ErrorToMqtt(errStr)
			return
		}
		if (addr > 15 || addr < 0) || (value > 255 || value < 0) {
			ErrorToMqtt(errStr)
			return
		}
		sendQueue <- PublishTopic{
			fmt.Sprintf("gpio_%s/dali", splitTopic[1]),
			string([]byte{uint8(0b10000000 | (addr << 1)), uint8(value)})}
	}

}

func daliRaw(msg mqtt.Message) {
	// gpio/2/dali/raw
	// payload 2 byte | 4 byte hex convert to 2 byte
	if len(msg.Payload()) == 4 {
		data, err := hex.DecodeString(string(msg.Payload()))
		if err == nil {
			fmt.Println(data)
		}
	}
}

func subscribe(client mqtt.Client) {
	topic := "gpio/#"
	token := client.Subscribe(topic, 1, gpio)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

func gpio(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(msg.Topic())
	splitTopic := strings.Split(msg.Topic(), "/")
	fmt.Println(splitTopic, len(splitTopic))
	if len(splitTopic) >= 6 {
		if splitTopic[2] == "dali" {
			if splitTopic[3] == "set" || splitTopic[3] == "set_grp" {
				daliSet(msg.Topic(), splitTopic)
			}
		}
	} else if len(splitTopic) == 4 {
		if splitTopic[2] == "dali" && splitTopic[3] == "raw" {
			daliRaw(msg)
		}
	}

}
