package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type listenerType struct {
	net.Listener
}

// type rawReceivedMessage struct {
// 	Text      string
// 	Sender    string
// 	Receiver  string
// 	Timestamp time.Time
// }

// type RawSentMessage struct {
// 	Text      string
// 	Sender    string
// 	Receiver  string
// 	Timestamp time.Time
// }

type peerConn struct {
	conn net.Conn
}

type listenerConn struct {
	conn net.Conn
}

type Message struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatModel struct {
	messages     []string
	input        string
	peerConn     net.Conn // TODO: You will implement an array of peers
	listenerConn net.Conn
	listener     net.Listener
	isClient     bool
	portNumber   string
	chatDB       ChatDB
}

func New() *ChatModel {
	return &ChatModel{
		messages:   []string{},
		input:      "",
		isClient:   false,
		portNumber: "8080",
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return nil
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Start and return listener
			suffix, success := parseStringSuffixFromPrefix(m.input, "start chat my port ")
			if success && m.listener == nil {
				m.input = ""
				if suffix == "" {
					suffix = m.portNumber
				}
				m.portNumber = suffix
				return m, startListenerCmd(suffix)
			}

			if m.input == "exit" {
				defer m.peerConn.Close()
				defer m.listenerConn.Close()
				return m, tea.Quit
			}

			// Need to connect to a client still, returns peerConnType
			suffix, success = parseStringSuffixFromPrefix(m.input, "connect to port ")
			if success && m.peerConn == nil {
				m.input = ""
				return m, createPeerConnCmd("", suffix)
			}

			// Send messages command returns rawSentMessage
			if m.input != "" && m.peerConn != nil {
				temp_input := m.input
				m.input = ""
				return m, sendMessageCmd(m.peerConn, temp_input)
			}

			log.Printf("Not all cases have been handled. There is an issue here.")
		case "backspace":
			m.input = deleteLastNCharacters(m.input, 1)
		default:
			m.input += msg.String()
		}
	case rawReceivedMessage:
		msg := msg
		m.messages = append(m.messages, msg.Text)
		return m, handleListenerConnCmd(m.listenerConn)
	case rawSentMessage:
		m.messages = append(m.messages, msg.Text)
	case peerConn:
		m.peerConn = msg.conn // TODO: Make into an array appending
	case listenerConn:
		log.Println("Connection read from listener on port: ", m.portNumber)
		m.listenerConn = msg.conn // TODO: Make into an array appending
		return m, handleListenerConnCmd(m.listenerConn)
	case listenerType:
		log.Println("Listener started on port: ", m.portNumber)
		m.listener = msg
		// Read new connections
		return m, readListenerCmd(m.listener)
	}

	return m, nil
}

func (m *ChatModel) View() string {
	// TODO: Consider asking for port number in a separate model/view
	var chatView string = "Type 'start chat my port ' and the port number for your server to join the chatroom." +
		"\nThen type 'connect to port ' and follow it with a port number. (Type 'exit' to quit.)\n\n"

	for _, msg := range m.messages {
		chatView += string(msg) + "\n----------------------\n"
	}

	return fmt.Sprintf("\n %s\n%s", chatView, m.input)
}

// Removes the last n number of characters from a string and returns the new string.
// 's' is defined as some string
func deleteLastNCharacters(s string, n int) string {
	size := len(s)

	if size-n < 0 {
		return s
	}

	s = s[:size-n]
	return s
}

// Return a True or False on success for whether a prefix was found in the result.
func parseStringSuffixFromPrefix(s string, prefix string) (string, bool) {
	suffix, success := strings.CutPrefix(strings.ToLower(s), prefix)
	if !success {
		return "", false
	}

	return suffix, true
}

func readListenerCmd(listener net.Listener) tea.Cmd {
	return func() tea.Msg {
		return readListener(listener)
	}
}

func readListener(listener net.Listener) listenerConn {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Println("Someone has connected: ", conn.LocalAddr().String())
		return listenerConn{conn: conn}
	}
}

func startListenerCmd(port string) tea.Cmd {
	return func() tea.Msg {
		return startServer(port)
	}
}

func startServer(port string) listenerType {
	for {
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("Error starting server:", err)
			continue
		}
		fmt.Printf("Server listening on port '%s'\n", port)

		return listener
	}
}

func createPeerConnCmd(address string, portNumber string) tea.Cmd {
	return func() tea.Msg {
		fullAddress := fmt.Sprintf("%s:%s", address, portNumber)
		conn, err := net.Dial("tcp", fullAddress)
		if err != nil {
			fmt.Println("error: ", err)
		}

		fmt.Println("Connected to port: ", portNumber)
		return peerConn{conn: conn}
	}
}

func handleListenerConnCmd(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		return handleListenerConn(conn)
	}
}

func handleListenerConn(conn net.Conn) Message {
	reader := bufio.NewReader(conn)

	for {
		log.Println("Listener is handling connection, awaiting read: ", conn.LocalAddr().String())
		message, err := reader.ReadBytes('\n')

		if err != nil {
			log.Println("Friend disconnected:", err)
			// TODO: Should be sending Struct instead of string for more info
			return Message{Text: "Friend disconnected."}
		}

		var receivedMsg Message //make this a struct
		err = json.Unmarshal(message, &receivedMsg)
		if err != nil {
			log.Println("Invalid message received:", message)
			continue
		}

		log.Println("Handle listener message received as: ", receivedMsg.Text)
		// TODO: Should be sending Struct instead of string for more info
		return receivedMsg
	}
}

// func sendMessage(conn net.Conn, text, sender, receiver string) error {
// 	msg := createMessage(text, sender, receiver)
// 	jsonData, err := json.Marshal(msg)
// 	if err != nil {
// 		return fmt.Errorf("error serializing message: %v", err)
// 	}
// }
// func sendMessageCmd(conn net.Conn, text string) tea.Cmd {
// 	return func() tea.Msg {
// 		message := sendMessage(conn, text)
// 		return message
// 	}
// }

// TODO: Refactor to be Struct returned, not raw string
// TODO: Break sendMessage into two functions; one to create a Message struct,
// TODO: another to actually send the JSON parsed message
// TODO: OR rename function to properly capture meaning of what it does
// NOTE: which is storing raw text into a crafted Message struct, serializing, and sending

// Example: prepareAndSendMessage

// Sends a message to connection `conn`, raw `text` message to send
// During this process, converts into a serialized JSON
// Also, stores Message into database
func createMessage(text, sender, receiver string) Message {
	return Message{
		Text:      text,
		Sender:    sender,
		Receiver:  receiver,
		Timestamp: time.Now(),
	}
}

func prepareAndSendMessage(text, sender, receiver string, conn net.Conn) ([]byte, error) {
	msg := createMessage(text, sender, receiver)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error serialize message: %v", err)
	}

	_, err = conn.Write(append(jsonData, '\n'))
	if err != nil {
		return nil, fmt.Errorf("error send message: %v", err)
	}
	return jsonData, nil
}

// func sendMessage(conn net.Conn, text string) rawSentMessage {
// 	msg := Message{
// 		Text:      text,
// 		Sender:    "User1",
// 		Receiver:  "User2",
// 		Timestamp: time.Now(),
// 	}
