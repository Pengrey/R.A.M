package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type Message struct {
	Type string
	Text string
}

func getArguments() (string, string) {
	// Get port to be used
	PORT := flag.String("port", "8080", "Port to be used for communication")

	// Defining a string flag
	PASSWD := flag.String("password", "", "Password to be used for authentication of the agents")

	// Call flag.Parse() to parse the command-line flags
	flag.Parse()

	return *PORT, *PASSWD
}

// Get preferred outbound ip of this machine
func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func sendMessage(message Message) {
	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:8080")

	if err != nil {
		println("[!] ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	// Start connection
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		println("[!] Dial failed:", err.Error())
		os.Exit(1)
	}

	// Get message json
	res, err := json.Marshal(message)

	// Send data
	_, err = conn.Write(res)
	if err != nil {
		println("[!] Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("[!] Read data failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("[*] Status:", string(received))

	conn.Close()
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	read_len, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Get json from bytes
	var msg Message
	err = json.Unmarshal(buffer[:read_len], &msg)
	if err != nil {
		fmt.Printf("[!] An error occured during json parsing: %v\n", err)
		conn.Write([]byte("ERROR"))
		os.Exit(1)
	}

	// Send Response
	if msg.Text != "RAM" {
		conn.Write([]byte("ERROR"))
	} else {
		conn.Write([]byte("OK"))
		fmt.Println("[+] RAM dump requested")
	}

	// close conn
	conn.Close()
}

func sayHello(name string) {
	fmt.Println("[+] Sending hello message.")
	helloMsg := Message{
		Type: "hello",
		Text: name,
	}
	sendMessage(helloMsg)
}

func startServer(port string) {
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Send Hello message
	sayHello(fmt.Sprintf("%s:%s", getIP(), port))

	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func main() {
	// Get arguments from user
	port, _ := getArguments()

	startServer(port)
}
