package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
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
	messages      []string
	input         string
	conn          net.Conn
	isClient      bool
	serverStarted bool
}

// New initializes a new ChatModel instance
func New() *ChatModel {
	return &ChatModel{
		messages:      []string{},
		input:         "",
		isClient:      false,
		serverStarted: false,
	}
}

// Init initializes the model and optionally starts the server
func (m *ChatModel) Init() tea.Cmd {
	// Do nothing initially, no goroutines
	return nil
}

// Update handles user input and sends messages
func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.input == "start chatting" && !m.serverStarted {
				m.serverStarted = true
				// Start the server directly without using goroutines
				startServer()
			}
			if m.input == "exit" {
				return m, tea.Quit
			}
			if m.conn != nil && m.input != "" {
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

// View renders the chat interface
func (m *ChatModel) View() string {
	var chatView string
	for _, msg := range m.messages {
		chatView += msg + "\n"
	}

	return fmt.Sprintf("Chat:\n%s\n\nType and press Enter to send.\n(Type 'exit' to quit)\n %s", chatView, m.input)
}

// RunChat starts the server and client for messaging
func RunChat(p *tea.Program) {
	// Send a message to the Bubble Tea program to simulate a key press (start chatting)
	p.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Simulate pressing Enter
}

// startServer starts the server and waits for incoming connections
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
		handleConnection(conn)
	}
}

// handleConnection handles the incoming messages from the client
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

// sendMessage sends a message to the connected peer
func sendMessage(conn net.Conn, text string) {
	msg := Model{
		Text:      text,
		Sender:    "User1",
		Receiver:  "User2",
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

	writer.Flush()
}
