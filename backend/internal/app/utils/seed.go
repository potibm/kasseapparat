package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type DatabaseSeed struct {
	products              []models.Product
	db                    *gorm.DB
	regularProduct        *models.Product
	reducedProduct        *models.Product
	freeProduct           *models.Product
	prepaidProduct        *models.Product
	reducedDkevGuestlist  *models.Guestlist
	reducedLdGuestlist    *models.Guestlist
	deineTicketsGuestlist *models.Guestlist
	demoUser              *models.User
	adminUser             *models.User
}

func NewDatabaseSeed(db *gorm.DB) *DatabaseSeed {
	return &DatabaseSeed{
		db: db,
	}
}

func (ds *DatabaseSeed) Seed(includeTestData bool) {
	const (
		DefaultGuestlistCount            = 38
		DefaultPurchaseCount             = 10
		MaxNotPresentEntriesPerGuestlist = 10
		MaxPresentEntriesPerGuestlist    = 2
	)

	_ = gofakeit.Seed(1)

	ds.seedUsers()
	ds.seedProducts()
	ds.seedGuestlists()

	if includeTestData {
		ds.seedGuests()
		ds.seedUserGuests(DefaultGuestlistCount, MaxNotPresentEntriesPerGuestlist, MaxPresentEntriesPerGuestlist)
		ds.seedPurchases(DefaultPurchaseCount)
	}
}

func (ds *DatabaseSeed) seedUsers() {
	ds.adminUser = &models.User{Username: "admin", Email: "admin@example.com", Password: "admin", Admin: true}

	ds.demoUser = &models.User{Username: "demo", Email: "demo@example.com", Password: "demo", Admin: false}

	ds.db.Create(ds.adminUser)
	ds.db.Create(ds.demoUser)
}

func (ds *DatabaseSeed) seedProducts() {
	vat0 := decimal.NewFromInt(0)
	vat7 := decimal.NewFromInt(7)
	vat19 := decimal.NewFromInt(19)

	price0 := decimal.NewFromInt(0)
	price40GrossAt7 := decimal.NewFromFloat(37.38)
	price20GrossAt7 := decimal.NewFromFloat(18.69)
	price1GrossAt19 := decimal.NewFromFloat(0.84)
	price20GrossAt19 := decimal.NewFromFloat(16.81)

	ds.regularProduct = &models.Product{Name: "🎟️ Regular", NetPrice: price40GrossAt7, VATRate: vat7, Pos: 1, ApiExport: true}
	ds.db.Create(ds.regularProduct)

	ds.reducedProduct = &models.Product{Name: "🎟️ Reduced", NetPrice: price20GrossAt7, VATRate: vat7, Pos: 2, ApiExport: true}
	ds.db.Create(ds.reducedProduct)

	ds.freeProduct = &models.Product{Name: "🎟️ Free", NetPrice: price0, VATRate: vat0, Pos: 3, ApiExport: true}
	ds.db.Create(ds.freeProduct)

	ds.prepaidProduct = &models.Product{Name: "🎟️ Prepaid", NetPrice: price0, VATRate: vat0, Pos: 4, WrapAfter: true, ApiExport: true}
	ds.db.Create(ds.prepaidProduct)

	ds.products = append(ds.products, *ds.prepaidProduct)

	ds.products = append(ds.products, models.Product{Name: "👕 Male S", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Male M", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Male L", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Male XL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Male XXL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Male 4XL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Female S", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Female M", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Female L", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Female XL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "👕 Female XXL", NetPrice: price20GrossAt19, VATRate: vat19, Pos: 10, TotalStock: gofakeit.IntRange(5, 30)})
	ds.products = append(ds.products, models.Product{Name: "☕ Coffee Mug", NetPrice: price1GrossAt19, VATRate: vat19, Pos: 30})

	for i := range ds.products {
		if ds.products[i].ID == 0 {
			ds.db.Create(&ds.products[i])
		}
	}
}

func (ds *DatabaseSeed) seedGuestlists() {
	ds.reducedDkevGuestlist = &models.Guestlist{Name: "Reduces Digitale Kultur", ProductID: ds.reducedProduct.ID}
	ds.db.Create(ds.reducedDkevGuestlist)

	ds.reducedLdGuestlist = &models.Guestlist{Name: "Long Distance", ProductID: ds.reducedProduct.ID}
	ds.db.Create(ds.reducedLdGuestlist)

	ds.deineTicketsGuestlist = &models.Guestlist{Name: "Deine Tickets", TypeCode: true, ProductID: ds.prepaidProduct.ID}
	ds.db.Create(ds.deineTicketsGuestlist)
}

func (ds *DatabaseSeed) seedGuests() {
	for i := 1; i < 5; i++ {
		ds.db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: ds.reducedDkevGuestlist.ID, AdditionalGuests: 0})
	}

	for i := 1; i < 15; i++ {
		ds.db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: ds.reducedLdGuestlist.ID, AdditionalGuests: 0})
	}

	for i := 1; i < 20; i++ {
		code := gofakeit.Password(false, true, true, false, false, 9)
		ds.db.Create(&models.Guest{Name: gofakeit.Name(), Code: &code, GuestlistID: ds.deineTicketsGuestlist.ID, AdditionalGuests: 0})
	}
}

func (ds *DatabaseSeed) seedUserGuests(guestlistCount, maxNotPresentEntries, maxPresentEntries int) {
	if guestlistCount <= 0 {
		return
	}

	for i := 1; i < guestlistCount; i++ {
		userGuestlist := &models.Guestlist{Name: "Guestlist " + gofakeit.FirstName(), ProductID: ds.freeProduct.ID}
		ds.db.Create(userGuestlist)

		for range gofakeit.Number(1, maxNotPresentEntries) {
			ds.db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: userGuestlist.ID, AdditionalGuests: uint(gofakeit.Number(0, 2))})
		}

		for range gofakeit.Number(1, maxPresentEntries) {
			arrivedAt := gofakeit.Date()
			ds.db.Create(&models.Guest{Name: gofakeit.Name(), GuestlistID: userGuestlist.ID, AdditionalGuests: uint(gofakeit.Number(0, 2)), AttendedGuests: 1, ArrivedAt: &arrivedAt})
		}
	}
}

func (ds *DatabaseSeed) seedPurchases(purchaseCount int) {
	if purchaseCount <= 0 {
		return
	}

	_ = ds.db.Transaction(func(tx *gorm.DB) error {
		for i := 1; i < purchaseCount; i++ {
			purchase := models.Purchase{
				PaymentMethod:   gofakeit.RandomString([]string{"CASH", "CC"}),
				TotalGrossPrice: decimal.NewFromInt(0),
				TotalNetPrice:   decimal.NewFromInt(0),
			}

			for j := 0; j < gofakeit.Number(1, 5); j++ {
				product := ds.products[gofakeit.Number(0, len(ds.products)-1)]

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

			purchase.CreatedByID = &ds.demoUser.ID

			ds.db.Create(&purchase)
		}

		return nil
	})
}
