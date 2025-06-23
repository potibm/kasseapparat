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

	params := transactions.ListTransactionsV21Params{
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

		if resp.Items == nil || len(*resp.Items) == 0 {
			return allItems, nil
		}

		allItems = append(allItems, ptrSliceToSlice(resp.Items)...)
		sortTransactionsByCreatedAt(allItems)

		nextHref := findNextHref(resp.Links)
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

func parseHrefToListTransactionsParams(href string) (*transactions.ListTransactionsV21Params, error) {
	values, err := url.ParseQuery(href)
	if err != nil {
		return nil, err
	}

	params := &transactions.ListTransactionsV21Params{
		Limit:           getIntPtr(values, "limit"),
		Order:           getStringPtr(values, "order"),
		OldestRef:       getStringPtr(values, "oldest_ref"),
		NewestRef:       getStringPtr(values, "newest_ref"),
		TransactionCode: getStringPtr(values, "transaction_code"),
		Users:           getStringSlicePtr(values, "users"),
		Statuses:        getStringSlicePtr(values, "statuses"),
		Types:           getStringSlicePtr(values, "types"),
		PaymentTypes:    getStringSlicePtr(values, "payment_types"),
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
	params := transactions.GetTransactionV21Params{
		Id: &transactionIdStr,
	}

	transaction, err := r.getTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID %s: %w", transactionId, err)
	}

	return transaction, nil
}

func (r *Repository) getTransaction(params transactions.GetTransactionV21Params) (*Transaction, error) {
	transactionResp, err := r.service.Client.Transactions.Get(context.Background(), r.service.MerchantCode, params)
	if err != nil {
		return nil, normalizeSumupError(err)
	}

	return fromSDKTransactionFull(transactionResp), nil
}

func (r *Repository) GetTransactionByClientTransactionId(clientTransactionId uuid.UUID) (*Transaction, error) {
	clientTransactionIdStr := clientTransactionId.String()
	params := transactions.GetTransactionV21Params{
		ClientTransactionId: &clientTransactionIdStr,
	}

	transaction, err := r.getTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ClientTransactionID %s: %w", clientTransactionIdStr, err)
	}

	return transaction, nil
}

func (r *Repository) RefundTransaction(transactionId uuid.UUID) error {
	body := transactions.RefundTransactionBody{}

	err := r.service.Client.Transactions.Refund(context.Background(), transactionId.String(), body)
	if err != nil {
		log.Printf("Error refunding transaction with ID %s: %v", transactionId, err)

		return normalizeSumupError(err)
	}

	return nil
}

func fromSDKTransaction(sdkCheckout *transactions.TransactionHistory) *Transaction {
	var transactionId uuid.UUID

	// Prefer parsing TransactionId if present
	if sdkCheckout.TransactionId != nil {
		if parsedId, err := uuid.Parse(string(*sdkCheckout.TransactionId)); err == nil {
			transactionId = parsedId
		}
	}

	return &Transaction{
		ID:              string(*sdkCheckout.Id),
		TransactionCode: string(*sdkCheckout.TransactionCode),
		TransactionID:   transactionId,
		Amount:          utils.F64PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		CardType:        string(*sdkCheckout.CardType),
		CreatedAt:       utils.TimePtr(sdkCheckout.Timestamp),
		Status:          string(*sdkCheckout.Status),
	}
}

func fromSDKTransactionFull(sdkCheckout *transactions.TransactionFull) *Transaction {
	var transactionId uuid.UUID

	if sdkCheckout.Id != nil {
		parsedId, err := uuid.Parse(*sdkCheckout.Id)
		if err == nil {
			transactionId = parsedId
		}
	}

	var events []TransactionEvent
	if sdkCheckout.Events != nil {
		events = make([]TransactionEvent, 0, len(*sdkCheckout.Events))
		for _, sdkEvent := range *sdkCheckout.Events {
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
		Amount:          utils.F64PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		CardType:        cardType,
		CreatedAt:       utils.TimePtr(sdkCheckout.Timestamp),
		Events:          events,
		Status:          string(*sdkCheckout.Status),
	}
}

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
		ID:        int(*sdkEvent.Id),
		Timestamp: timestamp,
		Type:      string(*sdkEvent.Type),
		Amount:    decimal.NewFromFloat(amount),
		Status:    string(*sdkEvent.Status),
	}
}
