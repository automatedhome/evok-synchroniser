package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	Status string `json:"status"`
	Data   []struct {
		Value   float64 `json:"value"`
		Circuit string  `json:"circuit"`
		Dev     string  `json:"dev"`
	} `json:"data"`
}

func main() {
	broker := flag.String("broker", "tcp://127.0.0.1:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	clientID := flag.String("clientid", "evoksync", "A clientid for the connection")
	evok := flag.String("evok", "http://127.0.0.1:8080/json/all", "Address of endpoint exposing all sensors data")
	interval := flag.Int("interval", 15, "Interval between synchronisations")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID(*clientID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)
	MQTTClient := mqtt.NewClient(opts)
	if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	time.Sleep(5 * time.Second)
	for {
		response, err := http.Get(*evok)
		if err != nil {
			log.Fatalf("Couldn't connect to EVOK: %v", err)
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Couldn't read EVOK data: %v", err)
		}

		data := Message{}
		err = json.Unmarshal([]byte(contents), &data)
		if err != nil {
			log.Printf("Failed to unmarshal JSON data from EVOK message: %v\n", err)
		}

		for _, sensor := range data.Data {
			if sensor.Dev != "temp" {
				continue
			}
			topic := "evok/" + sensor.Dev + "/" + sensor.Circuit + "/value"
			token := MQTTClient.Publish(topic, 0, false, fmt.Sprintf("%v", sensor.Value))
			token.Wait()
			if token.Error() != nil {
				log.Printf("Failed to publish packet: %s", token.Error())
			}
		}

		time.Sleep(time.Duration(*interval*60) * time.Minute)
	}
}
