package models

import (
	"context"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"app/utils"
)

type CustomerManager struct {
	Client *mongo.Client
	DB *mongo.Database
}

func NewCustomerManager(cfg *utils.ConfigReader) *CustomerManager {
	cm := &CustomerManager{}
	cm.Client, cm.DB = cfg.GetMongoDB()
	return cm
}

func (mgr *CustomerManager) GetCustomerCollection() *mongo.Collection {
	return mgr.DB.Collection("customers")
}

func (mgr *CustomerManager) GetCustomerCompanyCollection() *mongo.Collection {
	return mgr.DB.Collection("customer_companies")
}

func (mgr *CustomerManager) FillInCustomerInfo(orderInfos []OrderInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := mgr.Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mgr.Client.Disconnect(ctx)

	customer_col := mgr.GetCustomerCollection()
	company_col := mgr.GetCustomerCompanyCollection()

	for i, oi := range orderInfos {
		var customerInfo CustomerInfo
		err := customer_col.FindOne(ctx, bson.M{"user_id": oi.CustomerID}).Decode(&customerInfo)
		if err != nil {
			log.Fatal(err)
		}
		err = company_col.FindOne(ctx, bson.M{"company_id": customerInfo.CompanyID}).Decode(&customerInfo)
		if err != nil {
			log.Fatal(err)
		}
		orderInfos[i].SetCustomerInfo(customerInfo)
	}
}