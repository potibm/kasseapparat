package sqlite

import (
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestProductGroup(
	t *testing.T,
	repo *Repository,
	sectionID int,
	name string,
	pos int,
) models.ProductGroup {
	t.Helper()

	group, err := repo.CreateProductGroup(models.ProductGroup{SectionID: sectionID, Name: name, Pos: pos})
	require.NoError(t, err)
	require.NotZero(t, group.ID)

	return group
}

func TestProductGroups_CreateAndGetByID(t *testing.T) {
	repo := setupSectionRepository(t)

	shirts := createTestSection(t, repo, "T-Shirts", 1)
	created := createTestProductGroup(t, repo, shirts.ID, "Male", 10)

	loaded, err := repo.GetProductGroupByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Male", loaded.Name)
	assert.Equal(t, shirts.ID, loaded.SectionID)
	assert.Equal(t, 10, loaded.Pos)
}

func TestProductGroups_GetByID_NotFound(t *testing.T) {
	repo := setupSectionRepository(t)

	_, err := repo.GetProductGroupByID(4711)

	assert.ErrorIs(t, err, ErrProductGroupNotFound)
}

func TestProductGroups_ListSortedAndFiltered(t *testing.T) {
	repo := setupSectionRepository(t)

	shirts := createTestSection(t, repo, "T-Shirts", 1)
	merch := createTestSection(t, repo, "Merch", 2)

	female := createTestProductGroup(t, repo, shirts.ID, "Female", 20)
	male := createTestProductGroup(t, repo, shirts.ID, "Male", 10)
	hoodies := createTestProductGroup(t, repo, merch.ID, "Hoodie", 5)

	all, err := repo.GetProductGroups(10, 0, "pos", "ASC", ProductFilters{})
	require.NoError(t, err)
	require.Len(t, all, 3)
	assert.Equal(t, []int{hoodies.ID, male.ID, female.ID}, []int{all[0].ID, all[1].ID, all[2].ID})

	filters := ProductFilters{SectionIDs: []int{shirts.ID}}

	inShirts, err := repo.GetProductGroups(10, 0, "name", "DESC", filters)
	require.NoError(t, err)
	require.Len(t, inShirts, 2)
	assert.Equal(t, "Male", inShirts[0].Name)

	total, err := repo.GetTotalProductGroups(filters)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)

	byID, err := repo.GetProductGroups(10, 0, "sectionId", "bogus", ProductFilters{IDs: []int{hoodies.ID}})
	require.NoError(t, err)
	require.Len(t, byID, 1)
	assert.Equal(t, "Hoodie", byID[0].Name)
}

func TestProductGroups_ListRejectsUnknownSortField(t *testing.T) {
	repo := setupSectionRepository(t)

	_, err := repo.GetProductGroups(10, 0, "drop table", "ASC", ProductFilters{})

	assert.ErrorIs(t, err, ErrInvalidSortField)
}

func TestProductGroups_UpdateMovesTheProductsWithTheGroup(t *testing.T) {
	repo := setupSectionRepository(t)

	shirts := createTestSection(t, repo, "T-Shirts", 1)
	merch := createTestSection(t, repo, "Merch", 2)
	group := createTestProductGroup(t, repo, shirts.ID, "Male", 10)

	grouped, err := repo.CreateProduct(models.Product{
		Name: "Male S", SectionID: shirts.ID, ProductGroupID: &group.ID, Variant: "S",
	})
	require.NoError(t, err)

	standalone, err := repo.CreateProduct(models.Product{Name: "Sticker", SectionID: shirts.ID})
	require.NoError(t, err)

	updated, err := repo.UpdateProductGroupByID(group.ID, models.ProductGroup{
		SectionID: merch.ID, Name: "Men", Pos: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, "Men", updated.Name)
	assert.Equal(t, 3, updated.Pos)
	assert.Equal(t, merch.ID, updated.SectionID)

	reloaded, err := repo.GetProductByID(grouped.ID)
	require.NoError(t, err)
	assert.Equal(t, merch.ID, reloaded.SectionID)

	untouched, err := repo.GetProductByID(standalone.ID)
	require.NoError(t, err)
	assert.Equal(t, shirts.ID, untouched.SectionID)
}

func TestProductGroups_UpdateNotFound(t *testing.T) {
	repo := setupSectionRepository(t)

	_, err := repo.UpdateProductGroupByID(4711, models.ProductGroup{Name: "Ghost"})

	assert.ErrorIs(t, err, ErrProductGroupNotFound)
}

func TestProductGroups_Delete(t *testing.T) {
	repo := setupSectionRepository(t)

	shirts := createTestSection(t, repo, "T-Shirts", 1)
	empty := createTestProductGroup(t, repo, shirts.ID, "Empty", 1)
	full := createTestProductGroup(t, repo, shirts.ID, "Male", 2)

	_, err := repo.CreateProduct(models.Product{Name: "Male S", SectionID: shirts.ID, ProductGroupID: &full.ID})
	require.NoError(t, err)

	require.NoError(t, repo.DeleteProductGroup(empty))

	_, err = repo.GetProductGroupByID(empty.ID)
	require.ErrorIs(t, err, ErrProductGroupNotFound)

	err = repo.DeleteProductGroup(full)
	require.ErrorIs(t, err, ErrProductGroupNotEmpty)

	_, err = repo.GetProductGroupByID(full.ID)
	assert.NoError(t, err, "a group with products must survive")
}

func TestSections_DeleteRefusesSectionWithProductGroups(t *testing.T) {
	repo := setupSectionRepository(t)

	shirts := createTestSection(t, repo, "T-Shirts", 1)
	createTestProductGroup(t, repo, shirts.ID, "Male", 10)

	err := repo.DeleteSection(shirts)

	assert.ErrorIs(t, err, ErrSectionNotEmpty)
}

func TestProducts_UpdateStoresGroupAndVariant(t *testing.T) {
	repo := setupSectionRepository(t)

	shirts := createTestSection(t, repo, "T-Shirts", 1)
	group := createTestProductGroup(t, repo, shirts.ID, "Male", 10)

	product, err := repo.CreateProduct(models.Product{Name: "Male S", SectionID: shirts.ID})
	require.NoError(t, err)

	product.ProductGroupID = &group.ID
	product.Variant = "S"

	updated, err := repo.UpdateProductByID(product.ID, product)
	require.NoError(t, err)
	require.NotNil(t, updated.ProductGroupID)
	assert.Equal(t, group.ID, *updated.ProductGroupID)
	assert.Equal(t, "S", updated.Variant)

	updated.ProductGroupID = nil

	ungrouped, err := repo.UpdateProductByID(product.ID, *updated)
	require.NoError(t, err)
	assert.Nil(t, ungrouped.ProductGroupID)
}
