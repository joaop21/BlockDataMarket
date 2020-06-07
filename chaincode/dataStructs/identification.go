package dataStructs

// Identification object that represents an identification in the World State
type Identification struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

// Constructor for Identification
func NewIdentification(id string, name string, publicKey string) *Identification {

	return &Identification{
		Type:            "Identification",
		Id:				 id,
		Name:            name,
		PublicKey:		 publicKey,
	}

}