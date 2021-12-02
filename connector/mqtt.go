package connector

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

type Node struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	Status        string `json:"status" bson:"status"`
	WorkerID      string `json:"worker_id" bson:"worker_id"`
	JobID         string `json:"job_id" bson:"job_id"`
	FriendlyNames struct {
		De string `json:"de" bson:"de"`
		En string `json:"en" bson:"en"`
	} `json:"friendly_names" bson:"friendly_names"`
	VoteParticipant bool   `json:"vote_participant" bson:"vote_participant"`
	Schema          string `json:"schema" bson:"schema,omitempty"`
}

type Message struct {
	LamportTs int    `json:"lamport_ts" bson:"lamport_ts"`
	WorkerID  string `json:"worker_id" bson:"worker_id"`
	MsgType   string `json:"msg_type" bson:"msg_type"`
	Content   Node   `json:"content" bson:"content"`
	OrderID   string `json:"order_id" bson:"order_id"`
}

func ConnectMQTT(host string) error {

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}

	var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		var m Message
		if err := json.Unmarshal(msg.Payload(), &m); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v \n\n", string(msg.Payload()))

		var node Node
		if err := FindNodeByWorkerID(m.WorkerID, &node); err != nil {
			if err == mongo.ErrNoDocuments {
				CreateNode(&m.Content)
			} else {
				log.Fatal(err)
			}

		}
	}

	opts := mqtt.NewClientOptions().
		AddBroker("ws://" + host + ":9001/mqtt").
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