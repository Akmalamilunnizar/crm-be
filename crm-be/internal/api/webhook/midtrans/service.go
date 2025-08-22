package midtrans

import (
	"skripsi-be/internal/models/entities"
	"strconv"
)

type WebhookMidtransServiceInterface interface {
	GetAllWebhookMidtransService() ([]entities.Areas, error)
	ReceivedWebhookMidtransService(request *MidtransNotification) error
}
type WebhookMidtransServiceStruct struct {
	repository WebhookMidtransRepositoryInterface
}

func NewWebhookMidtransService(repository WebhookMidtransRepositoryInterface) *WebhookMidtransServiceStruct {
	return &WebhookMidtransServiceStruct{repository}
}

func (s WebhookMidtransServiceStruct) GetAllWebhookMidtransService() ([]entities.Areas, error) {
	areas, err := s.repository.FindWebhookMidtransRepository()

	if err != nil {
		return areas, err
	}

	return areas, err
}

func (s WebhookMidtransServiceStruct) ReceivedWebhookMidtransService(request *MidtransNotification) error {
	_, err := s.repository.UpdatePaidWebhookMidtransRepository(request.OrderID)
	if err != nil {
		return err
	}

	requestTransaction := CreateWebhookMidtransRequest{}
	requestTransaction.AccountID = "b82074e7-7acb-40e2-ab33-014f4b09c1f8"
	requestTransaction.TypeCash = string(entities.TransactionsTypeCashInternet)
	requestTransaction.TypeInOut = string(entities.TransactionsTypeInOutIn)
	amount, err := strconv.ParseFloat(request.GrossAmount, 64)
	if err != nil {
		return err
	}
	requestTransaction.Amount = int64(amount)
	_, err = s.repository.CreateWebhookMidtransRepository(requestTransaction)

	if err != nil {
		return err
	}
	return nil
}
