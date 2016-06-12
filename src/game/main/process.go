package main

import (
	"game/db"
	"game/handles"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"lib/structures"
	"lib/packets"
	"lib/security"
)

// OnConnect is called by the game server to initialize the client structure upon
// connection. Client ciphers should be initialized here. If the server has been 
// upgraded to use the DH exchange (after 5017), then the exchange may be 
// initialized here as well. In addition, the OnExchange event will need to be 
// defined.
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
		
	/* 1001: MsgRegister */ 
	case packets.MSGREGISTER:
		packet := new(packets.MsgRegister)
		err := packets.Read(buffer, packet)
		if err != nil { fmt.Println(err) } else { 
			handles.ProcRegister(client, packet) 
		}
	/* 1009: MsgItem */ 
	case packets.MSGITEM:
		packet := new(packets.MsgItem)
		err := packets.Read(buffer, packet)
		if err != nil { fmt.Println(err) } else { 
			handles.ProcItem(client, packet, b) 
		}
	/* 1010: MsgAction */ 
	case packets.MSGACTION:
		packet := new(packets.MsgAction)
		err := packets.Read(buffer, packet)
		if err != nil { fmt.Println(err) } else { 
			handles.ProcAction(client, packet, b) 
		}
	/* 1052: MsgConnect */ 
	case packets.MSGCONNECT:
		packet := new(packets.MsgConnect)
		err := packets.Read(buffer, packet)
		if err != nil { fmt.Println(err) } else { 
			handles.ProcConnect(client, packet) 
		}
	default:
		fmt.Println("Missing packet handle:", identity, "length", length)
		fmt.Println(hex.Dump(b))
	}
}

// OnDisconnect is called by the game server to dispose of client structures
// and stop in-progress actions from the client (such as trading or being a map 
// entity) after disconnect.
func OnDisconnect(client *structures.Client) {
	if client != nil && client.Character != nil {
		if db.Kernel.CharacterCreationPool.Contains(client.Identity) {
			db.Kernel.CharacterCreationPool.Remove(client.Identity)
		}
		fmt.Printf("%s disconnected.\n", client.Character.Name)
	}
}	
