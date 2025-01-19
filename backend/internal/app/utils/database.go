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
	err := db.Migrator().DropTable(&models.Product{}, &models.Purchase{}, &models.PurchaseItem{}, &models.User{}, models.Guestlist{}, models.Guest{}, models.ProductInterest{})
	if err != nil {
		panic(err)
	}
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Product{}, &models.Purchase{}, &models.PurchaseItem{}, &models.User{}, models.Guestlist{}, models.Guest{}, models.ProductInterest{})
	if err != nil {
		panic(err)
	}
}

func SeedDatabase(db *gorm.DB) {

	const (
		DefaultGuestlistCount            = 38
		MaxNotPresentEntriesPerGuestlist = 10
		MaxPresentEntriesPerGuestlist    = 2
	)

	_ = gofakeit.Seed(0)

	db.Create(&models.User{Username: "admin", Email: "admin@example.com", Password: "admin", Admin: true})
	db.Create(&models.User{Username: "demo", Email: "demo@example.com", Password: "demo", Admin: false})

	db.Create(&models.Product{Name: "🎟️ Regular", Price: 40, Pos: 1, ApiExport: true})
	reducedProduct := &models.Product{Name: "🎟️ Reduced", Price: 20, Pos: 2, ApiExport: true}
	db.Create(reducedProduct)
	freeProduct := &models.Product{Name: "🎟️ Free", Price: 0, Pos: 3, ApiExport: true}
	db.Create(freeProduct)
	prepaidProduct := &models.Product{Name: "🎟️ Prepaid", Price: 0, Pos: 4, WrapAfter: true, ApiExport: true}
	db.Create(prepaidProduct)
	db.Create(&models.Product{Name: "👕 Male S", Price: 20, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Male M", Price: 20, Pos: 11, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Male L", Price: 20, Pos: 12, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Male XL", Price: 20, Pos: 13, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Male XXL", Price: 20, Pos: 15, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Male XXXL", Price: 20, Pos: 16, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Male 4XL", Price: 20, Pos: 17, WrapAfter: true, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Female S", Price: 20, Pos: 20, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Female M", Price: 20, Pos: 21, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Female L", Price: 20, Pos: 22, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Female XL", Price: 20, Pos: 23, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "👕 Female XXL", Price: 20, Pos: 24, WrapAfter: true, TotalStock: gofakeit.IntRange(5, 30)})
	db.Create(&models.Product{Name: "☕ Coffee Mug", Price: 1, Pos: 30})

	reducedDkevGuestlist := &models.Guestlist{Name: "Reduces Digitale Kultur", ProductID: reducedProduct.ID}
	db.Create(reducedDkevGuestlist)
	for i := 1; i < 5; i++ {
		db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: reducedDkevGuestlist.ID, AdditionalGuests: 0})
	}

	reducedLdGuestlist := &models.Guestlist{Name: "Long Distance", ProductID: reducedProduct.ID}
	db.Create(reducedLdGuestlist)
	for i := 1; i < 15; i++ {
		db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: reducedLdGuestlist.ID, AdditionalGuests: 0})
	}

	deineTicketsGuestlist := &models.Guestlist{Name: "Deine Tickets", TypeCode: true, ProductID: prepaidProduct.ID}
	db.Create(deineTicketsGuestlist)
	for i := 1; i < 20; i++ {
		code := gofakeit.Password(false, true, true, false, false, 9)
		db.Create(&models.Guest{Name: gofakeit.Name(), Code: &code, GuestlistID: deineTicketsGuestlist.ID, AdditionalGuests: 0})
	}

	for i := 1; i < DefaultGuestlistCount; i++ {
		userGuestlist := &models.Guestlist{Name: "Guestlist " + gofakeit.FirstName(), ProductID: freeProduct.ID}
		db.Create(userGuestlist)

		for j := 0; j < gofakeit.Number(1, MaxNotPresentEntriesPerGuestlist); j++ {

			db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: userGuestlist.ID, AdditionalGuests: uint(gofakeit.Number(0, 2))})
		}
		for j := 0; j < gofakeit.Number(1, MaxPresentEntriesPerGuestlist); j++ {

			arrivedAt := gofakeit.Date()
			db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: userGuestlist.ID, AdditionalGuests: uint(gofakeit.Number(0, 2)), AttendedGuests: 1, ArrivedAt: &arrivedAt})
		}

	}
}
