package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	p2p "github.com/leprosus/golang-p2p"
)

type Message struct {
	Type string
	Text string
}

// Rewrite Logger
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type stdLogger struct{}

func NewStdLogger() (l *stdLogger) {
	return &stdLogger{}
}

func (l *stdLogger) Info(msg string) {
	return
}

func (l *stdLogger) Warn(msg string) {
	return
}

func (l *stdLogger) Error(msg string) {
	return
}

func getArguments() (int, string, bool) {
	// Get port to be used
	PORT := flag.Int("port", 8080, "Port to be used for communication")

	// Defining a string flag
	PASSWD := flag.String("password", "", "Password to be used for authentication of the agents")

	// Define if program is to be run on silent mode
	SILENT := flag.Bool("silent", false, "Remove prompt from startup")

	// Call flag.Parse() to parse the command-line flags
	flag.Parse()

	return *PORT, *PASSWD, *SILENT
}

func getPrompt() {
	version := "1.0"
	prompt := " (                          *     \n )\\ )         (           (  `    \n(()/(         )\\          )\\))(   \n /(_))     ((((_)(       ((_)()\\  \n(_))        )\\ _ )\\      (_()((_) \n| _ \\       (_)_\\(_)     |  \\/  | \n|   /   _    / _ \\    _  | |\\/| | \n|_|_\\  (_)  /_/ \\_\\  (_) |_|  |_| "

	fmt.Println(prompt)
	fmt.Printf("\n[Remote Anamnestic Mapper (v%s)]\n\n", version)
}

func addAgent() {
	fmt.Println("[+] Listing Agents.")
	receiveMessage()
}

func receiveMessage() {
	tcp := p2p.NewTCP("localhost", "8080")

	server, err := p2p.NewServer(tcp)
	if err != nil {
		log.Panicln(err)
	}

	server.SetLogger(NewStdLogger())

	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		message := Message{}
		err = req.GetGob(&message)
		if err != nil {
			return
		}

		fmt.Printf("[+] Added agent: %s\n", message.Text)

		return
	})

	err = server.Serve()
	if err != nil {
		log.Panicln(err)
	}
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

func requestRAM(port int, password string) {
	fmt.Println("[+] Requesting RAM.")

	// Check vars given
	fmt.Println("[*] Port being used: ", port)
	if password != "" {
		fmt.Println("[*] Password being used: ", password)
	} else {
		fmt.Println("[!] WARNING! The server is being run without a password!")
	}

}

func menu(port int, password string) {
	var inpt string
	for true {
		fmt.Println("[?] Choose one of the options:\n       1 - Add Agent\n       2 - Request Agent's RAM\n       3 - Quit")
		fmt.Print("> ")

		fmt.Scanln(&inpt)

		switch inpt {
		case "1":
			addAgent()
		case "2":
			requestRAM(port, password)
		case "3":
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
