package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatModel struct {
	messages []string
	input    string
	conn     net.Conn
	isClient bool
}

// New initializes a new Model instance
func New() ChatModel {
	return ChatModel{
		messages: []string{"Welcome to the chat!"},
		input:    "",
		isClient: false,
	}
}

// Init is required for Bubble Tea but isn't used here
func (m ChatModel) Init() tea.Cmd {
	go startServer()
	return nil
}

// Update handles user input (not needed for chat functionality in this case)
func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.input == "" {
				return m, nil
			}
			if m.input == "exit" {
				return m, tea.Quit
			}
			if m.conn != nil {
				sendMessage(m.conn, m.input)
			}

			m.messages = append(m.messages, "You: "+m.input)
			m.input = ""

		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

// View renders the application UI
func (m ChatModel) View() string {
	var chatView string
	for _, msg := range m.messages {
		chatView += msg + "\n"
	}
	return fmt.Sprintf(
		"Chat:\n%s\n\nType and press Enter to send.\n(Type 'exit' to quit)\n %s",
		chatView, m.input,
	)
}

// RunChat starts the server and client for messaging
func RunChat() {
	go startServer() // Run server in a goroutine
	startClient()    // Run client in the main thread
}

// Start a TCP server to listen for incoming messages
func startServer() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the port to listen on (default 8080):")
	scanner.Scan()
	port := scanner.Text()
	if port == "" {
		port = "8080" // Default port
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Println("Someone has connected")
		go handleConnection(conn)
	}
}

// Handles incoming messages from a connected client
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Friend disconnected:", err)
			return
		}

		var receivedMsg Model
		err = json.Unmarshal([]byte(message), &receivedMsg)
		if err != nil {
			log.Println("Invalid message received:", message)
			continue
		}

		log.Printf("[%s] %s: %s\n", receivedMsg.Timestamp.Format("15:04"), receivedMsg.Sender, receivedMsg.Text)
	}
}

func startClient() {
	log.Println("Starting client...")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Enter the peer's address (e.g., 127.0.0.1:8081):")
		scanner.Scan()
		address := scanner.Text()

		if !strings.Contains(address, ":") {
			fmt.Println("Invalid address. Format: hostname:port")
			continue
		}

		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Printf("Error connecting to peer: %v", err)
			fmt.Println("Failed to connect. Try again.")
			continue
		}
		defer conn.Close()

		fmt.Println("Connected! Type messages (type 'exit' to quit):")
		go handleConnection(conn)

		for scanner.Scan() {
			text := scanner.Text()
			if text == "exit" {
				log.Println("Exiting client...")
				return
			}
			sendMessage(conn, text) // Use sendMessage function
		}
		break
	}
}

func sendMessage(conn net.Conn, text string) {
	msg := Model{
		Text:      text,
		Sender:    "User1", // Replace with actual sender ID
		Receiver:  "User2", // Replace with actual receiver ID
		Type:      "chat",
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	writer := bufio.NewWriter(conn)
	_, err = writer.Write(append(jsonData, '\n'))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	writer.Flush() // Ensure message is sent immediately
}
