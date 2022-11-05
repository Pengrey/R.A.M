package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	p2p "github.com/leprosus/golang-p2p"
)

type Hello struct {
	Text string
}

type Buy struct {
	Text string
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

func listAgents() {
	fmt.Println("[+] Listing Agents.")
}

func requestRAM() {
	fmt.Println("[+] Requesting RAM.")
}

func menu() {
	var inpt string
	for true {
		fmt.Println("[?] Choose one of the options:\n       1 - List Agents\n       2 - Request Agent's RAM\n       3 - Quit")
		fmt.Print("> ")

		fmt.Scanln(&inpt)

		switch inpt {
		case "1":
			listAgents()
		case "2":
			requestRAM()
		case "3":
			os.Exit(0)
		default:
			fmt.Printf("[!] Option %s doesnt exist!\n", inpt)
		}
	}
}

func serverStart() {
	tcp := p2p.NewTCP("localhost", "8080")

	server, err := p2p.NewServer(tcp)
	if err != nil {
		log.Panicln(err)
	}

	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		hello := Hello{}
		err = req.GetGob(&hello)
		if err != nil {
			return
		}

		fmt.Printf("> Hello: %s\n", hello.Text)

		res = p2p.Data{}
		err = res.SetGob(Buy{
			Text: hello.Text,
		})

		return
	})

	err = server.Serve()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	// Get arguments from user
	port, password, silent := getArguments()

	// Print Prompt
	if !silent {
		getPrompt()
	}

	// Check vars
	fmt.Println("[*] Port being used: ", port)
	if password != "" {
		fmt.Println("[*] Password being used: ", password)
	} else {
		fmt.Println("[!] WARNING! The server is being run without a password!")
	}

	// Start Menu
	menu()
}
