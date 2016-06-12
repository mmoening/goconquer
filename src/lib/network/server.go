// Conquer Online uses the TCP/IPv4 protocol and the client-server distributed 
// application structure to establish a connection between the player and 
// multiple processing servers. This package provides the server projects with 
// portable interfaces for network management and network data structuring.
package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"lib/structures"
	"net"
)

// Server contains function pointers to events which the server calls to process
// client requests. These function pointers should be initialized upon creating
// the structure. Not all events are required to be defined.
type Server struct {
	OnConnect    func(*structures.Client)
	OnExchange   func(*structures.Client, *bytes.Buffer)
	OnReceive    func(*structures.Client, *bytes.Buffer)
	OnDisconnect func(*structures.Client)
}

// Listen binds the new server to a hostname and port, which is a network
// interface and port on the local machine; defaulted to localhost. It then
// immediately starts the server and accepts new connections from clients on a 
// separate go routine (non-blocking). It returns true on the channel if the 
// listen was successful and the server terminates without error.
func (s Server) Listen(host string, ch chan bool) {

	// Listen for incoming connections from the game client.
	listener, err := net.Listen("tcp", host)
	if err != nil { fmt.Println(err.Error()); ch <- false; return }
	defer listener.Close()	
	
	for { // While the application is running, accept new connections.
		connection, err := listener.Accept()
		if err != nil { fmt.Println(err.Error())
		} else { go s.Accept(connection) }
	}
	ch <- true
}

// Accept is called by the listener's go routine, created when the server accepts
// a new connection from the remote client. It receives data from the client's
// remote descriptor as long as the client is connected by calling into the 
// server's receive function. 
func (s Server) Accept(connection net.Conn) {
	defer connection.Close()
	client := &structures.Client { Connection: connection }
	if s.OnConnect != nil { s.OnConnect(client) }
	s.Receive(client)
}

// Receive is called by the Accept function to receive data from the client. If 
// the OnExchange event is defined, a buffer will be created for receiving data
// for the initial key exchange; else, the event will be skipped and go straight 
// to the receive loop. In the receive loop, packets are received by reading in 
// the header and then reading in the body using the expected packet length. 
// Since golang buffers data automatically, this shouldn't be any more costly 
// than handling packet splitting and the client's packet fragmentation using 
// pointer arithmetic and buffer persistence. 
func (s Server) Receive(client *structures.Client) {
	defer s.Disconnect(client)
	if s.OnExchange != nil {
		
		// Create the buffer and receive for an unknown length.
		buffer := make([]byte, 4096)
		length, err := client.Connection.Read(buffer)
		if err != nil || length == 0 { return }
		client.Cipher.Decrypt(buffer)
		packet := bytes.NewBuffer(buffer)
		s.OnExchange(client, packet)
	}
	if s.OnReceive == nil { return }
	for {
		// The first two bytes contains the expected length.
		buffer := make([]byte, 2)
		length, err := io.ReadFull(client.Connection, buffer)
		if err != nil || length != 2 { return }
		client.Cipher.Decrypt(buffer)
		
		// Sanity check against expected length.
		length = int(binary.LittleEndian.Uint16(buffer[0:2]))
		if (length < 4 || length > 4096) { 
			fmt.Println("invalid packet length") 
			return 
		}
		
		// Read the remaining bytes of the packet.
		packet := make([]byte, length)
		fulllength, err := io.ReadFull(client.Connection, packet[2:length])
		if err != nil || fulllength != length - 2 { return }
		
		// Combine the buffers and process.
		binary.LittleEndian.PutUint16(packet[0:2], uint16(length))
		client.Cipher.Decrypt(packet[2:length])
		p := bytes.NewBuffer(packet)
		s.OnReceive(client, p)
	}
}

// Disconnect is called by the defer in the receive function once an error 
// occurs or the client has disconnected from the server. This calls the 
// disconnect event for processing client features upon disconnect (such as 
// discontinuing trade transactions, removing the character from the map, etc).
func (s Server) Disconnect(client *structures.Client) {
	if s.OnDisconnect != nil { s.OnDisconnect(client) }
}
