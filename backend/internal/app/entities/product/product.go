package product

type Product struct {
	ID                        uint
	Username                  string
	Email                     string
	Password                  string
	Admin                     bool
	ChangePasswordToken       *string
	ChangePasswordTokenExpiry *int64
}
