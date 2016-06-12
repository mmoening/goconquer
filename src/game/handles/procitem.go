package handles

import (
	"fmt"
	"encoding/hex"
	"lib/packets"
	"lib/structures"
)

// ProcItem processes a request for item action processing or ping processing from
// a player. All inventory actions, weapon management, warehouse actions, etc.
// including a random ping-back request are handled here.
func ProcItem(c *structures.Client, p *packets.MsgItem, b []byte) {
	switch (p.Action) {
	
	default:
		fmt.Println("Missing packet handle:", p.Identifier, "length", p.Length)
		fmt.Println(hex.Dump(b))
		c.Send(p)
	}
}