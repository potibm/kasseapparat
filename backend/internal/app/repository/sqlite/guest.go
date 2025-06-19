package sqlite

import (
	"errors"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/gorm"
)

var (
	ErrGuestNotFound  = errors.New("guest not found")
	ErrGuestsNotFound = errors.New("guests not found")
)

type GuestFilters struct {
	Query       string
	GuestlistID int
	ListGroupId int
	Present     bool
	NotPresent  bool
	IDs         []int
}

var guestSortFieldMappings = map[string]string{
	"id":             "Guests.ID",
	"name":           "Guests.Name",
	"guestlist.name": "Guestlist.Name",
	"arrivedAt":      "Guests.arrived_at",
}

func (filters GuestFilters) AddWhere(query *gorm.DB) *gorm.DB {
	if len(filters.IDs) > 0 {
		query = query.Where("Guests.ID IN ?", filters.IDs)
	}

	if filters.Query != "" {
		query = query.Where("Guests.Name LIKE ? OR Guests.Code LIKE ?", "%"+filters.Query+"%", filters.Query+"%")
	}

	if filters.GuestlistID != 0 {
		query = query.Where("Guests.guestlist_id = ?", filters.GuestlistID)
	}

	if filters.Present {
		query = query.Where("Guests.attended_guests > 0")
	}

	if filters.NotPresent {
		query = query.Where("Guests.attended_guests = 0")
	}

	return query
}

func (repo *Repository) GetGuests(limit int, offset int, sort string, order string, filters GuestFilters) ([]models.Guest, error) {
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	sort, err := getGuestsValidSortFieldName(sort)
	if err != nil {
		return nil, err
	}

	var guests []models.Guest

	query := repo.db.Joins("Guestlist").Order(sort + " " + order + ", Guests.ID ASC").Limit(limit).Offset(offset)
	query = filters.AddWhere(query)

	if err := query.Find(&guests).Error; err != nil {
		return nil, ErrGuestsNotFound
	}

	return guests, nil
}

func getGuestsValidSortFieldName(input string) (string, error) {
	if field, exists := guestSortFieldMappings[input]; exists {
		return field, nil
	}

	return "", errors.New("invalid sort field name")
}

func (repo *Repository) GetTotalGuests(filters *GuestFilters) (int64, error) {
	var totalRows int64

	query := repo.db.Model(&models.Guest{}).Joins("Guestlist")
	if filters != nil {
		query = filters.AddWhere(query)
	}

	query.Count(&totalRows)

	return totalRows, nil
}

func (repo *Repository) GetUnattendedGuestsByProductID(productId int, q string) (models.GuestSummarySlice, error) {
	var guests models.GuestSummarySlice

	query := repo.db.Model(&models.Guest{}).
		Select("Guests.id, Guests.name, Guests.code, Guestlists.name AS list_name, Guests.additional_guests, Guests.arrival_note").
		Joins("JOIN guestlists ON Guests.guestlist_id = Guestlists.id").
		Joins("JOIN products ON Guestlists.product_id = Products.id").
		Where("Products.id = ? AND Guests.attended_guests = ?", productId, 0).
		Order("guests.name ASC")
	if q != "" {
		query = query.Where("Guests.name LIKE ? OR Guests.code = ?", "%"+q+"%", q)
	}

	if err := query.Scan(&guests).Error; err != nil {
		return nil, ErrGuestsNotFound
	}

	return guests, nil
}

func (repo *Repository) GetGuestByID(id int) (*models.Guest, error) {
	var guest models.Guest
	if err := repo.db.First(&guest, id).Error; err != nil {
		return nil, ErrGuestNotFound
	}

	return &guest, nil
}

func (repo *Repository) GetGuestByCode(code string) (*models.Guest, error) {
	var guest models.Guest
	if err := repo.db.Where("code = ?", code).First(&guest).Error; err != nil {
		return nil, ErrGuestNotFound
	}

	return &guest, nil
}

func (repo *Repository) GetFullGuestByID(id int) (*models.Guest, error) {
	var guest models.Guest
	if err := repo.db.Preload("Guestlist").Preload("Guestlist.Product").First(&guest, id).Error; err != nil {
		return nil, ErrGuestNotFound
	}

	return &guest, nil
}

func (repo *Repository) UpdateGuestByID(id int, updatedGuest models.Guest) (*models.Guest, error) {
	var guest models.Guest
	if err := repo.db.First(&guest, id).Error; err != nil {
		return nil, ErrGuestNotFound
	}

	updatedGuest.ID = guest.ID

	if err := repo.db.Save(&updatedGuest).Error; err != nil {
		return nil, errors.New("failed to update guest")
	}

	return &updatedGuest, nil
}

func (repo *Repository) CreateGuest(guest models.Guest) (models.Guest, error) {
	result := repo.db.Create(&guest)

	return guest, result.Error
}

func (repo *Repository) DeleteGuest(guest models.Guest, deletedBy models.User) {
	repo.db.Model(&models.Guest{}).Where(whereIDEquals, guest.ID).Update("DeletedByID", deletedBy.ID)

	repo.db.Delete(&guest)
}

func (repo *Repository) RollbackVisitedGuestsByPurchaseID(purchaseId uuid.UUID) error {
	repo.db.Model(&models.Guest{}).
		Where("purchase_id = ?", purchaseId.String()).
		Updates(map[string]interface{}{"purchase_id": nil, "attended_guests": 0, "arrived_at": nil})

	return nil
}
