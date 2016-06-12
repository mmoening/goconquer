package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"lib/packets"
	"lib/security"
)

// Client encapsulates the remote client's endpoint and used throughout the server 
// to send and receive data from the client. The structure is inherited by the 
// server projects' client structure to extend functionality for network actions.
type Client struct {
	Account	   	*Account
	Character	*Character
	Cipher     	security.Cipher
	Connection 	net.Conn
	Identity   	uint32
}

// Send an encrypted packet to the client. The encryption used is any cipher which 
// meets the Cipher interface. A copy buffer is created to encrypt the packet 
// without encrypting the original packet buffer.
func (c *Client) Send(packet interface{}) {
	
	// Create a copy of the buffer.
	writer := bytes.NewBuffer(nil)
	err := packets.Write(writer, packet)
	if err != nil { fmt.Println(err) }
	buffer := writer.Bytes()
	
	// Write the length of the packet to offset 0 (NetDragon byte ordering).
	binary.LittleEndian.PutUint16(buffer[0:2], uint16(len(buffer)))
	
	// Encrypt and send the buffer to the client.
	c.Cipher.Encrypt(buffer)
	c.Connection.Write(buffer)
}