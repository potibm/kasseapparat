package sqlite

import (
	"context"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"gorm.io/gorm"
)

const whereIDEquals = "id = ?"

type Repository struct {
	db            *gorm.DB
	decimalPlaces int32
}

type RepositoryInterface interface {
	GetDB() *gorm.DB
	WithTransaction(ctx context.Context, fn func(repo RepositoryInterface) error) error
	GetGuests(limit int, offset int, sort string, order string, filters GuestFilters) ([]models.Guest, error)
	GetTotalGuests(filters *GuestFilters) (int64, error)
	GetUnattendedGuestsByProductID(productId int, q string) (models.GuestSummarySlice, error)
	GetGuestByID(id int) (*models.Guest, error)
	GetGuestByCode(code string) (*models.Guest, error)
	GetFullGuestByID(id int) (*models.Guest, error)
	UpdateGuestByID(id int, updatedGuest models.Guest) (*models.Guest, error)
	CreateGuest(guest models.Guest) (models.Guest, error)
	DeleteGuest(guest models.Guest, deletedBy models.User)
	RollbackVisitedGuestsByPurchaseID(purchaseId uuid.UUID) error
	GetGuestlists(limit int, offset int, sort string, order string, filters GuestlistFilters) ([]models.Guestlist, error)
	GetTotalGuestlists() (int64, error)
	GetGuestlistByID(id int) (*models.Guestlist, error)
	GetGuestlistWithTypeCode() (*models.Guestlist, error)
	UpdateGuestlistByID(id int, updatedGuestlist models.Guestlist) (*models.Guestlist, error)
	CreateGuestlist(guestlist models.Guestlist) (models.Guestlist, error)
	DeleteGuestlist(guestlist models.Guestlist, deletedBy models.User)
	GetProductInterests(limit int, offset int, ids []int) ([]models.ProductInterest, error)
	GetTotalProductInterests() (int64, error)
	GetProductInterestByID(id int) (*models.ProductInterest, error)
	DeleteProductInterest(productInterest models.ProductInterest, deletedBy models.User)
	CreateProductInterest(productInterest models.ProductInterest, createdBy models.User) (models.ProductInterest, error)
	GetProductInterestCountByProductID(productID uint) (int, error)
	GetProductStats() ([]response.ProductStats, error)
	GetProducts(limit int, offset int, sort string, order string, ids []int) ([]models.Product, error)
	GetTotalProducts() (int64, error)
	GetProductByID(id int) (*models.Product, error)
	UpdateProductByID(id int, updatedProduct models.Product) (*models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	DeleteProduct(product models.Product, deletedBy models.User)
	GetAttendedGuestSumByProductID(productID uint) (int, error)
	StorePurchases(purchase models.Purchase) (models.Purchase, error)
	DeletePurchaseByID(id uuid.UUID, deletedBy models.User)
	GetPurchaseByID(id uuid.UUID) (*models.Purchase, error)
	UpdatePurchaseStatusByID(id uuid.UUID, status models.PurchaseStatus) (*models.Purchase, error)
	UpdatePurchaseSumupTransactionIDByID(id uuid.UUID, sumupTransactionID uuid.UUID) (*models.Purchase, error)
	UpdatePurchaseSumupClientTransactionIDByID(id uuid.UUID, sumupClientTransactionID uuid.UUID) (*models.Purchase, error)
	GetPurchases(limit int, offset int, sort string, order string, filters PurchaseFilters) ([]models.Purchase, error)
	GetFilteredPurchases(filters PurchaseFilters) ([]models.PurchaseItem, error)
	GetTotalPurchases(filters PurchaseFilters) (int64, error)
	GetPurchaseStats() ([]ProductPurchaseStats, error)
	GetPurchasedQuantitiesByProductID(productID uint) (int, error)
}

var _ RepositoryInterface = (*Repository)(nil)

func NewRepository(db *gorm.DB, decimalPlaces int32) *Repository {
	return &Repository{db: db, decimalPlaces: decimalPlaces}
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *Repository) WithTransaction(ctx context.Context, fn func(repo RepositoryInterface) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := r.cloneWithDB(tx)
		txRepo.db.Debug()

		return fn(txRepo)
	})
}

func (r *Repository) cloneWithDB(tx *gorm.DB) *Repository {
	return &Repository{db: tx, decimalPlaces: r.decimalPlaces}
}
