package handles

import (
	"account/db"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"lib/packets"
	"lib/structures"
	"lib/security"
	"time"
	"strings"
)

// AuthenticateLogin checks the user's account and password combination after
// decrypting the password sent across the MsgAccount packet. The user's account 
// loaded from the flat-file database, then sent to the game server for granted
// access.
func AuthenticateLogin(client *structures.Client, p *packets.MsgAccount) {
	client.Account = new(structures.Account)
	if db.LoadAccount(client.Account, p.Account) {
	
		// Check if the user is banned.
		if client.Account.Status == structures.ACCTSTATUS_BANNED {
			response := packets.NewMsgConnectEx()
			response.Token = 12
			copy(response.Address[:], packets.MSGCONNECTEX_BANNED_ACCOUNT)
			client.Send(response)
		} else {
			// Create a new account for the client.
			client.Identity = client.Account.Identity
			cipher := security.RC5 { }
			cipher.Init()
			
			// Decrypt the password from the client and hash it. Decrypting the 
			// password here since the ciphertext is so weak. A plain SHA1 is already 
			// effective without a thin middle layer of ciphertext.
			hash := sha1.New()
			cipher.Decrypt(p.Password[:])
			hash.Write(p.Password[:])
			password := fmt.Sprintf("%x", hash.Sum(nil))
	
			// Verify that the password is correct.
			if strings.Compare(password, client.Account.Password) == 0 {
				gameserver, exists := db.Kernel.GameServers[p.Server]
				if exists && gameserver.Connection != nil { 
					
					// Send authentication details to the game server.
					transfer := structures.Transfer {}
					transfer.Account = *client.Account
					transfer.IPAddress = client.Connection.RemoteAddr().String()
					transfer.Requested = time.Now()
					encoder := gob.NewEncoder(gameserver.Connection)
					err := encoder.Encode(transfer)
					if err != nil {
						
						// Failed to send. Notify the player of server downtime.
						response := packets.NewMsgConnectEx()
						response.Token = 10
						copy(response.Address[:], packets.MSGCONNECTEX_SERVER_DOWN[:])
						client.Send(response)
						
					} else { // Correct response. 
						// Forward the client to the game server.
						response := packets.NewMsgConnectEx()
						response.Identity = client.Account.Identity
						response.Token = client.Account.Authority
						copy(response.Address[:], gameserver.Host)
						response.Port = gameserver.Port
						client.Send(response)
						
						// Success message to make Matt feel better about himself.
						fmt.Printf("Authenticated %s for %s\n", 
							client.Account.Username, gameserver.Name)
					}
				} else { // The server doesn't exist or is offline.
					response := packets.NewMsgConnectEx()
					response.Token = 10
					copy(response.Address[:], packets.MSGCONNECTEX_SERVER_DOWN)
					client.Send(response)
				}
			} else { // Invalid username or password.
				response := packets.NewMsgConnectEx()
				response.Token = 1
				copy(response.Address[:], packets.MSGCONNECTEX_INVALID_ACCOUNT)
				client.Send(response)
			}
		}
	} else { // Invalid username or password.
		response := packets.NewMsgConnectEx()
		response.Token = 1
		copy(response.Address[:], packets.MSGCONNECTEX_INVALID_ACCOUNT)
		client.Send(response)
	}
}