package structures

// Character is saved to the flat-file database for persistent character data. 
// Temporary character data should not be stored here, but instead should be stored
// in other structures, linked to from the Client structure. 
type Character struct {
	Identity uint32
	Name, Spouse string
	Model, Avatar, Hairstyle uint16
	Silver, CPs uint32
	Level, Rebirths byte
	Experience uint64
	Class, PreviousClass byte
	Map uint32
	X, Y uint16
	Health, Mana uint16
	Attributes, Strength, Agility, Vitality, Spirit, PkPoints uint16
}
