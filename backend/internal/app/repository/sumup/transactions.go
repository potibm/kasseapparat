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
	sumup "github.com/sumup/sumup-go"
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

func (r *Repository) fetchPagedTransactions(ctx context.Context, oldestFrom *time.Time, maxPages, pageSize int) ([]*sumup.TransactionHistory, error) {
	var allItems []*sumup.TransactionHistory

	pageCount := 0

	params := sumup.TransactionsListParams{
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

func sortTransactionsByCreatedAt(transactions []*sumup.TransactionHistory) {
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

func ptrSliceToSlice(ptrSlice *[]sumup.TransactionHistory) []*sumup.TransactionHistory {
	if ptrSlice == nil {
		return nil
	}

	out := make([]*sumup.TransactionHistory, len(*ptrSlice))
	for i := range *ptrSlice {
		out[i] = &(*ptrSlice)[i]
	}

	return out
}

func findNextHref(links *[]sumup.Link) string {
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

func parseHrefToListTransactionsParams(href string) (*sumup.TransactionsListParams, error) {
	values, err := url.ParseQuery(href)
	if err != nil {
		return nil, err
	}

	paymentTypes := []sumup.PaymentType{}

	if pts, exists := values["payment_types"]; exists {
		for _, pt := range pts {
			paymentTypes = append(paymentTypes, sumup.PaymentType(pt))
		}
	}

	params := &sumup.TransactionsListParams{
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
	params := sumup.TransactionsGetParams{
		ID: &transactionIdStr,
	}

	transaction, err := r.getTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID %s: %w", transactionId, err)
	}

	return transaction, nil
}

func (r *Repository) getTransaction(params sumup.TransactionsGetParams) (*Transaction, error) {
	transactionResp, err := r.service.Client.Transactions.Get(context.Background(), r.service.MerchantCode, params)
	if err != nil {
		return nil, normalizeSumupError(err)
	}

	return fromSDKTransactionFull(transactionResp), nil
}

func (r *Repository) GetTransactionByClientTransactionId(clientTransactionId uuid.UUID) (*Transaction, error) {
	clientTransactionIdStr := clientTransactionId.String()
	params := sumup.TransactionsGetParams{
		ClientTransactionID: &clientTransactionIdStr,
	}

	transaction, err := r.getTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ClientTransactionID %s: %w", clientTransactionIdStr, err)
	}

	return transaction, nil
}

func (r *Repository) RefundTransaction(transactionId uuid.UUID) error {
	body := sumup.TransactionsRefundParams{}

	err := r.service.Client.Transactions.Refund(context.Background(), transactionId.String(), body)
	if err != nil {
		log.Printf("Error refunding transaction with ID %s: %v", transactionId, err)

		return normalizeSumupError(err)
	}

	return nil
}

func fromSDKTransaction(sdkCheckout *sumup.TransactionHistory) *Transaction {
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

func fromSDKTransactionFull(sdkCheckout *sumup.TransactionFull) *Transaction {
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

func fromSDKTransactionEvent(sdkEvent *sumup.Event) TransactionEvent {
	timestamp := time.Time{}

	if sdkEvent.Timestamp != nil {
		timestamp = *sdkEvent.Timestamp
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
