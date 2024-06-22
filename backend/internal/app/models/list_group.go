package models

// List represents a guestlist
type ListGroup struct {
	GormOwnedModel
	ListID 	uint `json:"listId"`
	List List `json:"list"`
	Name      string  `json:"name"`
}
