package types

type Profile struct {
	ID   ProfileID `json:"id"`   // Id of the profile.
	Name string    `json:"name"` // Name of the profile.
	Icon string    `json:"icon"` // Path to the profile's icon.
}
