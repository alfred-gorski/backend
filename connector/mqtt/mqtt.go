package mqtt_conn

import (
	"backend/config"
	"backend/controller"
	"backend/models"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConnectMQTT(m config.MQTTConfig) error {

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}

	var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		var m models.Message
		if err := json.Unmarshal(msg.Payload(), &m); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v \n\n", string(msg.Payload()))

		var node models.Node
		if err := controller.FindNodeByWorkerID(m.WorkerID, &node); err != nil {
			if err == mongo.ErrNoDocuments {
				m.Content.Schema = "{}"
				m.Content.FormData = "{}"
				controller.CreateNode(&m.Content)
			} else {
				log.Fatal(err)
			}

		}
	}
	url := fmt.Sprintf("ws://%s:%d/%s", m.Addr, m.Port, m.Name)

	opts := mqtt.NewClientOptions().
		AddBroker(url).
		SetClientID("go-simple")

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	topics := map[string]byte{
		"gaia/status/#":                         0,
		"gaia/systemoverview":                   0,
		"gaia/systemoverview/systemoverview1/#": 0,
		"gaia/+/response/systemoverview1":       0,
	}

	if token := c.SubscribeMultiple(topics, messageSubHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
