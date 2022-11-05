package main

import (
	"log"
	"math/rand"

	p2p "github.com/leprosus/golang-p2p"
)

type Message struct {
	Type string
	Text string
}

func RandomName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func sendMessage(message Message) {
	tcp := p2p.NewTCP("localhost", "8080")

	client, err := p2p.NewClient(tcp)
	if err != nil {
		log.Panicln(err)
	}

	req := p2p.Data{}
	err = req.SetGob(message)
	if err != nil {
		log.Panicln(err)
	}

	_, err = client.Send("dialog", req)
	if err != nil {
		log.Panicln(err)
	}
}

func sayHello() {
	helloMsg := Message{
		Type: "hello",
		Text: RandomName(10),
	}
	sendMessage(helloMsg)
}

func main() {
	sayHello()
}
