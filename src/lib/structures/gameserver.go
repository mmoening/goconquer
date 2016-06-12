package structures

import (
	"fmt"
	"net"
	"time"
)

// GameServer is used during the account to game server transfer. The client 
// specifies which game server to connect to in the MsgAccount packet. Then, the 
// account server sends a MsgConnectEx to forward the client to the correct game 
// server.
type GameServer struct {
	Name       string
	Host       string
	Port       uint32
	Backend    string
	Connection net.Conn
}

// Connect establishes a connection from the account server to the specified game 
// server. The game server must white list the account server in order to establish
// the connection.
func (g *GameServer) Connect() {
	for { // Reattempt after failure.
		for { // While the connection fails, reattempt once a second.
			connection, err := net.Dial("tcp", fmt.Sprintf("%s", g.Backend))
			if err == nil { 
				g.Connection = connection
				break
			} else { time.Sleep(time.Second) }
		}
		fmt.Printf("Connection established with %s\n", g.Name)

		for {
			// Attempt to read from the connection. Once this fails, the connection
			// has been broken and will need to be re-established.
			buffer := make([]byte, 1024)
			length, err := g.Connection.Read(buffer)
			if length > 0 && err == nil {
				
				// handle packet
				
			} else { break }
		}
		fmt.Printf("Connection lost with %s\n", g.Name)
		g.Connection = nil
	}
}
