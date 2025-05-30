package sumup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sumup/sumup-go/transactions"
)

func TestParseHrefToListTransactionsParams(t *testing.T) {
	href := "limit=1&oldest_ref=0dd170c7-d82a-4fec-b2c0-e6de01c631a8&order=ascending&skip_tx_result=true&changes_since=2024-01-01T12%3A00%3A00Z&users=test1%40example.com&users=test2%40example.com"

	params, err := parseHrefToListTransactionsParams(href)
	assert.NoError(t, err)
	assert.NotNil(t, params)

	assert.NotNil(t, params.Limit)
	assert.Equal(t, 1, *params.Limit)

	assert.NotNil(t, params.OldestRef)
	assert.Equal(t, "0dd170c7-d82a-4fec-b2c0-e6de01c631a8", *params.OldestRef)

	assert.NotNil(t, params.Order)
	assert.Equal(t, "ascending", *params.Order)

	assert.NotNil(t, params.ChangesSince)

	expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	assert.True(t, params.ChangesSince.Equal(expectedTime))

	assert.NotNil(t, params.Users)
	assert.Equal(t, []string{"test1@example.com", "test2@example.com"}, *params.Users)

	// real world example
	href = "limit=1&oldest_ref=0dd170c7-d82a-4fec-b2c0-e6de01c631a8&order=ascending&skip_tx_result=true"
	params, err = parseHrefToListTransactionsParams(href)
	assert.NoError(t, err)

	assert.NotNil(t, params.Limit)
	assert.Equal(t, 1, *params.Limit)
	assert.NotNil(t, params.OldestRef)
	assert.Equal(t, "0dd170c7-d82a-4fec-b2c0-e6de01c631a8", *params.OldestRef)
	assert.NotNil(t, params.Order)
	assert.Equal(t, "ascending", *params.Order)
	assert.Nil(t, params.ChangesSince)
	assert.Nil(t, params.Users)
	assert.Nil(t, params.NewestRef)
	assert.Nil(t, params.TransactionCode)
	assert.Nil(t, params.OldestTime)
	assert.Nil(t, params.NewestTime)
	assert.Nil(t, params.Statuses)
	assert.Nil(t, params.Types)
}

func TestFindNextHref(t *testing.T) {
	links := []transactions.Link{
		{Rel: nil, Href: nil},
		{Rel: nil, Href: nil},
		{Rel: nil, Href: nil},
	}

	nextHref := findNextHref(&links)
	assert.Equal(t, "", nextHref)

	links = []transactions.Link{
		{Rel: nil, Href: nil},
		{Rel: &[]string{"next"}[0], Href: &[]string{"https://example.com/next"}[0]},
	}

	nextHref = findNextHref(&links)
	assert.Equal(t, "https://example.com/next", nextHref)
}
