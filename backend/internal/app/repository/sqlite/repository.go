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

type TransactionalRepository interface {
	GetDB() *gorm.DB
	WithTransaction(ctx context.Context, fn func(repo RepositoryInterface) error) error
}

type GuestRepository interface {
	GuestCRUDRepository
	GetGuestsByPurchaseID(purchaseID uuid.UUID) ([]models.Guest, error)
	GetUnattendedGuestsByProductID(productId int, q string) (models.GuestSummarySlice, error)
	GetGuestByCode(code string) (*models.Guest, error)
	GetFullGuestByID(id int) (*models.Guest, error)
	RollbackVisitedGuestsByPurchaseID(purchaseId uuid.UUID) error
}

type GuestCRUDRepository interface {
	GetGuests(limit int, offset int, sort string, order string, filters GuestFilters) ([]models.Guest, error)
	GetGuestByID(id int) (*models.Guest, error)
	UpdateGuestByID(id int, updatedGuest models.Guest) (*models.Guest, error)
	CreateGuest(guest models.Guest) (models.Guest, error)
	DeleteGuest(guest models.Guest, deletedBy models.User)
	GetTotalGuests(filters *GuestFilters) (int64, error)
}

type GuestlistRepository interface {
	GetGuestlists(limit int, offset int, sort string, order string, filters GuestlistFilters) ([]models.Guestlist, error)
	GetTotalGuestlists() (int64, error)
	GetGuestlistByID(id int) (*models.Guestlist, error)
	GetGuestlistWithTypeCode() (*models.Guestlist, error)
	UpdateGuestlistByID(id int, updatedGuestlist models.Guestlist) (*models.Guestlist, error)
	CreateGuestlist(guestlist models.Guestlist) (models.Guestlist, error)
	DeleteGuestlist(guestlist models.Guestlist, deletedBy models.User)
}

type ProductInterestRepository interface {
	GetProductInterests(limit int, offset int, ids []int) ([]models.ProductInterest, error)
	GetTotalProductInterests() (int64, error)
	GetProductInterestByID(id int) (*models.ProductInterest, error)
	DeleteProductInterest(productInterest models.ProductInterest, deletedBy models.User)
	CreateProductInterest(productInterest models.ProductInterest, createdBy models.User) (models.ProductInterest, error)
	GetProductInterestCountByProductID(productID uint) (int, error)
}

type ProductRepository interface {
	GetProductStats() ([]response.ProductStats, error)
	GetProducts(limit int, offset int, sort string, order string, ids []int) ([]models.Product, error)
	GetTotalProducts() (int64, error)
	GetProductByID(id int) (*models.Product, error)
	UpdateProductByID(id int, updatedProduct models.Product) (*models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	DeleteProduct(product models.Product, deletedBy models.User)
	GetAttendedGuestSumByProductID(productID uint) (int, error)
}

type PurchaseRepository interface {
	PurchaseCRUDRepository
	
	GetPurchaseBySumupClientTransactionID(sumupTransactionID uuid.UUID) (*models.Purchase, error)
	UpdatePurchaseStatusByID(id uuid.UUID, status models.PurchaseStatus) (*models.Purchase, error)
	UpdatePurchaseSumupTransactionIDByID(id uuid.UUID, sumupTransactionID uuid.UUID) (*models.Purchase, error)
	UpdatePurchaseSumupClientTransactionIDByID(id uuid.UUID, sumupClientTransactionID uuid.UUID) (*models.Purchase, error)
	GetFilteredPurchases(filters PurchaseFilters) ([]models.PurchaseItem, error)
	GetPurchaseStats() ([]ProductPurchaseStats, error)
	GetPurchasedQuantitiesByProductID(productID uint) (int, error)
}

type PurchaseCRUDRepository interface {
	StorePurchases(purchase models.Purchase) (models.Purchase, error)
	DeletePurchaseByID(id uuid.UUID, deletedBy models.User)
	GetPurchaseByID(id uuid.UUID) (*models.Purchase, error)
	GetTotalPurchases(filters PurchaseFilters) (int64, error)
	GetPurchases(limit int, offset int, sort string, order string, filters PurchaseFilters) ([]models.Purchase, error)
}


type UserRepository interface {
	GetUserByID(id int) (*models.User, error)
	GetUsers(limit int, offset int, sort string, order string, filters UserFilters) ([]models.User, error)
	GetTotalUsers(filters *UserFilters) (int64, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUserByID(id int, updatedUser models.User) (*models.User, error)
	DeleteUser(user models.User)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error)
}

type RepositoryInterface interface {
	TransactionalRepository
	GuestRepository
	GuestlistRepository
	ProductInterestRepository
	ProductRepository
	PurchaseRepository
	UserRepository
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

		return fn(txRepo)
	})
}

func (r *Repository) cloneWithDB(tx *gorm.DB) *Repository {
	return &Repository{db: tx, decimalPlaces: r.decimalPlaces}
}
