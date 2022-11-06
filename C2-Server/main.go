package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var agentsList = []string{}

type Message struct {
	Type string
	Text string
}

func getArguments() (string, string, bool) {
	// Get port to be used
	PORT := flag.String("port", "8080", "Port to be used for communication")

	// Defining a string flag
	PASSWD := flag.String("password", "", "Password to be used for authentication of the agents")

	// Define if program is to be run on silent mode
	SILENT := flag.Bool("s", false, "Remove prompt from startup")

	// Call flag.Parse() to parse the command-line flags
	flag.Parse()

	return *PORT, *PASSWD, *SILENT
}

func getPrompt() {
	version := "0.1 beta"
	prompt := " (                          *     \n )\\ )         (           (  `    \n(()/(         )\\          )\\))(   \n /(_))     ((((_)(       ((_)()\\  \n(_))        )\\ _ )\\      (_()((_) \n| _ \\       (_)_\\(_)     |  \\/  | \n|   /   _    / _ \\    _  | |\\/| | \n|_|_\\  (_)  /_/ \\_\\  (_) |_|  |_| "

	fmt.Println(prompt)
	fmt.Printf("\n[Remote Anamnestic Mapper (v%s)]\n\n", version)
	fmt.Println("[?] Available commands:\n       addA - Add Agent\n       reqR - Request Agent's RAM\n       help - See available commands\n       quit - Quit")
}

func addAgent(port string, password string) {
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
	listen, err := net.Listen("tcp", "localhost:"+port)
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

func requestRAM(port string, password string) {
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
				break
			} else {
				fmt.Println("[!] Invalid agent index!")
			}
		}
	} else {
		fmt.Println("[!] No Agents are saved!")
	}
}

func menu(port string, password string) {
	var inpt string
	for true {
		fmt.Print("> ")

		fmt.Scanln(&inpt)

		switch inpt {
		case "help":
			fmt.Println("[?] Available commands:\n       addA - Add Agent\n       reqR - Request Agent's RAM\n       help - See available commands\n       quit - Quit")
		case "addA":
			addAgent(port, password)
		case "reqR":
			requestRAM(port, password)
		case "quit":
			os.Exit(0)
		default:
			fmt.Printf("[!] Option %s doesnt exist!\n", inpt)
		}
	}
}

func main() {
	// Get arguments from user
	port, password, silent := getArguments()

	// Print Prompt
	if !silent {
		getPrompt()
	}

	// Start Menu
	menu(port, password)
}
