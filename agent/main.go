package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

type Message struct {
	Type     string
	Text     string
	LimePort string
}

func getArguments() (string, string, string, string) {
	// Get port to be used
	LPORT := flag.String("LPORT", "8081", "Port to be used for communication")

	// Get ip of the server to be used
	RHOST := flag.String("RHOST", getIP(), "IP of the server")

	// Get port of the server to be used
	RPORT := flag.String("RPORT", "8080", "Port of the server")

	LimePORT := flag.String("LIMEP", "4444", "Port LiME sends RAM through")
	// Call flag.Parse() to parse the command-line flags
	flag.Parse()

	return *LPORT, *RHOST, *RPORT, *LimePORT
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

func sendRam(RHOST string, LPORT string) {
	fmt.Println("[+] Sending RAM dump.")

	// Prepare command

	//cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("memdump -s 409600 > /dev/tcp/%s/%s", RHOST, RPORT)) // Used for tests
	//cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("insmod  -s 409600 > /dev/tcp/%s/%s", RHOST, RPORT)) // Used for lime tests
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("insmod /var/lib/dkms/lime-forensics/1.9.1-3/6.0.6-76060006-generic/x86_64/module/lime.ko \"path=tcp:%s format=lime\"", LPORT)) // Used for lime tests
	// cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("memdump > /dev/tcp/%s/%s", RHOST, RPORT)) // Used in production

	// Execute command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[+] RAM dump sent.")
	exec.Command("/bin/bash", "-c", "rmmod lime").Run()
}

func handleRequest(conn net.Conn, RHOST string, LPORT string) {
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
		sendRam(RHOST, LPORT)
	}

	// close conn
	conn.Close()
}

func sayHello(RHOST string, RPORT string, name string, LimePORT string) {
	fmt.Println("[+] Sending hello message.")

	tcpServer, err := net.ResolveTCPAddr("tcp", RHOST+":"+RPORT)

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

	// Prepare message
	helloMsg := Message{
		Type:     "hello",
		Text:     name,
		LimePort: LimePORT,
	}

	// Get message json
	res, err := json.Marshal(helloMsg)

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

func startServer(LPORT string, RHOST string, RPORT string, LimePORT string) {
	fmt.Println("[+] Starting server.")
	listen, err := net.Listen("tcp", "0.0.0.0:"+LPORT)
	fmt.Println("[+] Listening on port", LPORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Send Hello message
	sayHello(RHOST, RPORT, fmt.Sprintf("%s:%s", getIP(), LPORT), LimePORT)

	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn, RHOST, LPORT)
	}
}

func main() {
	// Get arguments from user
	LPORT, RHOST, RPORT, LimePORT := getArguments()

	startServer(LPORT, RHOST, RPORT, LimePORT)
}
