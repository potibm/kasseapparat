package sumup

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/sumup/sumup-go/shared"
	"github.com/sumup/sumup-go/transactions"
)

const (
	TransactionPageSize = 100
	TransactionMaxPages = 5
)

func (r *Repository) GetTransactions(oldestFrom *time.Time) ([]Transaction, error) {
	ctx := context.Background()

	sdkItems, err := r.fetchPagedTransactions(ctx, oldestFrom, TransactionMaxPages, TransactionPageSize)
	if err != nil {
		return nil, err
	}

	result := make([]Transaction, 0, len(sdkItems))
	for _, sdkTx := range sdkItems {
		result = append(result, *fromSDKTransaction(sdkTx))
	}

	return result, nil
}

func (r *Repository) fetchPagedTransactions(ctx context.Context, oldestFrom *time.Time, maxPages, pageSize int) ([]*transactions.TransactionHistory, error) {
	var allItems []*transactions.TransactionHistory

	pageCount := 0

	params := transactions.ListParams{
		Limit: &pageSize,
	}
	if oldestFrom != nil {
		params.OldestTime = oldestFrom
	}

	for {
		if pageCount >= maxPages {
			return allItems, nil
		}

		pageCount++

		resp, err := r.service.Client.Transactions.List(ctx, r.service.MerchantCode, params)
		if err != nil {
			return nil, err
		}

		if resp.Items == nil || len(resp.Items) == 0 {
			return allItems, nil
		}

		allItems = append(allItems, ptrSliceToSlice(&resp.Items)...)
		sortTransactionsByCreatedAt(allItems)

		nextHref := findNextHref(&resp.Links)
		if nextHref == "" {
			return allItems, nil
		}

		nextParams, err := parseHrefToListTransactionsParams(nextHref)
		if err != nil {
			log.Printf("Error parsing next page link: %v", err)

			return allItems, nil
		}

		params = *nextParams
	}
}

func sortTransactionsByCreatedAt(transactions []*transactions.TransactionHistory) {
	if transactions == nil {
		return
	}

	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].Timestamp == nil || transactions[j].Timestamp == nil {
			return false
		}

		return transactions[i].Timestamp.After(*transactions[j].Timestamp)
	})
}

func ptrSliceToSlice(ptrSlice *[]transactions.TransactionHistory) []*transactions.TransactionHistory {
	if ptrSlice == nil {
		return nil
	}

	out := make([]*transactions.TransactionHistory, len(*ptrSlice))
	for i := range *ptrSlice {
		out[i] = &(*ptrSlice)[i]
	}

	return out
}

// findNextHref extracts the Href of the link whose Rel is "next".
// If links is nil or no such link is present, it returns an empty string.
func findNextHref(links *[]transactions.Link) string {
	if links == nil {
		return ""
	}

	for _, link := range *links {
		if link.Rel != nil && *link.Rel == "next" && link.Href != nil {
			return *link.Href
		}
	}

	return ""
}

// parseHrefToListTransactionsParams parses a URL query string into a transactions.ListParams.
// 
// The function reads common list query parameters (limit, order, oldest_ref, newest_ref,
// transaction_code, users, statuses, types) and maps them to the corresponding ListParams
// fields. The "payment_types" parameter is converted to a slice of shared.PaymentType.
// Time-related parameters ("changes_since", "oldest_time", "newest_time") are parsed and
// assigned as time pointers when present. It returns an error if the query string cannot
// be parsed.
func parseHrefToListTransactionsParams(href string) (*transactions.ListParams, error) {
	values, err := url.ParseQuery(href)
	if err != nil {
		return nil, err
	}

	paymentTypes := []shared.PaymentType{}

	if pts, exists := values["payment_types"]; exists {
		for _, pt := range pts {
			paymentTypes = append(paymentTypes, shared.PaymentType(pt))
		}
	}

	params := &transactions.ListParams{
		Limit:           getIntPtr(values, "limit"),
		Order:           getStringPtr(values, "order"),
		OldestRef:       getStringPtr(values, "oldest_ref"),
		NewestRef:       getStringPtr(values, "newest_ref"),
		TransactionCode: getStringPtr(values, "transaction_code"),
		Users:           getStringSlice(values, "users"),
		Statuses:        getStringSlice(values, "statuses"),
		Types:           getStringSlice(values, "types"),
		PaymentTypes:    paymentTypes,
	}

	if t := getTimePtr(values, "changes_since"); t != nil {
		params.ChangesSince = t
	}

	if t := getTimePtr(values, "oldest_time"); t != nil {
		params.OldestTime = t
	}

	if t := getTimePtr(values, "newest_time"); t != nil {
		params.NewestTime = t
	}

	return params, nil
}

func (r *Repository) GetTransactionById(transactionId uuid.UUID) (*Transaction, error) {
	transactionIdStr := transactionId.String()
	params := transactions.GetParams{
		ID: &transactionIdStr,
	}

	transaction, err := r.getTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID %s: %w", transactionId, err)
	}

	return transaction, nil
}

func (r *Repository) getTransaction(params transactions.GetParams) (*Transaction, error) {
	transactionResp, err := r.service.Client.Transactions.Get(context.Background(), r.service.MerchantCode, params)
	if err != nil {
		return nil, normalizeSumupError(err)
	}

	return fromSDKTransactionFull(transactionResp), nil
}

func (r *Repository) GetTransactionByClientTransactionId(clientTransactionId uuid.UUID) (*Transaction, error) {
	clientTransactionIdStr := clientTransactionId.String()
	params := transactions.GetParams{
		ClientTransactionID: &clientTransactionIdStr,
	}

	transaction, err := r.getTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ClientTransactionID %s: %w", clientTransactionIdStr, err)
	}

	return transaction, nil
}

func (r *Repository) RefundTransaction(transactionId uuid.UUID) error {
	body := transactions.Refund{}

	err := r.service.Client.Transactions.Refund(context.Background(), transactionId.String(), body)
	if err != nil {
		log.Printf("Error refunding transaction with ID %s: %v", transactionId, err)

		return normalizeSumupError(err)
	}

	return nil
}

// fromSDKTransaction converts an SDK TransactionHistory into a domain Transaction.
// 
// If the SDK TransactionID field is present and parses as a UUID, it is set on the returned Transaction;
// otherwise the TransactionID field will be the zero UUID. All other returned fields are mapped directly
// from the corresponding SDK fields.
func fromSDKTransaction(sdkCheckout *transactions.TransactionHistory) *Transaction {
	var transactionId uuid.UUID

	// Prefer parsing TransactionId if present
	if sdkCheckout.TransactionID != nil {
		if parsedId, err := uuid.Parse(string(*sdkCheckout.TransactionID)); err == nil {
			transactionId = parsedId
		}
	}

	return &Transaction{
		ID:              string(*sdkCheckout.ID),
		TransactionCode: string(*sdkCheckout.TransactionCode),
		TransactionID:   transactionId,
		Amount:          utils.F32PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		CardType:        string(*sdkCheckout.CardType),
		CreatedAt:       utils.TimePtr(sdkCheckout.Timestamp),
		Status:          string(*sdkCheckout.Status),
	}
}

// fromSDKTransactionFull converts an SDK TransactionFull into a domain Transaction.
// It parses the SDK ID into a UUID (zero UUID on parse failure or if missing), maps amount, currency, status, timestamp, and card type, and converts SDK events into domain TransactionEvent values; nil SDK slices and fields produce empty or zero-valued domain fields.
func fromSDKTransactionFull(sdkCheckout *transactions.TransactionFull) *Transaction {
	var transactionId uuid.UUID

	if sdkCheckout.ID != nil {
		parsedId, err := uuid.Parse(*sdkCheckout.ID)
		if err == nil {
			transactionId = parsedId
		}
	}

	var events []TransactionEvent
	if sdkCheckout.Events != nil {
		events = make([]TransactionEvent, 0, len(sdkCheckout.Events))
		for _, sdkEvent := range sdkCheckout.Events {
			events = append(events, fromSDKTransactionEvent(&sdkEvent))
		}
	} else {
		events = make([]TransactionEvent, 0)
	}

	var cardType string
	if sdkCheckout.Card != nil && sdkCheckout.Card.Type != nil {
		cardType = string(*sdkCheckout.Card.Type)
	}

	return &Transaction{
		ID:              transactionId.String(),
		TransactionCode: string(*sdkCheckout.TransactionCode),
		TransactionID:   transactionId,
		Amount:          utils.F32PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		CardType:        cardType,
		CreatedAt:       utils.TimePtr(sdkCheckout.Timestamp),
		Events:          events,
		Status:          string(*sdkCheckout.Status),
	}
}

// fromSDKTransactionEvent converts an SDK transactions.Event into a domain TransactionEvent.
// It parses the event Timestamp using RFC3339; on parse failure the error is logged and the zero time is used.
// If Amount is nil it is treated as 0. ID, Type, and Status are extracted from the corresponding SDK fields.
func fromSDKTransactionEvent(sdkEvent *transactions.Event) TransactionEvent {
	timestamp := time.Time{}

	if sdkEvent.Timestamp != nil {
		var err error
		if sdkEvent.Timestamp != nil {
			timestamp, err = time.Parse(time.RFC3339, string(*sdkEvent.Timestamp))
			if err != nil {
				log.Printf("Error parsing timestamp %s: %v", *sdkEvent.Timestamp, err)
			}
		}
	}

	amount := float64(0)
	if sdkEvent.Amount != nil {
		amount = float64(*sdkEvent.Amount)
	}

	return TransactionEvent{
		ID:        int(*sdkEvent.ID),
		Timestamp: timestamp,
		Type:      string(*sdkEvent.Type),
		Amount:    decimal.NewFromFloat(amount),
		Status:    string(*sdkEvent.Status),
	}
}