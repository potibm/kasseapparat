package models

// List represents a guestlist
type ListEntry struct {
	GormOwnedModel
	ListID 	uint	`json:"listId"`
	List 	List	`json:"list"`
	Name      string  `json:"name" `
	Code 		string`json:"code"`
	ListGroupID	*uint 	`json:"listGroupId"`
	ListGroup	*ListGroup 	`json:"listGroup"`
	AdditionalGuests uint `json:"additionalGuests" gorm:"default:0"`
	AttendedGuests uint `json:"attendedGuests" gorm:"default:0"`
}
