package sumup

import "time"

type Reader struct {
	ID               string
	Name             string
	Status           string
	DeviceIdentifier string
	DeviceModel      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
