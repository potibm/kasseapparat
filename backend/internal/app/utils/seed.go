package utils

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/kasseapparat/internal/app/models"
	gormaudit "github.com/potibm/kasseapparat/internal/app/store/gorm"
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
}

func NewDatabaseSeed(db *gorm.DB) *DatabaseSeed {
	return &DatabaseSeed{
		db: db,
	}
}

func (ds *DatabaseSeed) Seed(includeTestData bool) error {
	const (
		DefaultGuestlistCount            = 38
		DefaultPurchaseCount             = 30
		MaxNotPresentEntriesPerGuestlist = 10
		MaxPresentEntriesPerGuestlist    = 2
	)

	_ = gofakeit.Seed(1)

	// Each step depends on the IDs the previous one populated, so a failure has to
	// stop the seed rather than leave half-written data behind: a product whose
	// Create failed has ID 0, and every guestlist seeded after it would point at
	// ProductID 0.
	if err := ds.seedProducts(); err != nil {
		return err
	}

	if err := ds.seedGuestlists(); err != nil {
		return err
	}

	if !includeTestData {
		return nil
	}

	if err := ds.seedGuests(); err != nil {
		return err
	}

	if err := ds.seedUserGuests(
		DefaultGuestlistCount,
		MaxNotPresentEntriesPerGuestlist,
		MaxPresentEntriesPerGuestlist,
	); err != nil {
		return err
	}

	return ds.seedPurchases(DefaultPurchaseCount)
}

func (ds *DatabaseSeed) seedProducts() error {
	vat0 := decimal.NewFromInt(0)
	vat7 := decimal.NewFromInt(7)
	vat19 := decimal.NewFromInt(19)

	price0 := decimal.NewFromInt(0)
	price40GrossAt7 := decimal.NewFromFloat(37.38)
	price20GrossAt7 := decimal.NewFromFloat(18.69)
	price1GrossAt19 := decimal.NewFromFloat(0.84)
	price20GrossAt19 := decimal.NewFromFloat(16.81)

	ds.regularProduct = &models.Product{
		Name:      "🎟️ Regular",
		NetPrice:  price40GrossAt7,
		VATRate:   vat7,
		Pos:       1,
		APIExport: true,
	}
	if err := ds.db.Create(ds.regularProduct).Error; err != nil {
		return fmt.Errorf("failed to create product %q: %w", ds.regularProduct.Name, err)
	}

	ds.reducedProduct = &models.Product{
		Name:      "🎟️ Reduced",
		NetPrice:  price20GrossAt7,
		VATRate:   vat7,
		Pos:       2,
		APIExport: true,
	}
	if err := ds.db.Create(ds.reducedProduct).Error; err != nil {
		return fmt.Errorf("failed to create product %q: %w", ds.reducedProduct.Name, err)
	}

	ds.freeProduct = &models.Product{Name: "🎟️ Free", NetPrice: price0, VATRate: vat0, Pos: 3, APIExport: true}
	if err := ds.db.Create(ds.freeProduct).Error; err != nil {
		return fmt.Errorf("failed to create product %q: %w", ds.freeProduct.Name, err)
	}

	ds.prepaidProduct = &models.Product{
		Name:      "🎟️ Prepaid",
		NetPrice:  price0,
		VATRate:   vat0,
		Pos:       4,
		WrapAfter: true,
		APIExport: true,
	}
	if err := ds.db.Create(ds.prepaidProduct).Error; err != nil {
		return fmt.Errorf("failed to create product %q: %w", ds.prepaidProduct.Name, err)
	}

	ds.products = append(ds.products, *ds.prepaidProduct)

	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Male S",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Male M",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Male L",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Male XL",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Male XXL",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Male 4XL",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Female S",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Female M",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Female L",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Female XL",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{
			Name:       "👕 Female XXL",
			NetPrice:   price20GrossAt19,
			VATRate:    vat19,
			Pos:        10,
			TotalStock: gofakeit.IntRange(5, 30),
		},
	)
	ds.products = append(
		ds.products,
		models.Product{Name: "☕ Coffee Mug", NetPrice: price1GrossAt19, VATRate: vat19, Pos: 30},
	)

	for i := range ds.products {
		if ds.products[i].ID == 0 {
			if err := ds.db.Create(&ds.products[i]).Error; err != nil {
				return fmt.Errorf("failed to create product %q: %w", ds.products[i].Name, err)
			}
		}
	}

	return nil
}

func (ds *DatabaseSeed) seedGuestlists() error {
	ds.reducedDkevGuestlist = &models.Guestlist{Name: "Reduces Digitale Kultur", ProductID: ds.reducedProduct.ID}
	if err := ds.db.Create(ds.reducedDkevGuestlist).Error; err != nil {
		return fmt.Errorf("failed to create guestlist %q: %w", ds.reducedDkevGuestlist.Name, err)
	}

	ds.reducedLdGuestlist = &models.Guestlist{Name: "Long Distance", ProductID: ds.reducedProduct.ID}
	if err := ds.db.Create(ds.reducedLdGuestlist).Error; err != nil {
		return fmt.Errorf("failed to create guestlist %q: %w", ds.reducedLdGuestlist.Name, err)
	}

	ds.deineTicketsGuestlist = &models.Guestlist{Name: "Deine Tickets", TypeCode: true, ProductID: ds.prepaidProduct.ID}
	if err := ds.db.Create(ds.deineTicketsGuestlist).Error; err != nil {
		return fmt.Errorf("failed to create guestlist %q: %w", ds.deineTicketsGuestlist.Name, err)
	}

	return nil
}

func (ds *DatabaseSeed) seedGuests() error {
	for i := 1; i < 5; i++ {
		guest := &models.Guest{Name: gofakeit.Name(), GuestlistID: ds.reducedDkevGuestlist.ID, AdditionalGuests: 0}
		if err := ds.db.Create(guest).Error; err != nil {
			return fmt.Errorf("failed to create guest %q: %w", guest.Name, err)
		}
	}

	for i := 1; i < 15; i++ {
		guest := &models.Guest{Name: gofakeit.Name(), GuestlistID: ds.reducedLdGuestlist.ID, AdditionalGuests: 0}
		if err := ds.db.Create(guest).Error; err != nil {
			return fmt.Errorf("failed to create guest %q: %w", guest.Name, err)
		}
	}

	for i := 1; i < 20; i++ {
		code := gofakeit.Password(false, true, true, false, false, 9)
		guest := &models.Guest{
			Name:             gofakeit.Name(),
			Code:             &code,
			GuestlistID:      ds.deineTicketsGuestlist.ID,
			AdditionalGuests: 0,
		}

		if err := ds.db.Create(guest).Error; err != nil {
			return fmt.Errorf("failed to create guest %q: %w", guest.Name, err)
		}
	}

	// for e2e test: create a guest with a known code in the deineTicketsGuestlist
	code := "ABCDEFGHI"
	knownCodeGuest := &models.Guest{
		Name:             "Jan Jansen",
		Code:             &code,
		GuestlistID:      ds.deineTicketsGuestlist.ID,
		AdditionalGuests: 0,
	}

	if err := ds.db.Create(knownCodeGuest).Error; err != nil {
		return fmt.Errorf("failed to create guest %q: %w", knownCodeGuest.Name, err)
	}

	return nil
}

func getOptionalArrivalNote() *string {
	if gofakeit.Number(1, 100) <= 20 {
		words := gofakeit.Sentence(5)

		return &words
	}

	return nil
}

func (ds *DatabaseSeed) seedUserGuests(guestlistCount, maxNotPresentEntries, maxPresentEntries int) error {
	if guestlistCount <= 0 {
		return nil
	}

	for i := 1; i < guestlistCount; i++ {
		if err := ds.seedUserGuestlist(maxNotPresentEntries, maxPresentEntries); err != nil {
			return err
		}
	}

	return ds.seedKnownE2EGuests()
}

func (ds *DatabaseSeed) seedUserGuestlist(maxNotPresentEntries, maxPresentEntries int) error {
	userGuestlist := &models.Guestlist{Name: "Guestlist " + gofakeit.FirstName(), ProductID: ds.freeProduct.ID}
	if err := ds.db.Create(userGuestlist).Error; err != nil {
		return fmt.Errorf("failed to create guestlist %q: %w", userGuestlist.Name, err)
	}

	for range gofakeit.Number(1, maxNotPresentEntries) {
		guest := &models.Guest{
			Name:             gofakeit.Name(),
			GuestlistID:      userGuestlist.ID,
			AdditionalGuests: gofakeit.UintRange(0, 2),
			ArrivalNote:      getOptionalArrivalNote(),
		}

		if err := ds.db.Create(guest).Error; err != nil {
			return fmt.Errorf("failed to create guest %q: %w", guest.Name, err)
		}
	}

	for range gofakeit.Number(1, maxPresentEntries) {
		arrivedAt := gofakeit.Date()
		guest := &models.Guest{
			Name:             gofakeit.Name(),
			GuestlistID:      userGuestlist.ID,
			AdditionalGuests: gofakeit.UintRange(0, 2),
			AttendedGuests:   1,
			ArrivedAt:        &arrivedAt,
		}

		if err := ds.db.Create(guest).Error; err != nil {
			return fmt.Errorf("failed to create guest %q: %w", guest.Name, err)
		}
	}

	return nil
}

// seedKnownE2EGuests creates the two guests the e2e suite searches for by name.
func (ds *DatabaseSeed) seedKnownE2EGuests() error {
	// for e2e test: create two guests with known names in a special guestlist
	userGuestlist := &models.Guestlist{Name: "E2E Guestlist " + gofakeit.FirstName(), ProductID: ds.freeProduct.ID}
	if err := ds.db.Create(userGuestlist).Error; err != nil {
		return fmt.Errorf("failed to create guestlist %q: %w", userGuestlist.Name, err)
	}

	jean := &models.Guest{
		Name:             "Jean Dupont",
		GuestlistID:      userGuestlist.ID,
		AdditionalGuests: uint(gofakeit.UintRange(0, 2)),
	}
	if err := ds.db.Create(jean).Error; err != nil {
		return fmt.Errorf("failed to create guest %q: %w", jean.Name, err)
	}

	note := "Ciao Mario!"

	mario := &models.Guest{
		Name:             "Mario Rossi",
		GuestlistID:      userGuestlist.ID,
		AdditionalGuests: uint(gofakeit.UintRange(0, 2)),
		ArrivalNote:      &note,
	}
	if err := ds.db.Create(mario).Error; err != nil {
		return fmt.Errorf("failed to create guest %q: %w", mario.Name, err)
	}

	return nil
}

func (ds *DatabaseSeed) seedPurchases(purchaseCount int) error {
	if purchaseCount <= 0 {
		return nil
	}

	ctx := gormaudit.WithUserID(context.Background(), "seed")

	if err := ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 1; i < purchaseCount; i++ {
			purchase := models.Purchase{
				// generate a random PaymentMethod from models.PaymentMethodCash and models.PaymentMethodCC
				PaymentMethod: models.PaymentMethod(
					gofakeit.RandomString([]string{string(models.PaymentMethodCash), string(models.PaymentMethodCC)}),
				),
				TotalGrossPrice: decimal.NewFromInt(0),
				TotalNetPrice:   decimal.NewFromInt(0),
			}

			for j := 0; j < gofakeit.Number(1, 5); j++ {
				product := ds.products[gofakeit.Number(0, len(ds.products)-1)]

				quantity := gofakeit.UintRange(1, 3)
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

			// Write through tx, never the outer ds.db: this transaction already
			// holds the only pool connection, so an outer-handle write would wait
			// for a connection that is never released. See ADR 003.
			if err := tx.WithContext(ctx).Create(&purchase).Error; err != nil {
				return fmt.Errorf("failed to seed purchase %d: %w", i, err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to seed purchases: %w", err)
	}

	return nil
}
