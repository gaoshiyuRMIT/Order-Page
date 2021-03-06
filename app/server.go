package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"app/models"
	"app/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Orders API is alive!")
}


func searchOrders(w http.ResponseWriter, r *http.Request) {
	queryVals := r.URL.Query()
	var orderQuery models.OrderInfoQuery
	pagination := models.NewPagination()
	schema.NewDecoder().Decode(&orderQuery, queryVals)
	schema.NewDecoder().Decode(pagination, queryVals)
	
	orderMgr := models.NewOrderManager(utils.Config)
	orders, err := orderMgr.Search(&orderQuery, pagination)
	if err != nil {
		utils.WriteHttpError(w, err)
		return
	}
	cstmrMgr, err := models.NewCustomerManager(utils.Config)
	if err != nil {
		utils.WriteHttpError(w, err)
		return
	}
	err = cstmrMgr.FillInCustomerInfo(orders)
	if err != nil {
		utils.WriteHttpError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(orders)
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Path("/api/heartbeat").HandlerFunc(heartbeat).Methods("GET")
	myRouter.Path("/api/orders/search").HandlerFunc(searchOrders).Methods("GET")
	port := utils.Config.GetAPIPort()
	log.Fatal(http.ListenAndServe(":" + port, myRouter))
}

func main() {
	handleRequest()
}
