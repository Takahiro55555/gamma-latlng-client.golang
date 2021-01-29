package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	client "github.com/Takahiro55555/location-based-mqtt-client.golang"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	host := flag.String("host", "localhost", "Gateway MQTT broker host")
	port := flag.Int("port", 1884, "Gateway MQTT broker port")
	flag.Parse()
	log.SetOutput(os.Stdout)
	log.Printf("[LOG] host: %v, port: %v", *host, *port)

	//////////////            軌跡データの読み込み準備           //////////////
	ch := make(chan mqtt.Message)

	dir := "./test/"
	files, _ := ioutil.ReadDir(dir)
	go func() {
		for i, f := range files {
			log.Printf("Loading file: %v, %v/%v", dir+f.Name(), i+1, len(files))
			publishTrajectory(ch, dir+f.Name(), *host, *port)
		}
	}()

	//////////////           勝手に終了しないようにする          //////////////
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	for {
		select {
		case m := <-ch:
			log.Printf("[LOG] Recieved Topic: %v, Message: %v\n", m.Topic(), string(m.Payload()))
		case <-signalCh:
			log.Printf("Interrupt detected.\n")
			return
		}
	}
}

func publishTrajectory(ch chan mqtt.Message, fileName string, host string, port int) {

	gatewayBrokerHost := fmt.Sprintf("tcp://%v:%v", host, port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(gatewayBrokerHost)

	// ゲートウェイブローカへ接続
	gatewayBroker := mqtt.NewClient(opts)
	if token := gatewayBroker.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}
	defer gatewayBroker.Disconnect(1000)

	var callback mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		ch <- msg
	}

	c := client.NewClient(gatewayBroker, 30.)

	log.Printf("Open: %v", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	counter := 1
	for {
		record, err := reader.Read()
		if err != nil {
			log.Print(err)
			break
		}
		log.Printf("Counter: %v (file: %v)", counter, fileName)
		lng, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		lat, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		if err := c.UpdateSubscribe(lat, lng, 0, callback); err != nil {
			log.Fatalf("Mqtt error: %s", err)
		}

		client_id := "hoge"
		payload := fmt.Sprintf("{\"client_id\":\"%v\",\"objects\":{\"cars\":[{\"lat\":%v,\"lng\":%v}],\"bicycles\":[],\"humans\":[],\"other\":{}}}", client_id, lat, lng)
		if err := c.Publish(lat, lng, 0, false, payload); err != nil {
			log.Fatalf("Mqtt error: %s", err)
		}
		time.Sleep(time.Millisecond * 10)
		counter++
	}
	c.Unsubscribe()
	log.Printf("Close: %v", fileName)
}
