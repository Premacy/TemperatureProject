package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type connection_data struct {
	ip   string
	port int
}

type sensor_data struct {
	temperature string
	pressure    string
}

func parseJsonMessage(bytes []byte, data *sensor_data) error {
	err := json.Unmarshal(bytes, data)

	return err
}

func getOpts(conn_data *connection_data) *mqtt.ClientOptions {

	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		var data = new(sensor_data)
		parseJsonMessage(msg.Payload(), data)

		fmt.Printf("Temperature: %s, pressure %s", data.temperature, data.pressure)
	}

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conn_data.ip, conn_data.port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return opts
}

func listen() {

}

func main() {
	fmt.Println("Start server")

	topics := [2]string{
		"sensor/temperature_pressure",
	}

	var conn_data = connection_data{"broker.emqx.io", 1883}
	var opts = getOpts(&conn_data)

	client := mqtt.NewClient(opts)

	const connection_timeout = 5 * time.Second

	var token = client.Connect()

	if !token.WaitTimeout(connection_timeout) {
		fmt.Println("Erorr. Connection timeout")
		return
	}

	client.Subscribe(topics[0], 1, nil)
	client.Subscribe(topics[1], 1, nil)

	for {

	}
}
