package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	//"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)


var j=0

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	j = j+1
	fmt.Printf("Received message on topic: %s\nMessage: %s\nCount: %d\n", message.Topic(), message.Payload(),j)
}

var i int64

func main() {
	c := make(chan os.Signal, 1)
	i = 0
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("signal received, exiting")
		os.Exit(0)
	}()

	hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://10.0.2.15:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	sub_multiple:= []string{"one","two","three"}
	topic := sub_multiple //flag.String("topic", sub_multiple, "Topic to subscribe to")
	flag.Parse()


	connOpts := &MQTT.ClientOptions{
		ClientID:             *clientid,
		CleanSession:         true,
		Username:             *username,
		Password:             *password,
		MaxReconnectInterval: 1 * time.Second,
		//KeepAlive:            10 * time.Second,
		TLSConfig:            tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert},
	}
	connOpts.AddBroker(*server)
	connOpts.OnConnect = func(c MQTT.Client) {
		for i,_:=range sub_multiple{

		if token := c.Subscribe(topic[i], byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
                      }
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *server)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
