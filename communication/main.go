package main

import (
	"fmt"
	"os"
	"time"
)

type Message struct {
	From    string
	Payload string
}

type Server struct {
	msgCh  chan Message
	quitCh chan struct{} // empty struct uses 0 bytes in memory
}

func (s *Server) StartAndListen() {
	// can name a for loop
running:
	for {
		select {
		// blocks here until channel receives a message
		case msg := <-s.msgCh:
			fmt.Printf("message received from: %s - payload: %s\n", msg.From, msg.Payload)
		case <-s.quitCh:
			fmt.Println("server is doing gracefull shutdown")
			// logic for gracefull shutdown
			break running
		default:
		}
	}
	fmt.Println("server down!")
	os.Exit(0)
}

func sendMessageToServer(msgCh chan Message, payload string) {
	msg := Message{
		From:    "Calypsa",
		Payload: payload,
	}
	msgCh <- msg
}

func gracefullQuitServer(quitCh chan struct{}) {
	close(quitCh)
}

func main() {
	fmt.Println("Concurrency communication!")
	server := &Server{
		msgCh:  make(chan Message),
		quitCh: make(chan struct{}),
	}
	go server.StartAndListen()

	go func() {
		time.Sleep(2 * time.Second)
		sendMessageToServer(server.msgCh, "Guau guau guau")
	}()

	go func() {
		time.Sleep(4 * time.Second)
		gracefullQuitServer(server.quitCh)
	}()

	select {}
}
