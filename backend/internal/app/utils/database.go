package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./data/kasseapparat.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func ConnectToLocalDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func PurgeDatabase(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Product{}, &models.Purchase{}, &models.PurchaseItem{}, &models.User{}, models.List{}, models.ListEntry{})
	if err != nil {
		panic(err)
	}
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Product{}, &models.Purchase{}, &models.PurchaseItem{}, &models.User{}, models.List{}, models.ListEntry{})
	if err != nil {
		panic(err)
	}
}

func SeedDatabase(db *gorm.DB) {
	_ = gofakeit.Seed(0)

	db.Create(&models.User{Username: "admin", Email: "admin@example.com", Password: "admin", PasswordChangeRequired: false, Admin: true})
	db.Create(&models.User{Username: "demo",  Email: "demo@example.com", Password: "demo", PasswordChangeRequired: false, Admin: false})

	db.Create(&models.Product{Name: "🎟️ Regular", Price: 40, Pos: 1, ApiExport: true})
	reducedProduct := &models.Product{Name: "🎟️ Reduced", Price: 20, Pos: 2, ApiExport: true}
	db.Create(reducedProduct)
	freeProduct := &models.Product{Name: "🎟️ Free", Price: 0, Pos: 3, ApiExport: true}
	db.Create(freeProduct)
	prepaidProduct := &models.Product{Name: "🎟️ Prepaid", Price: 0, Pos: 4, WrapAfter: true, ApiExport: true}
	db.Create(prepaidProduct)
	db.Create(&models.Product{Name: "👕 T-Shirt Male S", Price: 20, Pos: 10})
	db.Create(&models.Product{Name: "👕 T-Shirt Male M", Price: 20, Pos: 11})
	db.Create(&models.Product{Name: "👕 T-Shirt Male L", Price: 20, Pos: 12})
	db.Create(&models.Product{Name: "👕 T-Shirt Male XL", Price: 20, Pos: 13, WrapAfter: true})
	db.Create(&models.Product{Name: "👕 T-Shirt Female S", Price: 20, Pos: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Female M", Price: 20, Pos: 21})
	db.Create(&models.Product{Name: "👕 T-Shirt Female L", Price: 20, Pos: 22})
	db.Create(&models.Product{Name: "👕 T-Shirt Female XL", Price: 20, Pos: 23, WrapAfter: true})
	db.Create(&models.Product{Name: "☕ Coffee Mug", Price: 1, Pos: 30})
	
	reducedDkevList := &models.List{Name: "Reduces Digitale Kultur", ProductID: reducedProduct.ID}
	db.Create(reducedDkevList)
	for i := 1; i < 5; i++ { 
		db.Create(&models.ListEntry{Name: gofakeit.Name(), ListID: reducedDkevList.ID, AdditionalGuests: 0})
	}

	reducedLdList := &models.List{Name: "Long Distance", ProductID: reducedProduct.ID}
	db.Create(reducedLdList)
	for i := 1; i < 15; i++ { 
		db.Create(&models.ListEntry{Name: gofakeit.Name(), ListID: reducedLdList.ID, AdditionalGuests: 0})
	}

	deineTicketsList := &models.List{Name: "Deine Tickets", TypeCode: true, ProductID: prepaidProduct.ID}
	db.Create(deineTicketsList)
	for i := 1; i < 20; i++ { 
		code :=gofakeit.Password(false, true, true, false, false, 9);
		db.Create(&models.ListEntry{Name: gofakeit.Name(), Code: &code, ListID: deineTicketsList.ID, AdditionalGuests: 0})
	}

	for i := 1; i < 8; i++ { 
		userGuestList := &models.List{Name: "Guestlist " + gofakeit.FirstName(), ProductID: freeProduct.ID}
		db.Create(userGuestList)

		for j := 0; j < gofakeit.Number(1, 10); j++ {

			db.Create(&models.ListEntry{Name: gofakeit.Name(), ListID: userGuestList.ID, AdditionalGuests: uint(gofakeit.Number(0, 2))})
		}
	}	
}
