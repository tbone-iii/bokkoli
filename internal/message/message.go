package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type peerConn struct {
	conn net.Conn
}

type listenerConn struct {
	conn net.Conn
}

type errorOnMessageReceive struct {
	err error
}
type errorOnMessageSend struct {
	err error
}

type (
	direction    string
	listener     net.Listener
	incomingJson []byte
)

const (
	Outgoing direction = "outgoing"
	Incoming direction = "incoming"
)

type Message struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Direction direction `json:"direction"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatModel struct {
	// TODO: If # of fields increase, break down into grouped sub-structs
	messages     []Message
	input        string
	peerConn     net.Conn // TODO: You will implement an array of peers
	listenerConn net.Conn
	listener     net.Listener
	isClient     bool
	portNumber   string
	chatDb       *ChatDb
	username     string
}

func New() *ChatModel {
	db, err := NewChatDB(DefaultDbFilePath)
	if err != nil {
		log.Println("Error upon DB creation: ", err)
	}

	return &ChatModel{
		messages:   []Message{},
		input:      "",
		isClient:   false,
		portNumber: "8080",
		username:   "Somy", // TODO: Pull this info from DB or user input [look at 'huh' bubbletea library]
		chatDb:     db,
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
				defer m.peerConn.Close()
				defer m.listenerConn.Close()
				return m, tea.Quit
			}

			// Need to connect to a client still
			suffix, success = parseStringSuffixFromPrefix(m.input, "connect to port ")
			if success && m.peerConn == nil {
				m.input = ""
				return m, createPeerConnCmd("", suffix)
			}

			// Send messages command
			if m.input != "" && m.peerConn != nil {
				temp_input := m.input
				m.input = ""
				message := createMessage(temp_input, m.username, "Pickle", Outgoing)
				return m, handleDbAndSendMessageCmd(message, m.peerConn, m.chatDb)
			}

			log.Printf("Not all cases have been handled. There is an issue here.")
		case "backspace":
			m.input = deleteLastNCharacters(m.input, 1)
			log.Printf("Not all cases have been handled. There is an issue here.")
		case "backspace":
			m.input = deleteLastNCharacters(m.input, 1)
		default:
			m.input += msg.String()
		}
	case Message:
		if msg.Direction == Outgoing {
			m.messages = append(m.messages, msg)
		} else if msg.Direction == Incoming {
			m.messages = append(m.messages, msg)
			return m, handleListenerConnCmd(m.listenerConn)
		} else {
			log.Fatal("There should not be any other directions. Crashing program.")
		}
	case peerConn:
		m.peerConn = msg.conn // TODO: Make into an array appending
	case listenerConn:
		log.Println("Connection read from listener on port: ", m.portNumber)
		m.listenerConn = msg.conn // TODO: Make into an array appending
		return m, handleListenerConnCmd(m.listenerConn)
	case incomingJson:
		return m, handleDbAndReceiveMessageCmd(msg, m.chatDb)
	case errorOnMessageSend:
		// Do something here based on that
	case errorOnMessageReceive:
		// Do something here based on that
	case listener:
		log.Println("Listener started on port: ", m.portNumber)
		m.listener = msg
		// Read new connections
		return m, readListenerCmd(m.listener)
	}

	return m, nil
}

func (m *ChatModel) View() string {
	// TODO: Consider asking for port number in a separate model/view
	var chatView string = "Type 'start chat my port ' and the port number for your server to join the chatroom.\n" +
		"Then type 'connect to port ' and follow it with a port number. (Type 'exit' to quit.)\n\n"

	for _, message := range m.messages {
		chatView += fmt.Sprintf("%s\n%s: %s\n-----------------\n", message.Timestamp.Format("2025-02-05 13:08"), message.Sender, message.Text)
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

func startServer(port string) listener {
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
		jsonMessage, err := handleListenerConn(conn)
		if err != nil {
			return errorOnMessageReceive{err: err}
		}

		return jsonMessage
	}
}

func handleListenerConn(conn net.Conn) (incomingJson, error) {
	reader := bufio.NewReader(conn)

	for {
		log.Println("Listener is handling connection, awaiting read: ", conn.LocalAddr().String())
		jsonMessage, err := reader.ReadBytes('\n')

		if err != nil {
			log.Println("Friend disconnected:", err)
			return jsonMessage, err
		}

		log.Println("Handle listener message received as: ", jsonMessage)
		return jsonMessage, nil
	}
}

func createMessage(text, sender, receiver string, direction direction) Message {
	return Message{
		Text:      text,
		Sender:    sender,
		Receiver:  receiver,
		Direction: direction,
		Timestamp: time.Now(),
	}
}

func serializeMessage(message Message) ([]byte, error) {
	jsonData, err := json.Marshal(message)
	return jsonData, err
}

func deserializeJsonMessage(jsonData []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(jsonData, &msg)
	return msg, err
}

func handleDbAndReceiveMessage(jsonData []byte, db *ChatDb) (Message, error) {
	msg, err := deserializeJsonMessage(jsonData)
	if err != nil {
		log.Println("Error deserializing JSON message: ", err)
		return msg, err
	}

	msg.Direction = Incoming

	err = db.saveMessage(msg)
	if err != nil {
		log.Println("Error saving message to DB: ", err)
		return msg, err
	}

	return msg, nil
}

func handleDbAndSendMessage(msg Message, conn net.Conn, db *ChatDb) (Message, error) {
	err := db.saveMessage(msg)
	if err != nil {
		log.Println("Error saving message to DB: ", err)
	}

	jsonData, err := serializeMessage(msg)
	if err != nil {
		log.Println("Error marshalling message: ", err)
		return msg, err
	}

	_, err = conn.Write(append(jsonData, '\n')) // Add new line for reader to actually parse the delimiter appropriately
	if err != nil {
		log.Printf("error sending messageL %v", err)
		return msg, err
	}

	return msg, nil
}

func handleDbAndSendMessageCmd(msg Message, conn net.Conn, db *ChatDb) tea.Cmd {
	return func() tea.Msg {
		msg, err := handleDbAndSendMessage(msg, conn, db)
		if err != nil {
			return errorOnMessageSend{err: err}
		}

		return msg
	}
}

func handleDbAndReceiveMessageCmd(jsonData []byte, db *ChatDb) tea.Cmd {
	return func() tea.Msg {
		msg, err := handleDbAndReceiveMessage(jsonData, db)
		if err != nil {
			return errorOnMessageReceive{err: err}
		}

		return msg
	}
}
