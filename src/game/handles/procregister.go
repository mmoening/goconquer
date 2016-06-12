package handles

import (
	"fmt"
	"game/db"
	"lib/packets"
	"lib/structures"
	"math/rand"
	"os"
	"strings"
	"sync"
)

// Global server variables for character creation.
var mutex sync.Mutex
var hairstyles []uint16 = []uint16{
	10, 11, 13, 14, 15, 24, 30, 35,
	37, 38, 39, 40, 43, 50, 72, 74}

// ProcRegister is sent by the game client to request character creation. The 
// character name, body, class ,etc. should be verified before saving the 
// character to the file system. This patch disconnects after creating the 
// character. More recent patches around 5300+ resend the MsgConnectEx packet.
func ProcRegister(client *structures.Client, p *packets.MsgRegister) {

	// Is this a flood attack?
	if client.Character != nil || 
		!db.Kernel.CharacterCreationPool.Contains(client.Identity) {
		client.Connection.Close()
		return	
	}
	
	// Validate input from the player.
	if (p.Model != 2001 && p.Model != 2002 && p.Model != 1003 && p.Model != 1004) ||
	    (p.Class != 10 && p.Class != 20 && p.Class != 40 && p.Class != 100) ||
		client.Identity != p.Identity {	
		client.Connection.Close()
		return
	}

	// Check the name of the character.
	if strings.Contains(p.Name, "GM") || strings.Contains(p.Name, "PM") {
		client.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
			"ALLUSERS", "Name is invalid.", packets.MSGTALK_ENTRANCE))
		return
	}

	// Perform file operations, lock to prevent race conditions.
	mutex.Lock()
	if _, err := os.Stat("./characters/" + p.Name + ".json"); err == nil {
		client.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
			"ALLUSERS", "Name is taken.", packets.MSGTALK_ENTRANCE))
		return
	}

	// Create the character file and unlock.
	file, err := os.Create("./characters/" + p.Name + ".json")
	file.Close()
	mutex.Unlock()
	if err != nil {
		fmt.Println("error: failed to create character file")
		client.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
			"ALLUSERS", "Database error.", packets.MSGTALK_ENTRANCE))
		return
	}

	// Initialize character.
	character := new(structures.Character)
	character.Model = p.Model
	character.Class = byte(p.Class)
	character.Identity = client.Identity
	character.Level = 1
	character.Map = 1010
	character.Name = p.Name
	character.Silver = 10000
	character.Spouse = "None"
	character.X = 61
	character.Y = 109

	// Obtain attributes from the database for that class.
	attrib := db.Attributes.Get(character.Class, 1)
	character.Agility = attrib[db.AGILITY]
	character.Spirit = attrib[db.SPIRIT]
	character.Strength = attrib[db.STRENGTH]
	character.Vitality = attrib[db.VITALITY]

	// Generate random characteristics for the character.
	character.Avatar = uint16(rand.Intn(50))
	character.Hairstyle = uint16((rand.Intn(7)) + 3) * 100 +
		hairstyles[rand.Intn(len(hairstyles))]
	character.Mana = character.Spirit * 5
	character.Health = character.Strength * 3 + character.Agility * 3 +
		character.Spirit * 3 + character.Vitality * 24
	if p.Model == 2001 || p.Model == 2002 {
		character.Avatar += 200
	}

	// Save the character to file.
	if !db.Characters.Save(character) ||
		!db.Characters.AppendIndex(character.Identity, character.Name) {
		client.Send(packets.NewMsgTalk(p.Identity, "SYSTEM",
			"ALLUSERS", "Database error.", packets.MSGTALK_ENTRANCE))
		fmt.Println("error: failed to save character file")
	}

	// Respond to the client.
	db.Kernel.CharacterCreationPool.Remove(client.Identity)
	client.Send(packets.NewMsgTalk(p.Identity, "SYSTEM", "ALLUSERS",
		"ANSWER_OK", packets.MSGTALK_ENTRANCE))
}
