package message

import (
	"bokkoli/internal/db"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	LOWERBOUND_PORT_NUMBER int = 1024
	UPPERBOUND_PORT_NUMBER int = 49151
)

var timestampStyle = lipgloss.NewStyle().Italic(true).
	Faint(true)

var senderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#a488f7"))

var timestampSenderStyle = lipgloss.NewStyle().
	PaddingBottom(2)

var maxLineLength = 50

var messageStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	MaxWidth(maxLineLength).
	Padding(1, 1, 1)

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

type (
	listener     net.Listener
	incomingJson []byte
)

type ChatModel struct {
	messages     []db.Message
	input        string
	peerConn     net.Conn // TODO: implement an array of peers
	listenerConn net.Conn
	listener     net.Listener
	isClient     bool
	dbHandler    *db.DbHandler
	settings     *db.Setup
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

	settings, err := dbHandler.ReadSetup()
	if err == nil {
		log.Println("User connection settings read from DB: ", settings)
	}

	return &ChatModel{
		messages:  []db.Message{},
		input:     "",
		isClient:  false,
		settings:  &settings,
		dbHandler: dbHandler,
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
					suffix = m.settings.Port
				}
				m.settings.Port = suffix
				return m, startListenerCmd(suffix)
			}

			if m.input == "exit" || m.input == "ctrl+c" {
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
				message := createMessage(temp_input, m.settings.Username, db.Outgoing)
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
		log.Println("Connection read from listener on port: ", m.settings.Port)
		m.listenerConn = msg.conn // TODO: Make into an array appending
		return m, handleListenerConnCmd(m.listenerConn)
	case incomingJson:
		return m, handleDbAndReceiveMessageCmd(msg, m.dbHandler)
	case errorOnMessageSend:
		// Do something here based on that
	case errorOnMessageReceive:
		// Do something here based on that
	case listener:
		log.Println("Listener started on port: ", m.settings.Port)
		m.listener = msg
		// Read new connections
		return m, readListenerCmd(m.listener)
	}

	return m, nil
}

func (m *ChatModel) View() string {
	// TODO: Consider asking for port number in a separate model/view

	var chatView strings.Builder

	chatView.WriteString(fmt.Sprintf("\n*** To start a chat server type %s followed by your port number. \n (If no port number is provided, the default entered in %s will be used.)\n\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f47d56")).Bold(true).Render("'start chat my port <port>'"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f47d56")).Bold(true).Render("'user settings'"),
	))
	chatView.WriteString(fmt.Sprintf("*** To connect to a chat server type %s followed by the port number.",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f47d56")).Bold(true).Render("'connect to port <port>'"),
	))

	chatView.WriteString(fmt.Sprintf("\n\n%s.\n\n",
		lipgloss.NewStyle().Faint(true).Render("Press 'esc' to return to main menu.\nTo exit, type 'exit' or press 'ctrl + c' to exit program"),
	))

	for _, message := range m.messages {
		tempChatView := fmt.Sprintf("%s - %s", timestampStyle.Render(message.Timestamp.Format("2006-01-02 15:04")), senderStyle.Render(message.Sender))
		tempChatView = timestampSenderStyle.Render(tempChatView) + "\n"
		tempChatView += wrapText(message.Text, maxLineLength-messageStyle.GetHorizontalPadding())
		chatView.WriteString(messageStyle.Render(tempChatView) + "\n")
	}

	indicator := inputLineIndicator.Render("> ")
	return fmt.Sprintf("%s\n\n%s", chatView.String(), inputStyle.Render(indicator, m.input))
}

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
	if !validatePort(port) {
		return func() tea.Msg {
			fmt.Printf("Only ports in the range %d-%d are allowed, re-enter", LOWERBOUND_PORT_NUMBER, UPPERBOUND_PORT_NUMBER)
			return errorOnMessageReceive{err: fmt.Errorf("invalid port: %s", port)}
		}
	}

	listener, err := startServer(port)
	if err != nil {
		return func() tea.Msg { return errorOnMessageReceive{err: fmt.Errorf("port already in use: %s", port)} }
	}

	return func() tea.Msg {
		return listener
	}
}

func startServer(port string) (listener, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting server, port '%s' is already in use, select a different port. ", port)
		return listener, err
	}
	fmt.Printf("Server listening on port '%s'\n", port)

	return listener, nil
}

func createPeerConnCmd(address string, portNumber string) tea.Cmd {
	return func() tea.Msg {
		fullAddress := net.JoinHostPort(address, portNumber)
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

func createMessage(text string, sender string, direction db.Direction) db.Message {
	return db.Message{
		Text:      text,
		Sender:    sender,
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

	// The newline enables reader to actually parse the delimiter appropriately
	_, err = conn.Write(append(jsonData, '\n'))

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

func validatePort(port string) bool {
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("%s is not valid port number due to not being integer:", port)
		return false
	}
	if !(LOWERBOUND_PORT_NUMBER <= portNumber && portNumber <= UPPERBOUND_PORT_NUMBER) {
		log.Printf("%d is out of range [%d, %d]", portNumber, LOWERBOUND_PORT_NUMBER, UPPERBOUND_PORT_NUMBER)
		return false
	}
	return true
}

// Intelligently wrap text based on max line length and breaks on spaces
// If a word is longer than the max line length, it will break the word into parts
func wrapText(text string, maxLineLength int) string {
	words := strings.Fields(text)

	const LONG_WORD_PADDING int = 2
	var result strings.Builder
	var line string

	for _, word := range words {
		if len(word) > maxLineLength {
			if len(line) > 0 {
				if len(line) < maxLineLength {
					line += " "
				}
				result.WriteString(line)
				result.WriteRune('\n')
				line = ""
			}
			for len(word) > maxLineLength {
				result.WriteString(word[:maxLineLength-LONG_WORD_PADDING])
				result.WriteRune('\n')
				word = word[maxLineLength-LONG_WORD_PADDING:]
			}
			if len(word) > 0 {
				line = word
			}
			continue
		}

		if len(line) == 0 {
			line = word
		} else if len(line)+1+len(word) <= maxLineLength {
			line += " " + word
		} else {
			result.WriteString(line)
			result.WriteRune('\n')
			line = word
		}
	}
	if len(line) > 0 {
		result.WriteString(line)
	}
	return result.String()
}
