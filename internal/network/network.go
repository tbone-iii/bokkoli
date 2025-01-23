package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Peer struct {
	address string
	conn    net.Conn
}

// Main function - program entry point
func main() {
	// Asynchronously start server
	go startServer()

	// os.Stdin is standard input
	scanner := bufio.NewScanner(os.Stdin)
	// Creates variable pointing to Peer struct, starts as nil
	var currentPeer *Peer

	// User interface
	fmt.Println("/connect <ip:port> - Connect to a peer")

	// Input reading
	for scanner.Scan() {
		input := scanner.Text()
		if strings.HasPrefix(input, "/connect ") {
			// Get address
			addr := strings.TrimSpace(strings.TrimPrefix(input, "/connect"))
			// Connect to address
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("Failed to connect to %s: %v\n", addr, err)
				continue
			}
			// Store connection
			currentPeer = &Peer{address: addr, conn: conn}
			fmt.Printf("Connected to %s\n", addr)
			// Start receiving msgs from peer
			go handleConnection(conn)

		} else if strings.HasPrefix(input, "/send ") && currentPeer != nil {
			message := strings.TrimPrefix(input, "/send ")
			_, err := currentPeer.conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			}

		} else if input == "/quit" {
			if currentPeer != nil {
				currentPeer.conn.Close()
			}
			os.Exit(0)
		} else {
			fmt.Println("Invalid command")
		}
	}
}

func startServer() {
	// Listen on IPv4 and IPv6
	listener, err := net.Listen("tcp", ":49999")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		// Exit code for why program terminated
		os.Exit(1)
	}
	// Ensure listener is properly closed when function exits
	defer listener.Close()

	fmt.Printf("Listening on port %s\n", ":49999")

	for {
		// Creates a conn object when a client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

// net.Conn is a type that represents a generic network connection
func handleConnection(conn net.Conn) {
	// defer ensures connection closes when function returns
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Infinite loop to continuously read messages
	for {
		// ReadString reads until it encounters the delimiter '\n'
		message, err := reader.ReadString('\n')
		if err != nil {
			// RemoteAddr returns the remote network address
			fmt.Printf("Connection closed from %s\n", conn.RemoteAddr().String())
			return
		}
		fmt.Printf("\nReceived from %s: %s", conn.RemoteAddr().String(), message)
	}
}
