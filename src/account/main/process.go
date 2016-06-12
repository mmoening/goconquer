package main

import (
	"account/handles"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"lib/structures"
	"lib/packets"
	"lib/security"
)

// OnConnect is called by the auth server to initialize the client structure upon
// connection. Client ciphers should be initialized here. If the server has been 
// upgraded to the new RC5 seed procedure from around patch 5180, then send the 
// MsgEncryptCode packet here after cipher initialization. 
func OnConnect(client *structures.Client) {
	client.Cipher = new(security.TQCipher)
	client.Cipher.Init()
}

// OnReceive is called by the auth server to handle packets from the player's
// game client. At this point, the server has assembled fragments and split the 
// packet buffer into multiple packets. This function accepts a decrypted packet 
// to be processed in the passed buffer.
func OnReceive(client *structures.Client, buffer *bytes.Buffer) {
	b := buffer.Bytes()
	length := binary.LittleEndian.Uint16(b[0:2])
	identity := binary.LittleEndian.Uint16(b[2:4])
	switch identity {
		
	// 1051: MsgAccount
	case packets.MSGACCOUNT:
		packet := new(packets.MsgAccount)
		err := packets.Read(buffer, packet)
		if err != nil { fmt.Println(err) } else { 
			handles.AuthenticateLogin(client, packet) 
		}
	default:
		fmt.Println("Missing packet handle:", identity, "length", length)
		fmt.Println(hex.Dump(b))
	}
}