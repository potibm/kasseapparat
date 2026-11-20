package sqlite

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupSectionRepository starts from empty tables: the in-memory database is shared by
// every test in this package.
func setupSectionRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := utils.ConnectToLocalDatabase()
	require.NoError(t, err)
	require.NoError(t, utils.PurgeDatabase(db))
	require.NoError(t, utils.MigrateDatabase(db))

	return NewRepository(db, 2)
}

func createTestSection(t *testing.T, repo *Repository, name string, pos int) models.Section {
	t.Helper()

	section, err := repo.CreateSection(models.Section{Name: name, Pos: pos})
	require.NoError(t, err)
	require.NotZero(t, section.ID)

	return section
}

func TestSections_CreateAndGetByID(t *testing.T) {
	repo := setupSectionRepository(t)

	created := createTestSection(t, repo, "Tickets", 1)

	loaded, err := repo.GetSectionByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Tickets", loaded.Name)
	assert.Equal(t, 1, loaded.Pos)
}

func TestSections_GetByID_NotFound(t *testing.T) {
	repo := setupSectionRepository(t)

	_, err := repo.GetSectionByID(4711)

	assert.ErrorIs(t, err, ErrSectionNotFound)
}

func TestSections_ListSortedByPosWithTotal(t *testing.T) {
	repo := setupSectionRepository(t)

	merch := createTestSection(t, repo, "Merch", 3)
	tickets := createTestSection(t, repo, "Tickets", 1)
	shirts := createTestSection(t, repo, "T-Shirts", 2)

	sections, err := repo.GetSections(10, 0, "pos", "ASC", nil)
	require.NoError(t, err)
	require.Len(t, sections, 3)
	assert.Equal(t, []int{tickets.ID, shirts.ID, merch.ID}, []int{sections[0].ID, sections[1].ID, sections[2].ID})

	total, err := repo.GetTotalSections()
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
}

func TestSections_ListFiltersByIDs(t *testing.T) {
	repo := setupSectionRepository(t)

	first := createTestSection(t, repo, "A", 1)
	createTestSection(t, repo, "B", 2)

	sections, err := repo.GetSections(10, 0, "pos", "ASC", []int{first.ID})
	require.NoError(t, err)
	require.Len(t, sections, 1)
	assert.Equal(t, first.ID, sections[0].ID)
}

func TestSections_ListRejectsUnknownSortField(t *testing.T) {
	repo := setupSectionRepository(t)

	_, err := repo.GetSections(10, 0, "drop table", "ASC", nil)

	assert.Error(t, err)
}

func TestSections_Update(t *testing.T) {
	repo := setupSectionRepository(t)

	section := createTestSection(t, repo, "Tickets", 1)

	updated, err := repo.UpdateSectionByID(section.ID, models.Section{Name: "", Pos: 5})
	require.NoError(t, err)
	assert.Empty(t, updated.Name, "an empty name is allowed: the section is a plain row break")
	assert.Equal(t, 5, updated.Pos)

	_, err = repo.UpdateSectionByID(4711, models.Section{Name: "x"})
	assert.ErrorIs(t, err, ErrSectionNotFound)
}

func TestSections_DeleteEmptySection(t *testing.T) {
	repo := setupSectionRepository(t)

	section := createTestSection(t, repo, "Empty", 1)

	require.NoError(t, repo.DeleteSection(section))

	_, err := repo.GetSectionByID(section.ID)
	assert.ErrorIs(t, err, ErrSectionNotFound)
}

func TestSections_DeleteRefusesSectionWithProducts(t *testing.T) {
	repo := setupSectionRepository(t)

	section := createTestSection(t, repo, "Tickets", 1)
	_, err := repo.CreateProduct(models.Product{
		Name:      "Regular",
		NetPrice:  decimal.NewFromInt(10),
		SectionID: section.ID,
	})
	require.NoError(t, err)

	err = repo.DeleteSection(section)

	require.ErrorIs(t, err, ErrSectionNotEmpty)

	_, err = repo.GetSectionByID(section.ID)
	assert.NoError(t, err, "the section must still exist")
}

func TestSections_GetOrCreateDefaultSection_CreatesOneWhenNoneExists(t *testing.T) {
	repo := setupSectionRepository(t)

	section, err := repo.GetOrCreateDefaultSection()
	require.NoError(t, err)
	assert.NotZero(t, section.ID)
	assert.Empty(t, section.Name)

	again, err := repo.GetOrCreateDefaultSection()
	require.NoError(t, err)
	assert.Equal(t, section.ID, again.ID, "the default section is created only once")
}

func TestSections_GetOrCreateDefaultSection_ReturnsTheFirstByPos(t *testing.T) {
	repo := setupSectionRepository(t)

	createTestSection(t, repo, "Merch", 2)
	tickets := createTestSection(t, repo, "Tickets", 1)

	section, err := repo.GetOrCreateDefaultSection()
	require.NoError(t, err)
	assert.Equal(t, tickets.ID, section.ID)

	total, err := repo.GetTotalSections()
	require.NoError(t, err)
	assert.Equal(t, int64(2), total, "no extra section is created")
}

func TestProducts_UpdateKeepsTheSection(t *testing.T) {
	repo := setupSectionRepository(t)

	tickets := createTestSection(t, repo, "Tickets", 1)
	merch := createTestSection(t, repo, "Merch", 2)

	product, err := repo.CreateProduct(models.Product{Name: "Mug", SectionID: tickets.ID})
	require.NoError(t, err)

	product.SectionID = merch.ID
	updated, err := repo.UpdateProductByID(product.ID, product)
	require.NoError(t, err)
	assert.Equal(t, merch.ID, updated.SectionID)
}

func TestProducts_FilterBySection(t *testing.T) {
	repo := setupSectionRepository(t)

	tickets := createTestSection(t, repo, "Tickets", 1)
	merch := createTestSection(t, repo, "Merch", 2)

	regular, err := repo.CreateProduct(models.Product{Name: "Regular", Pos: 1, SectionID: tickets.ID})
	require.NoError(t, err)

	mug, err := repo.CreateProduct(models.Product{Name: "Mug", Pos: 2, SectionID: merch.ID})
	require.NoError(t, err)

	filters := ProductFilters{SectionIDs: []int{merch.ID}}

	products, err := repo.GetProducts(10, 0, "pos", "ASC", filters)
	require.NoError(t, err)
	require.Len(t, products, 1)
	assert.Equal(t, mug.ID, products[0].ID)

	total, err := repo.GetTotalProducts(filters)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)

	all, err := repo.GetTotalProducts(ProductFilters{})
	require.NoError(t, err)
	assert.Equal(t, int64(2), all)

	byID, err := repo.GetProducts(10, 0, "pos", "ASC", ProductFilters{IDs: []int{regular.ID}})
	require.NoError(t, err)
	require.Len(t, byID, 1)
	assert.Equal(t, "Regular", byID[0].Name)
}

func TestProducts_SortBySection(t *testing.T) {
	repo := setupSectionRepository(t)

	tickets := createTestSection(t, repo, "Tickets", 1)
	merch := createTestSection(t, repo, "Merch", 2)

	_, err := repo.CreateProduct(models.Product{Name: "Mug", Pos: 1, SectionID: merch.ID})
	require.NoError(t, err)

	_, err = repo.CreateProduct(models.Product{Name: "Regular", Pos: 2, SectionID: tickets.ID})
	require.NoError(t, err)

	products, err := repo.GetProducts(10, 0, "sectionId", "ASC", ProductFilters{})
	require.NoError(t, err)
	require.Len(t, products, 2)
	assert.Equal(t, "Regular", products[0].Name)
}
