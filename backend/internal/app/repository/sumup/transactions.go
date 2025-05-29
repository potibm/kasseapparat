package sumup

import (
	"context"
	"log"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/sumup/sumup-go/transactions"
)

func (r *Repository) GetTransactions() ([]Transaction, error) {
	result := []Transaction{}
	params := transactions.ListTransactionsV21Params{}

	transactionResp, err := r.service.Client.Transactions.List(context.Background(), r.service.MerchantCode, params)
	if err != nil {
		return nil, err
	}

	// iterate over the checkouts and convert them to our Checkout type
	for _, transaction := range *transactionResp.Items {
		result = append(result, *fromSDKTransaction(&transaction))
	}

	return result, nil
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
