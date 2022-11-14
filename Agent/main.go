package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
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

	// Start connection, if error
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

func sendRam() {
	fmt.Println("[+] Sending RAM dump.")
	// Send RAM file size to server
	ramFile, err := os.Open("ram.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer ramFile.Close()

	// Get file size
	ramFileStat, err := ramFile.Stat()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Send file chunk number
	chunkNumber := ramFileStat.Size()/512 + 1
	ramSize := Message{
		Type: "size",
		Text: fmt.Sprintf("%d", chunkNumber),
	}
	// Sleep to avoid flooding
	time.Sleep(1 * time.Second)
	// Send file
	sendMessage(ramSize)

	// Read file content chunk by chunk, without loading the whole file into memory
	file, err := os.Open("ram.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Send file content
	buffer := make([]byte, 512)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			break
		}
		// Send chunk to server where the Text field is the chunk encoded in base64
		ramChunk := Message{
			Type: "ram",
			Text: base64.StdEncoding.EncodeToString(buffer),
		}

		// Sleep to avoid flooding
		time.Sleep(300 * time.Millisecond)
		sendMessage(ramChunk)

		// Clear buffer
		buffer = make([]byte, 512)
	}
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
		sendRam()
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
