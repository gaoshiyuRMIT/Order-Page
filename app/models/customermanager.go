package models

import (
	"context"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"app/utils"
)

type CustomerManager struct {
	Client *mongo.Client
	DB *mongo.Database
}

func NewCustomerManager(cfg *utils.ConfigReader) (*CustomerManager, error) {
	cm := &CustomerManager{}
	var err error
	cm.Client, cm.DB, err = cfg.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (mgr *CustomerManager) GetCustomerCollection() *mongo.Collection {
	return mgr.DB.Collection("customers")
}

func (mgr *CustomerManager) GetCustomerCompanyCollection() *mongo.Collection {
	return mgr.DB.Collection("customer_companies")
}

func (mgr *CustomerManager) FillInCustomerInfo(orderInfos []OrderInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := mgr.Client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Cannot connect Mongo client to context. %w", err)
	}
	defer mgr.Client.Disconnect(ctx)

	customer_col := mgr.GetCustomerCollection()
	company_col := mgr.GetCustomerCompanyCollection()

	for i, oi := range orderInfos {
		var customerInfo CustomerInfo
		err := customer_col.FindOne(ctx, bson.M{"user_id": oi.CustomerID}).Decode(&customerInfo)
		if err != nil {
			return fmt.Errorf("Cannot query for one customer with customer id %s. %w", oi.CustomerID, err)
		}
		err = company_col.FindOne(ctx, bson.M{"company_id": customerInfo.CompanyID}).Decode(&customerInfo)
		if err != nil {
			return fmt.Errorf("Cannot query for one customer with company id %d. %w", customerInfo.CompanyID, err)
		}
		orderInfos[i].SetCustomerInfo(customerInfo)
	}
	return nil
}