package main

import (
	"fmt"
	jrpc "github.com/gumeniukcom/golang-jsonrpc2"
	"gpio_mqtt/gpio_handler"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config = gpio_handler.Config

var config Config

func loadConfig() {
	// $MQTT_HOST $MqttClientId $MQTT_USER $MQTT_PASSWORD $MqttPort $MqttPortWsSsl $LogLevel
	//	1          2             3			4				5			6			7
	var err error

	randId := randStr(10)

	if len(os.Args) <= 5 {
		config.MqttPortWsSsl = 8884
		config.MqttPort = 1883
		config.MqttHost = "192.168.1.117"
		config.MqttClientId = "mqtt_gpio" + "_" + randId
		config.MqttUsername = "mqtt"
		config.MqttPassword = "mqtt"
		config.LogLevel = "warning"
		config.CertFile = "server.crt"
		config.KeyFile = "server.key"
	} else {
		config.MqttPort, err = strconv.Atoi(os.Args[5])
		if err != nil {
			config.MqttPort = 1883
		}
		config.MqttPortWsSsl, err = strconv.Atoi(os.Args[6])
		if err != nil {
			config.MqttPortWsSsl = 8884
		}
		config.MqttHost = os.Args[1]
		config.MqttClientId = os.Args[2] + "_" + randId
		config.MqttUsername = os.Args[3]
		config.MqttPassword = os.Args[4]
		config.LogLevel = os.Args[7]
		config.CertFile = "/ssl/fullchain.pem"
		config.KeyFile = "/ssl/privkey.pem"
	}

	fmt.Printf(
		"MqttHost: %s\nMqttClientId: %s\nMqttUsername: %s\nMqttPassword: %s\nLogLevel: %s\nMqttPort: %d\nMqttPortWsSsl: %d\n",
		config.MqttHost, config.MqttClientId, config.MqttUsername, config.MqttPassword, config.LogLevel, config.MqttPort, config.MqttPortWsSsl,
	)
}

func main() {
	loadConfig()
	config.Rpc = jrpc.New()

	go gpio_handler.Run(config)

	go http.HandleFunc("/rpc", rpcServe)
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	_ = http.ListenAndServe("0.0.0.0:8099", nil)

}

func rpcServe(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	panic(err)
	//}
	//defer r.Body.Close()
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//if _, err = w.Write(config.Rpc.HandleRPCJsonRawMessage(ctx, body)); err != nil {
	//	panic(err)
	//}
	time.Sleep(time.Second)
	w.Write([]byte("ok"))
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

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		lp := filepath.Join("web", "index.html")
		tmpl, _ := template.ParseFiles(lp)
		_ = tmpl.Execute(w, nil)
	}
}
