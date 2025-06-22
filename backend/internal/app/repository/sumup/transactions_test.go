package sumup

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/sumup/sumup-go/shared"
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

func TestFromSDKTransaction_WithUUIDId(t *testing.T) {
	id := "2b5cd782-0733-4fb2-bf22-5a12345bd94f"
	tid := shared.TransactionId(id)
	tc := "TAAAABCP2SA"
	amount := 40.0
	currency := shared.Currency("EUR")
	cardType := transactions.TransactionHistoryCardType("MASTERCARD")
	status := transactions.TransactionHistoryStatus("SUCCESSFUL")

	timestamp := parseTime(t, "2025-06-15T20:45:27.588Z")

	sdk := &transactions.TransactionHistory{
		Id:              &id,
		TransactionCode: &tc,
		TransactionId:   &tid,
		Amount:          &amount,
		Currency:        &currency,
		CardType:        &cardType,
		Timestamp:       &timestamp,
		Status:          &status,
	}

	tx := fromSDKTransaction(sdk)

	assert.Equal(t, id, tx.ID)
	assert.Equal(t, tc, tx.TransactionCode)
	assert.Equal(t, uuid.MustParse(id), tx.TransactionID)
	assert.Equal(t, "EUR", tx.Currency)
	assert.Equal(t, "MASTERCARD", tx.CardType)
	assert.Equal(t, "SUCCESSFUL", tx.Status)
	assert.WithinDuration(t, timestamp, tx.CreatedAt, time.Second)
}

func TestFromSDKTransaction_WithNonUUIDId(t *testing.T) {
	id := "8119994131" // not a UUID
	tc := "TAAAABCP2SA"
	tid := shared.TransactionId("2b5cd782-0733-4fb2-bf22-5a12345bd94f")
	amount := 40.0
	currency := shared.Currency("EUR")
	cardType := transactions.TransactionHistoryCardType("MASTERCARD")
	status := transactions.TransactionHistoryStatus("REFUNDED")
	timestamp := parseTime(t, "2025-06-15T21:00:13.536Z")

	sdk := &transactions.TransactionHistory{
		Id:              &id,
		TransactionCode: &tc,
		TransactionId:   &tid,
		Amount:          &amount,
		Currency:        &currency,
		CardType:        &cardType,
		Timestamp:       &timestamp,
		Status:          &status,
	}

	tx := fromSDKTransaction(sdk)

	assert.Equal(t, id, tx.ID)
	assert.Equal(t, tc, tx.TransactionCode)
	assert.Equal(t, uuid.MustParse(string(tid)), tx.TransactionID)
	assert.Equal(t, "EUR", tx.Currency)
	assert.Equal(t, "MASTERCARD", tx.CardType)
	assert.Equal(t, "REFUNDED", tx.Status)
	assert.WithinDuration(t, timestamp, tx.CreatedAt, time.Second)
}

func parseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("invalid time: %v", err)
	}

	return ts
}
