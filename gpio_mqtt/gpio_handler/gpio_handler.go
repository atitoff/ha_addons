package gpio_handler

import (
	"encoding/hex"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"regexp"
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

type CompiledRegex struct {
	daliSet    *regexp.Regexp
	daliSetGrp *regexp.Regexp
	daliRaw    *regexp.Regexp
	gpioSet    *regexp.Regexp
}

var compiledRegex CompiledRegex

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

// 0 <= index <= len(a)
func shiftArray(array *[]string, position int, value string) {
	//  extend array by one
	*array = append(*array, "")

	// shift values
	copy((*array)[position+1:], (*array)[position:])

	// insert value
	(*array)[position] = value
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Run(settings Config) {
	compileRegex()
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
	// var errStr string = fmt.Sprintf("dali mismatch: %s", topic)
	// GPIO/2/DALI/SET/25/255
	// 25 адрес устройства (00..63)
	// 255 устанавливаемое значение (0..255)
	fmt.Println(topic)
	addr, err := strconv.Atoi(splitTopic[4])
	if err != nil {
		return
	}
	value, err := strconv.Atoi(splitTopic[5])
	if err != nil {
		return
	}
	payload := []byte{uint8(addr << 1), uint8(value)}
	sendQueue <- PublishTopic{
		fmt.Sprintf("GPIO/SUB/%s", splitTopic[1]),
		"0003" + hex.EncodeToString(payload)}

}

func daliSetGrp(topic string, splitTopic []string) {
	// GPIO/2/DALI/SET_GRP/12/255
	fmt.Println(topic)
	addr, err := strconv.Atoi(splitTopic[4])
	if err != nil {
		return
	}
	value, err := strconv.Atoi(splitTopic[5])
	if err != nil {
		return
	}

	payload := []byte{uint8(0b10000000 | (addr << 1)), uint8(value)}
	sendQueue <- PublishTopic{
		fmt.Sprintf("GPIO/SUB/%s", splitTopic[1]),
		"0003" + hex.EncodeToString(payload)}
}

func daliRaw(msg mqtt.Message, splitTopic []string) {
	// GPIO/2/DALI/RAW
	// payload 4 byte hex convert to 2 byte
	_, err := strconv.ParseUint(string(msg.Payload()), 16, 64)
	if err != nil {
		return
	}
	if len(msg.Payload()) != 4 {
		return
	}

	sendQueue <- PublishTopic{
		fmt.Sprintf("GPIO/SUB/%s", splitTopic[1]),
		"0003" + string(msg.Payload())}
}

func gpioSet(topic string, splitTopic []string) {
	// GPIO/2/0/SET/3/255
	// 0 - подмодуль
	// 3 - порт
	// 255 устанавливаемое значение (0..255)
	port, err := strconv.Atoi(splitTopic[4])
	if err != nil {
		return
	}
	value, err := strconv.Atoi(splitTopic[5])
	if err != nil {
		return
	}
	subModule, err := strconv.Atoi(splitTopic[2])
	if err != nil {
		return
	}
	payload := []byte{uint8(subModule), 0x00, uint8(port), uint8(value)}
	sendQueue <- PublishTopic{
		fmt.Sprintf("GPIO/SUB/%s", splitTopic[1]),
		string(payload)}

}

func subscribe(client mqtt.Client) {
	topic := "GPIO/PUB/+"
	token := client.Subscribe(topic, 0, evGpio)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
	topic = "GPIO/+/SET/#"
	token = client.Subscribe(topic, 0, setGpio)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
	topic = "GPIO/+/+/SET/#"
	token = client.Subscribe(topic, 0, setGpio)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
	topic = "GPIO/+/DALI/#"
	token = client.Subscribe(topic, 0, sendToDali)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)

}

// send data to GPIO module GPIO/SUB/+
func setGpio(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(msg.Topic(), msg.Payload())
	splitTopic := strings.Split(msg.Topic(), "/")

	if len(splitTopic) == 5 {
		shiftArray(&splitTopic, 2, "0")
	}

	fmt.Println(splitTopic)
	if compiledRegex.gpioSet.MatchString(msg.Topic()) {
		gpioSet(msg.Topic(), splitTopic)
	}
}

func sendToDali(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(msg.Topic(), msg.Payload())
	splitTopic := strings.Split(msg.Topic(), "/")
	if compiledRegex.daliSet.MatchString(msg.Topic()) {
		daliSet(msg.Topic(), splitTopic)
	} else if compiledRegex.daliSetGrp.MatchString(msg.Topic()) {
		daliSetGrp(msg.Topic(), splitTopic)
	} else if compiledRegex.daliRaw.MatchString(msg.Topic()) {
		daliRaw(msg, splitTopic)
	}
}

// receive GPIO/PUB/+ from GPIO modules
// convert to GPIO/+/EV/...
func evGpio(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(msg.Topic(), msg.Payload())
}

func compileRegex() {
	v0255 := `(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
	v015 := `(\d|1[0-5])`
	v063 := `(\d|[1-5]\d|6[0-3])`
	compiledRegex.daliSet = regexp.MustCompile(`^GPIO/` + v0255 + `/DALI/SET/` + v063 + `/` + v0255 + `$`)
	compiledRegex.daliSetGrp = regexp.MustCompile(`^GPIO/` + v0255 + `/DALI/SET_GRP/` + v015 + `/` + v0255 + `$`)
	compiledRegex.daliRaw = regexp.MustCompile(`^GPIO/` + v0255 + `/DALI/RAW$`)
	compiledRegex.gpioSet = regexp.MustCompile(`^GPIO/(` + v0255 + `|` + v0255 + `/` + v0255 + `)/SET/` + v0255 + `/` + v0255 + `$`)
}
