package sumup

import (
	"context"
	"log"
	"net/url"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/sumup/sumup-go/transactions"
)

const (
	TransactionPageSize = 100
	TransactionMaxPages = 3
)

func (r *Repository) GetTransactions() ([]Transaction, error) {
	ctx := context.Background()

	sdkItems, err := r.fetchPagedTransactions(ctx, TransactionMaxPages, TransactionPageSize)
	if err != nil {
		return nil, err
	}

	result := make([]Transaction, 0, len(sdkItems))
	for _, sdkTx := range sdkItems {
		result = append(result, *fromSDKTransaction(sdkTx))
	}

	return result, nil
}

func (r *Repository) fetchPagedTransactions(ctx context.Context, maxPages, pageSize int) ([]*transactions.TransactionHistory, error) {
	var allItems []*transactions.TransactionHistory

	pageCount := 0

	params := transactions.ListTransactionsV21Params{
		Limit: &pageSize,
	}

	for {
		if pageCount >= maxPages {
			return allItems, nil
		}

		pageCount++

		log.Printf("Fetching transactions, page %d with limit %d", pageCount, pageSize)

		resp, err := r.service.Client.Transactions.List(ctx, r.service.MerchantCode, params)
		if err != nil {
			return nil, err
		}

		if resp.Items == nil || len(*resp.Items) == 0 {
			return allItems, nil
		}

		allItems = append(allItems, ptrSliceToSlice(resp.Items)...)

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

func (r *Repository) GetTransactionById(transactionId string) (*Transaction, error) {
	params := transactions.GetTransactionV21Params{Id: &transactionId}

	return r.getTransaction(params)
}

func (r *Repository) getTransaction(params transactions.GetTransactionV21Params) (*Transaction, error) {
	transactionResp, err := r.service.Client.Transactions.Get(context.Background(), r.service.MerchantCode, params)
	if err != nil {
		log.Printf("Error retrieving transaction: %v", err)

		return nil, err
	}

	return fromSDKTransactionFull(transactionResp), nil
}

func (r *Repository) RefundTransaction(transactionId string) error {
	body := transactions.RefundTransactionBody{}

	_, err := r.service.Client.Transactions.Refund(context.Background(), transactionId, body)
	if err != nil {
		log.Printf("Error retrieving checkout with ID %s: %v", transactionId, err)

		return err
	}

	return nil
}

func fromSDKTransaction(sdkCheckout *transactions.TransactionHistory) *Transaction {
	return &Transaction{
		ID:              utils.StrPtr(sdkCheckout.Id),
		TransactionCode: string(*sdkCheckout.TransactionCode),
		Amount:          utils.F64PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		CreatedAt:       utils.TimePtr(sdkCheckout.Timestamp),
		Status:          string(*sdkCheckout.Status),
	}
}

func fromSDKTransactionFull(sdkCheckout *transactions.TransactionFull) *Transaction {
	return &Transaction{
		ID:              utils.StrPtr(sdkCheckout.Id),
		TransactionCode: string(*sdkCheckout.TransactionCode),
		Amount:          utils.F64PtrToDecimal(sdkCheckout.Amount),
		Currency:        string(*sdkCheckout.Currency),
		CreatedAt:       utils.TimePtr(sdkCheckout.Timestamp),
		Status:          string(*sdkCheckout.Status),
	}
}
