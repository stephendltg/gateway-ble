package broker

import (
	"encoding/json"
	"fmt"
	"gateway-ble/store"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// Var
var client = mqtt.NewClient(mqtt.NewClientOptions())

// Connect to broker
func Connect(host *string, port *string, username *string, password *string) {

	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// debugger
	debugger := log.WithFields(log.Fields{"package": "BROKER"})

	// Default port
	PORT := "1883"
	if port != nil {
		PORT = *port
	}

	// Default hostname
	HOST := "127.0.0.1"
	if host != nil {
		HOST = *host
	}

	// Default username
	USERNAME := "emqx"
	if username != nil {
		USERNAME = *username
	}

	// Default password
	PASSWORD := "public"
	if password != nil {
		PASSWORD = *password
	}

	// Options broker
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", HOST, PORT))
	opts.SetClientID("gateway_client_" + uuid.New().String())
	opts.SetUsername(USERNAME)
	opts.SetPassword(PASSWORD)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(3 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.WithFields(log.Fields{
			"playload": fmt.Sprintf("%s", msg.Payload()),
			"topic":    msg.Topic(),
			"qos":      msg.Qos(),
			"retain":   msg.Retained(),
			"Id":       msg.MessageID(),
		}).Info("Broker: Received message")
	})

	// Connection lost handler
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		store.Set("mqtt:status", "lost error")
		debugger.Trace("Connection lost error", err.Error())
	})
	// Broker reconnection
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		store.Set("mqtt:status", "reconnect")
		debugger.Trace("Reconnecting")
	})
	// Broker connect
	opts.OnConnect = func(c mqtt.Client) {
		debugger.Info("Connected")
		store.Set("mqtt:status", "connected")
	}
	//Broker connection lost
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		store.Set("mqtt:status", "lost")
		debugger.Warn("Connect lost: ", err.Error())
	}

	// Connection
	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		debugger.Warn(token.Error())
	} else {
		// Subscribe
		// topic := "gateway/epyo/status"
		// subscribe := client.Subscribe(topic, 1, nil)
		// subscribe.Wait()
		// debugger.Info(fmt.Sprintf("Subscribed to topic %s", topic))
	}

	// Cron jobs
	c := cron.New()
	interval := store.Get("mqtt:interval", nil)
	c.AddFunc("@every "+interval, func() {
		Publish(store.BeaconsList())
	})
	c.Start()

	// Test
	// publish(client)
}

// Publish
func Publish(playload interface{}) {

	debugger := log.WithFields(log.Fields{"package": "BROKER:PUBLISH"})

	// b, err := json.Marshal(playload)
	// if err != nil {
	// 	debugger.Warn("JSON Marsham Error", err)
	// } else {
	// 	token := client.Publish("gateway/ble/status", 0, false, b)
	// 	token.Wait()
	// }

	jsonString, err := json.Marshal(playload)
	if err != nil {
		debugger.Warn("JSON Marsham Error", err)
	}
	token := client.Publish("gateway/ble/status", 0, false, jsonString)
	token.Wait()

}

// TEST
// func publish(client mqtt.Client) {
// 	num := 10
// 	for i := 0; i < num; i++ {
// 		text := fmt.Sprintf("Message %d", i)
// 		token := client.Publish("gateway/ble/status", 0, false, text)
// 		token.Wait()
// 		time.Sleep(time.Second)
// 	}
// }
