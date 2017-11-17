package main

import (
	//"bufio"
	//"crypto/tls"
	"flag"
	"fmt"
	//"io"
	//"log"
	//"os"
	//"strconv"
	//"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	INFO "github.com/shirou/gopsutil/mem"
)

type deviceinfo struct{
	clientid:
	deviceid:
	memused:
	memavailable:
}

type  meminfo struct{
	clientid:
	deviceid:
	memused:
	memavailable:
}


	
	

func main() {
	server := flag.String("server", "tcp://10.0.2.15:1883", "The full URL of the MQTT server to connect to")
	topic1 := flag.String("topic1", "one", "Topic to publish the messages on")
	topic2 := flag.String("topic2", "two", "Topic to publish the messages on")
	topic := []*string{topic1,topic2}
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	
	clientid1 := flag.String("clientid1", "client1", "A clientid1 for the connection")
	clientid2 := flag.String("clientid2", "client2", "A clientid2 for the connection")
	clientid3 := flag.String("clientid3", "client3", "A clientid3 for the connection")
	clientid4 := flag.String("clientid4", "client4", "A clientid4 for the connection")
	clients := []*string{clientid1,clientid2,clientid3,clientid4}
	flag.Parse()
	 
	c:=make(chan string)
	for k:=0; k<3; k++{ 
		go create_client(clients[k], c,server,topic[0],qos,retained)
	}
	go create_client(clients[3],c,server,topic[1],qos,retained)

	for l:=0; l<13; l++ {
		fmt.Println(<-c)
	}
}

func create_client(clientids *string, c chan string,server *string,topic *string,qos *int,retained *bool){
   
		connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientids).SetCleanSession(true)

		client := MQTT.NewClient(connOpts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			return
	        	}
		fmt.Printf("Connected to %s\n", *server)

		for j:=0;j<3;j++ {
		client.Publish(*topic, byte(*qos), *retained, "message")
		c <- "message received"
		}
}
