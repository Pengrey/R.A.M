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

var agentsList = []string{}

type Message struct {
	Type string
	Text string
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

func getArguments() (string, bool) {
	// Get port to be used
	PORT := flag.String("port", "8080", "Port to be used for communication")

	// Define if program is to be run on silent mode
	SILENT := flag.Bool("s", false, "Remove prompt from startup")

	// Call flag.Parse() to parse the command-line flags
	flag.Parse()

	return *PORT, *SILENT
}

func getPrompt() {
	version := "0.1 beta"
	prompt := " (                          *     \n )\\ )         (           (  `    \n(()/(         )\\          )\\))(   \n /(_))     ((((_)(       ((_)()\\  \n(_))        )\\ _ )\\      (_()((_) \n| _ \\       (_)_\\(_)     |  \\/  | \n|   /   _    / _ \\    _  | |\\/| | \n|_|_\\  (_)  /_/ \\_\\  (_) |_|  |_| "

	fmt.Println(prompt)
	fmt.Printf("\n[Remote Anamnestic Mapper (v%s)]\n\n", version)
	fmt.Println("[?] Available commands:\n       info - Get server info\n       addA - Add Agent\n       reqR - Request Agent's RAM\n       help - See available commands\n       quit - Quit")
}

func addAgent(port string) {
	fmt.Println("[+] Listing Agents.")
	// Receive msg
	msg := receiveMessage(port)

	// Check if message is hello type
	if msg.Type == "hello" {
		agentsList = append(agentsList, msg.Text)
		fmt.Printf("[+] Agent [%s] added.\n", msg.Text)
	} else {
		fmt.Println("[!] ERROR: message given is misstyped.")
	}
}

func receiveMessage(port string) Message {
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	conn, err := listen.Accept()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// close listener
	listen.Close()

	// incoming request
	buffer := make([]byte, 1024)
	read_len, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Respond OK
	conn.Write([]byte("OK"))

	// close conn
	conn.Close()

	// Get json from bytes
	var msg Message
	err = json.Unmarshal(buffer[:read_len], &msg)
	if err != nil {
		fmt.Printf("[!] An error occured during json parsing: %v\n", err)
		os.Exit(1)
	}

	return msg
}

func sendMessage(message Message, addr string) {
	tcpServer, err := net.ResolveTCPAddr("tcp", addr)

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

func requestRAM(port string) {
	// Send message requesting ram
	msg := Message{
		Type: "ram",
		Text: "RAM",
	}

	if len(agentsList) > 0 {
		fmt.Println("[+] Listing Agents saved:")
		fmt.Println("+-------+--------------------------+")
		fmt.Println("| Index |    Agent's IP Address    |")
		fmt.Println("+-------+--------------------------+")
		for index, element := range agentsList {
			fmt.Printf("|   %d   |   %s    |\n", index, element)
			fmt.Println("+-------+--------------------------+")
		}

		var inpt int

		for true {
			fmt.Println("[?] Insert Index of the agent to request")
			fmt.Print("> ")

			fmt.Scanln(&inpt)

			if inpt >= 0 && inpt < len(agentsList) {
				fmt.Println("[+] Requesting RAM.")
				sendMessage(msg, agentsList[inpt])
				retreiveRAM(port)
				break
			} else {
				fmt.Println("[!] Invalid agent index!")
			}
		}
	} else {
		fmt.Println("[!] No Agents are saved!")
	}
}

func retreiveRAM(port string) {
	fmt.Println("[+] Retreiving RAM from agent.")
	// Prepare command
	cmd := fmt.Sprintf("nc -l -p %s > ram.txt", port)

	// Run command with shell
	err := exec.Command("sh", "-c", cmd).Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[+] RAM saved into ram.txt.")
}

func menu(port string) {
	var inpt string
	for true {
		fmt.Print("> ")

		fmt.Scanln(&inpt)

		switch inpt {
		case "help":
			fmt.Println("[?] Available commands:\n       info - Get server info\n       addA - Add Agent\n       reqR - Request Agent's RAM\n       help - See available commands\n       quit - Quit")
		case "info":
			fmt.Printf("[*] Server info:\n       IP: %s\n       Port: %s\n", getIP(), port)
		case "addA":
			addAgent(port)
		case "reqR":
			requestRAM(port)
		case "quit":
			os.Exit(0)
		default:
			fmt.Printf("[!] Option %s doesnt exist!\n", inpt)
		}
	}
}

func main() {
	// Get arguments from user
	port, silent := getArguments()

	// Print Prompt
	if !silent {
		getPrompt()
	}

	// Start Menu
	menu(port)
}