package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Create a type that represents the command type
type commandType uint8

// Create an enumeration of the different types of commands that can be sent
const (
	// Command type for sending a message
	ping commandType = iota
	// Command type for sending close
	pong
	// Command type for sending a ping
	dataMessage
	// Command type for sending a pong
	close
)

// Message is a struct that represents a message that can be sent over the network
type Message struct {
	// Type is the type of message that is being sent
	Type commandType
	// The length of the data that is being sent
	Length uint16
	// Data is the data that is being sent
	Data []byte
}

// String will return a string representation of the message
func (m Message) String() string {
	return fmt.Sprintf("Type: %d, Data: %s", m.Type, string(m.Data))
}

func (m Message) MarshallBinary() ([]byte, error) {
	// Create a buffer to write the data to
	var buf bytes.Buffer

	// Write the type to the buffer
	err := binary.Write(&buf, binary.BigEndian, m.Type)
	if err != nil {
		return nil, err
	}

	// Write the length of the data to the buffer
	err = binary.Write(&buf, binary.BigEndian, m.Length)
	if err != nil {
		return nil, err
	}

	// Write the data to the buffer
	_, err = buf.Write(m.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Message) UnmarshallBinary(data []byte) error {
	// Create a buffer to read the data from
	buf := bytes.NewBuffer(data)

	// Read the type from the buffer
	err := binary.Read(buf, binary.BigEndian, &m.Type)
	if err != nil {
		return err
	}

	// Read the length of the data from the buffer
	err = binary.Read(buf, binary.BigEndian, &m.Length)
	if err != nil {
		return err
	}

	// Read the data from the buffer
	m.Data = make([]byte, m.Length)
	_, err = buf.Read(m.Data)
	if err != nil {
		return err
	}

	return nil
}

// main will parse the command line arguments to determine if this is a client or server
func main() {
	// Parse the command line to work out if this is a client or server
	isServer := flag.Bool("s", false, "Run as server")

	flag.Parse()

	if *isServer {
		fmt.Println("Running as server")

		l, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}

		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}

			go handleConnection(conn, isServer)
		}
	} else {
		fmt.Println("Running as client")

		// Connect to the server
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Error dialing:", err.Error())
			os.Exit(1)
		}

		handleConnection(conn, isServer)
	}
}

// handleConnection will read the data from the connection and print it to the console
// It will also write data from the console to the connection
func handleConnection(conn net.Conn, isServer *bool) {
	// Create a channel to handle SIGINT or SIGTERM
	signals := make(chan os.Signal, 1)

	// Register the channel to receive the signals
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Close the connection when the function returns
	defer conn.Close()

	// Create a channel to wait for the connection to close
	done := make(chan struct{})

	// Start a goroutine to read from the connection
	go func() {
		receiver(conn)
		done <- struct{}{}
	}()

	if *isServer {
		// Start a goroutine to write to the connection
		go func() {
			sender(conn)
			done <- struct{}{}
		}()
	}

	// Wait for the connection to close or a signal to be received
	select {
	case <-done:
		fmt.Println("Connection closed")
	case sig := <- signals:
		fmt.Println("Received signal:", sig)
	}
}

// sender will read from the console and write to the connection in a loop
func sender(conn net.Conn) {
	// Send a ping message once a second
	for {
		sendMessage(conn, ping, []byte("Hello, World!"))

		// Sleep for a second
		time.Sleep(time.Second)
	}
}

// receiver will read from the connection and write to the console in a loop
func receiver(conn net.Conn) {
	for {
		// Read a message from the connection
		message, err := readMessage(conn)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// Print the message to the console
		fmt.Println("Received message:", message)

		// Switch on the type of message
		switch message.Type {
		case ping:
			// Send a pong message
			sendMessage(conn, pong, message.Data)
		case dataMessage:
			// Send a close message
			sendMessage(conn, close, message.Data)
			return
		case pong:
			// Do nothing
		case close:
			// Close the connection
			return
		}
	}
}

// sendMessage will send a message over the connection
func sendMessage(conn net.Conn, messageType commandType, data []byte) {
	// Create a message to send
	message := Message{
		Type:   messageType,
		Length: uint16(len(data)),
		Data:   data,
	}

	// Print the message to the console
	fmt.Println("Sending message:", message)

	// Marshall the message into a byte array
	messageData, err := message.MarshallBinary()
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	// Write the message to the connection
	_, err = conn.Write(messageData)
	if err != nil {
		fmt.Println("Error writing message:", err)
		return
	}
}

// readMessage will read a message from the connection
func readMessage(conn net.Conn) (Message, error) {
	// Create a buffer to read the data into
	buf := make([]byte, 1024)

	// Read the data from the connection
	n, err := conn.Read(buf)
	if err != nil {
		return Message{}, err
	}

	// Unmarshall the data into a message
	message := Message{}
	err = message.UnmarshallBinary(buf[:n])
	if err != nil {
		return Message{}, err
	}

	return message, nil
}
