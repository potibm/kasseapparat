package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
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

func SeedDatabase(db *gorm.DB, includeTestData bool) {
	const (
		DefaultGuestlistCount            = 38
		DefaultPurchaseCount             = 10
		MaxNotPresentEntriesPerGuestlist = 10
		MaxPresentEntriesPerGuestlist    = 2
	)

	_ = gofakeit.Seed(0)

	vat0 := decimal.NewFromInt(0)
	vat7 := decimal.NewFromInt(7)
	vat19 := decimal.NewFromInt(19)

	price0 := decimal.NewFromInt(0)
	price40GrossAt7 := decimal.NewFromFloat(37.38)
	price20GrossAt7 := decimal.NewFromFloat(18.69)
	price1GrossAt19 := decimal.NewFromFloat(0.84)
	price20GrossAt19 := decimal.NewFromFloat(16.81)

	// USERS
	adminUser := &models.User{Username: "admin", Email: "admin@example.com", Password: "admin", Admin: true}
	db.Create(adminUser)

	demoUser := &models.User{Username: "demo", Email: "demo@example.com", Password: "demo", Admin: false}
	db.Create(demoUser)

	// PRODUCTS
	products := []models.Product{}
	regularProduct := models.Product{Name: "🎟️ Regular", NetPrice: price40GrossAt7, VATRate: vat7, Pos: 1, ApiExport: true}
	db.Create(&regularProduct)
	products = append(products, regularProduct)

	reducedProduct := models.Product{Name: "🎟️ Reduced", NetPrice: price20GrossAt7, VATRate: vat7, Pos: 2, ApiExport: true}
	db.Create(&reducedProduct)
	products = append(products, reducedProduct)

	freeProduct := models.Product{Name: "🎟️ Free", NetPrice: price0, VATRate: vat0, Pos: 3, ApiExport: true}
	db.Create(&freeProduct)
	products = append(products, freeProduct)

	prepaidProduct := models.Product{Name: "🎟️ Prepaid", NetPrice: price0, VATRate: vat0, Pos: 4, WrapAfter: true, ApiExport: true}
	db.Create(&prepaidProduct)
	products = append(products, prepaidProduct)

	products = append(products, models.Product{Name: "👕 Male S", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Male M", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Male L", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Male XL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Male XXL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Male 4XL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Femal S", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Femal M", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Femal L", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Femal XL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "👕 Femal XXL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	products = append(products, models.Product{Name: "☕ Coffee Mug", NetPrice: price1GrossAt19, VATRate: vat19, Pos: 30})

	for i := range products {
		if products[i].ID == 0 {
			db.Create(&products[i])
		}
	}

	// GUESTLISTS
	reducedDkevGuestlist := &models.Guestlist{Name: "Reduces Digitale Kultur", ProductID: reducedProduct.ID}
	db.Create(reducedDkevGuestlist)

	reducedLdGuestlist := &models.Guestlist{Name: "Long Distance", ProductID: reducedProduct.ID}
	db.Create(reducedLdGuestlist)

	deineTicketsGuestlist := &models.Guestlist{Name: "Deine Tickets", TypeCode: true, ProductID: prepaidProduct.ID}
	db.Create(deineTicketsGuestlist)

	// GUESTS
	if includeTestData {
		for i := 1; i < 5; i++ {
			db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: reducedDkevGuestlist.ID, AdditionalGuests: 0})
		}

		for i := 1; i < 15; i++ {
			db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: reducedLdGuestlist.ID, AdditionalGuests: 0})
		}

		for i := 1; i < 20; i++ {
			code := gofakeit.Password(false, true, true, false, false, 9)
			db.Create(&models.Guest{Name: gofakeit.Name(), Code: &code, GuestlistID: deineTicketsGuestlist.ID, AdditionalGuests: 0})
		}
	}

	if includeTestData {
		for i := 1; i < DefaultGuestlistCount; i++ {
			userGuestlist := &models.Guestlist{Name: "Guestlist " + gofakeit.FirstName(), ProductID: freeProduct.ID}
			db.Create(userGuestlist)

			for range gofakeit.Number(1, MaxNotPresentEntriesPerGuestlist) {
				db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: userGuestlist.ID, AdditionalGuests: uint(gofakeit.Number(0, 2))})
			}

			for range gofakeit.Number(1, MaxPresentEntriesPerGuestlist) {
				arrivedAt := gofakeit.Date()
				db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: userGuestlist.ID, AdditionalGuests: uint(gofakeit.Number(0, 2)), AttendedGuests: 1, ArrivedAt: &arrivedAt})
			}
		}
	}

	// PURCHASES
	if includeTestData {
		_ = db.Transaction(func(tx *gorm.DB) error {
			// create some purchases
			for i := 1; i < DefaultPurchaseCount; i++ {
				purchase := models.Purchase{
					PaymentMethod:   gofakeit.RandomString([]string{"CASH", "CC"}),
					TotalGrossPrice: decimal.NewFromInt(0),
					TotalNetPrice:   decimal.NewFromInt(0),
				}

				// select productsCount-times a random product and a quantity
				for j := 0; j < gofakeit.Number(1, 5); j++ {
					product := products[gofakeit.Number(0, len(products)-1)]

					quantity := gofakeit.Number(1, 3)
					purchaseItem := models.PurchaseItem{
						ProductID: product.ID,
						Quantity:  quantity,
						NetPrice:  product.NetPrice,
						VATRate:   product.VATRate,
					}

					purchase.TotalGrossPrice = purchase.TotalGrossPrice.Add(purchaseItem.TotalGrossPrice(2))
					purchase.TotalNetPrice = purchase.TotalNetPrice.Add(purchaseItem.TotalNetPrice(2))

					purchase.PurchaseItems = append(purchase.PurchaseItems, purchaseItem)
				}

				purchase.CreatedBy = demoUser

				db.Create(&purchase)
			}

			return nil
		})
	}
}
