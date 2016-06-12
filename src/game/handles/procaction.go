package handles

import (
	"encoding/hex"
	"fmt"
	"lib/structures"
	"lib/packets"
)

// ProcessMsgAction uses a subtype called Action which determines which routine 
// will be called to process the packet. These subtypes are defined with the 
// packet structure. The buffer is passed into the function for default hex dump.
func ProcAction(c *structures.Client, p *packets.MsgAction, b []byte) {
	switch (p.Action) {
	
	case packets.ACTION_SETLOCATION: 	SetLocation(c, p)
	
	default:
		fmt.Println("Missing packet handle:", p.Identifier, "length", p.Length)
		fmt.Println(hex.Dump(b))
		c.Send(p)
	}
}

// SetLocation is called after character initialization on login to initialize
// the location of the character. It isn't called after the login sequence.
func SetLocation(c *structures.Client, p *packets.MsgAction) {
	p.X = c.Character.X
	p.Y = c.Character.Y
	p.Data = c.Character.Map
	c.Send(p)
}