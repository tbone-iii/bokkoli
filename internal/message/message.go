package message

import (
	"bokkoli/internal/db"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var timestampStyle = lipgloss.NewStyle().Italic(true)

var senderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

var messageStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	MaxWidth(30)
var inputStyle = lipgloss.NewStyle().Faint(true)
var inputLineIndicator = lipgloss.NewStyle().
	Blink(true).
	Foreground(lipgloss.Color("#1379af"))

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

type Setup struct {
	Username string
	Port     string
}

type (
	listener     net.Listener
	incomingJson []byte
)

type ChatModel struct {
	// TODO: If # of fields increase, break down into grouped sub-structs
	messages     []db.Message
	input        string
	peerConn     net.Conn // TODO: You will implement an array of peers
	listenerConn net.Conn
	listener     net.Listener
	isClient     bool
	portNumber   string
	dbHandler    *db.DbHandler
	username     string
}

func New() *ChatModel {
	dbHandler, err := db.NewDbHandler(db.DefaultDbFilePath)
	if err != nil {
		log.Fatal("Error upon DB creation: ", err)
	}

	err = dbHandler.SetupSchemas()
	if err != nil {
		log.Fatal("Error upon schema creation: ", err)
	}

	var portNumber string //= "8080"
	var username string   //= "DefaultUser"

	setup, err := dbHandler.ReadSetup()
	if err == nil {
		portNumber = setup.Port
		username = setup.Username
		log.Println("Setting port number and username read from DB: ", portNumber, username)
	}

	return &ChatModel{
		messages:   []db.Message{},
		input:      "",
		isClient:   false,
		portNumber: portNumber,
		username:   username,
		dbHandler:  dbHandler,
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
				message := createMessage(temp_input, m.username, "Pickle", db.Outgoing)
				return m, handleDbAndSendMessageCmd(message, m.peerConn, m.dbHandler)
			}

			log.Printf("Not all cases have been handled. There is an issue here.")
		case "backspace":
			m.input = deleteLastNCharacters(m.input, 1)
		default:
			m.input += msg.String()
		}
	case db.Message:
		switch msg.Direction {
		case db.Outgoing:
			m.messages = append(m.messages, msg)
		case db.Incoming:
			m.messages = append(m.messages, msg)
			return m, handleListenerConnCmd(m.listenerConn)
		default:
			log.Fatal("There should not be any other directions. Crashing program.")
		}
	case peerConn:
		m.peerConn = msg.conn // TODO: Make into an array appending
	case listenerConn:
		log.Println("Connection read from listener on port: ", m.portNumber)
		m.listenerConn = msg.conn // TODO: Make into an array appending
		return m, handleListenerConnCmd(m.listenerConn)
	case incomingJson:
		return m, handleDbAndReceiveMessageCmd(msg, m.dbHandler)
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

	var chatView strings.Builder
	chatView.WriteString("Type 'start chat my port ' and the port number for your server to join the chatroom.\n")
	chatView.WriteString("Then type 'connect to port ' and follow it with a port number. (Type 'exit' to quit.)\n\n")

	// TODO: If the time is the same and user is the same, join all the messages into one block
	// TODO: If the time is different and the user is the same, show only the first time
	// TODO: Perhaps make the last message blink?
	for _, message := range m.messages {
		tempChatView := fmt.Sprintf(
			"%s\n%s: %s",
			timestampStyle.Render(message.Timestamp.Format("2006-01-02 15:04")),
			senderStyle.Render(message.Sender),
			message.Text,
		)
		chatView.WriteString(messageStyle.Render(tempChatView) + "\n")
	}

	indicator := inputLineIndicator.Render("> ")
	return fmt.Sprintf("%s\n\n%s", chatView.String(), inputStyle.Render(indicator, m.input))
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

// func handleListenerConnCmd(conn net.Conn) tea.Cmd {
// 	return func() tea.Msg {
// 		jsonMessage, err := handleListenerConn(conn)
// 		if err != nil {
// 			return errorOnMessageReceive{err: err}
// 		}

//			return jsonMessage
//		}
//	}
func handleListenerConnCmd(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		if conn == nil {
			return errorOnMessageReceive{err: fmt.Errorf("connection is nil")}
		}

		jsonMessage, err := handleListenerConn(conn)
		if err != nil {
			return errorOnMessageReceive{err: err}
		}

		return jsonMessage
	}
}

// func handleListenerConn(conn net.Conn) (incomingJson, error) {
// 	reader := bufio.NewReader(conn)

// 	for {
// 		log.Println("Listener is handling connection, awaiting read: ", conn.LocalAddr().String())
// 		jsonMessage, err := reader.ReadBytes('\n')

// 		if err != nil {
// 			log.Println("Friend disconnected:", err)
// 			return jsonMessage, err
// 		}

//			log.Println("Handle listener message received as: ", jsonMessage)
//			return jsonMessage, nil
//		}
//	}
func handleListenerConn(conn net.Conn) (incomingJson, error) {
	reader := bufio.NewReader(conn)

	log.Println("Listener is handling connection, awaiting read: ", conn.LocalAddr().String())
	jsonMessage, err := reader.ReadBytes('\n')

	if err != nil {
		log.Println("Friend disconnected:", err)
		return nil, err
	}

	log.Println("Handle listener message received as: ", jsonMessage)
	return jsonMessage, nil
}

func createMessage(text, sender, receiver string, direction db.Direction) db.Message {
	return db.Message{
		Text:      text,
		Sender:    sender,
		Receiver:  receiver,
		Direction: direction,
		Timestamp: time.Now(),
	}
}

func serializeMessage(message db.Message) ([]byte, error) {
	jsonData, err := json.Marshal(message)
	return jsonData, err
}

func deserializeJsonMessage(jsonData []byte) (db.Message, error) {
	var message db.Message
	err := json.Unmarshal(jsonData, &message)
	return message, err
}

func handleDbAndReceiveMessage(jsonData []byte, dbHandler *db.DbHandler) (db.Message, error) {
	message, err := deserializeJsonMessage(jsonData)
	if err != nil {
		log.Println("Error deserializing JSON message: ", err)
		return message, err
	}

	message.Direction = db.Incoming

	err = dbHandler.SaveMessage(message)
	if err != nil {
		log.Println("Error saving message to DB: ", err)
		return message, err
	}

	return message, nil
}

func handleDbAndSendMessage(message db.Message, conn net.Conn, dbHandler *db.DbHandler) (db.Message, error) {
	err := dbHandler.SaveMessage(message)
	if err != nil {
		log.Println("Error saving message to DB: ", err)
	}

	jsonData, err := serializeMessage(message)
	if err != nil {
		log.Println("Error marshalling message: ", err)
		return message, err
	}

	_, err = conn.Write(append(jsonData, '\n')) // Add new line for reader to actually parse the delimiter appropriately
	if err != nil {
		log.Printf("error sending messageL %v", err)
		return message, err
	}

	return message, nil
}

func handleDbAndSendMessageCmd(message db.Message, conn net.Conn, dbHandler *db.DbHandler) tea.Cmd {
	return func() tea.Msg {
		msg, err := handleDbAndSendMessage(message, conn, dbHandler)
		if err != nil {
			return errorOnMessageSend{err: err}
		}

		return msg
	}
}

func handleDbAndReceiveMessageCmd(jsonData []byte, dbHandler *db.DbHandler) tea.Cmd {
	return func() tea.Msg {
		msg, err := handleDbAndReceiveMessage(jsonData, dbHandler)
		if err != nil {
			return errorOnMessageReceive{err: err}
		}

		return msg
	}
}
