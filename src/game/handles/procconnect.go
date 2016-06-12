package handles

import (
	"encoding/gob"
	"fmt"
	"game/db"
	"lib/packets"
	"lib/structures"
	"net"
	"strings"
)

// ProcConnect initializes the game client after the session has been
// authenticated by the account server. The game server implements this function
// to check if the player's character exists. If the character exists, it will
// be sent to the client to continue the login sequence; else, the login
// sequence will be interrupted for character creation.
func ProcConnect(c *structures.Client, p *packets.MsgConnect) {

	// Does the client exist in the authentication pool?
	if db.Kernel.AuthenticatedClients.Contains(p.Identity) {
		
		// Pull the transfer structure from the authentication pool.
		t := db.Kernel.AuthenticatedClients.Get(p.Identity).(*structures.Transfer)
		transferip := t.IPAddress[0:strings.Index(t.IPAddress, ":")]
		clientip := c.Connection.RemoteAddr().String()[0:strings.Index(
			c.Connection.RemoteAddr().String(), ":")]
		
		// Verify that the origins of the requests are the same.
		if strings.Compare(transferip, clientip) == 0 {
			c.Account = &t.Account
		} else { c.Connection.Close(); return }
	} else { c.Connection.Close(); return }
	db.Kernel.AuthenticatedClients.Remove(p.Identity)

	// Does an observer with the same account already exist on the server?
	observer := db.Kernel.ConnectedClients.Remove(p.Identity)
	if observer != nil { observer.(*structures.Client).Connection.Close() }
	if db.Kernel.ConnectedClients.Add(p.Identity, c) {

		// Generate keys for the client.
		c.Identity = p.Identity
		c.Cipher.Generate(p.Token, p.Identity)

		// Does the player's character exist?
		exists, err := db.Characters.Load(c)
		if exists {
			c.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
				"ALLUSERS", "ANSWER_OK", packets.MSGTALK_REGISTRATION))

			// Send character info to the client.
			packet := packets.NewMsgUserInfo()
			packet.Identity = c.Identity
			packet.Mesh = uint32(c.Character.Model) +
				uint32(c.Character.Avatar) * 10000
			packet.Hairstyle = c.Character.Hairstyle
			packet.Silver = c.Character.Silver
			packet.Experience = c.Character.Experience
			packet.Strength = c.Character.Strength
			packet.Agility = c.Character.Agility
			packet.Vitality = c.Character.Vitality
			packet.Spirit = c.Character.Spirit
			packet.Attributes = c.Character.Attributes
			packet.Health = c.Character.Health
			packet.Mana = c.Character.Mana
			packet.PkPoints = c.Character.PkPoints
			packet.Level = c.Character.Level
			packet.Class = c.Character.Class
			packet.Autoallot = c.Character.Rebirths > 0
			packet.Rebirths = c.Character.Rebirths
			packet.Strings = make([]string, 2)
			packet.Strings[0] = c.Character.Name
			packet.Strings[1] = c.Character.Spouse
			c.Send(packet)

		} else if err == nil {
			db.Kernel.CharacterCreationPool.Add(c.Identity, nil)
			c.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
				"ALLUSERS", "NEW_ROLE", packets.MSGTALK_REGISTRATION))
		} else {
			fmt.Println(err)
			c.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
				"ALLUSERS", "Database error", packets.MSGTALK_REGISTRATION))
		}
	} else { c.Connection.Close() }
}

// OpenAuthenticationChannel accepts a connection from the whitelisted account
// server, then receives the account information for authenticated players as 
// they login. The authentication channel helps cover a security vulnerability 
// where a player can impersonate another player. If you consider removing this
// system, recall that encryption is not a form of authentication. 
func OpenAuthenticationChannel() {

	// Listen for a new connection from the account server.
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d",
		db.Configuration.AuthPort))
	if err != nil { fmt.Println(err.Error()); return }
	defer listener.Close()

	for { // Reconnect to the account server on failure.
		// Accept the new connection
		connection, err := listener.Accept()
		if err != nil { fmt.Println(err.Error())
		} else { // Check if the connection is whitelisted.

			whitelisted := db.Configuration.AuthHost
			remote := connection.RemoteAddr().String()[
				0:strings.Index(connection.RemoteAddr().String(), ":")]
			if strings.Compare(remote, whitelisted) == 0 {

				fmt.Println("Connection established with account server")
				for { // Receive transfers from the connection.
					decoder := gob.NewDecoder(connection)
					transfer := &structures.Transfer{}
					err := decoder.Decode(transfer)
					if err != nil {
						fmt.Println("Disconnected from account server!")
						break
					}

					// Add the transfer to the accepted connections pool.
					db.Kernel.AuthenticatedClients.Add(
						transfer.Account.Identity, transfer)
				}
			}
		}
	}
}
