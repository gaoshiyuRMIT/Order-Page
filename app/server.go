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

func orders(w http.ResponseWriter, r *http.Request) {
	orderMgr := models.NewOrderManager(utils.Config)
	orders := orderMgr.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func searchOrders(w http.ResponseWriter, r *http.Request) {
	queryVals := r.URL.Query()
	var orderQuery models.OrderInfoQuery
	schema.NewDecoder().Decode(&orderQuery, queryVals)
	
	orderMgr := models.NewOrderManager(utils.Config)
	orders := orderMgr.Search(&orderQuery)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Path("/api/heartbeat").HandlerFunc(heartbeat).Methods("GET")
	myRouter.Path("/api/orders").HandlerFunc(orders).Methods("GET")
	myRouter.Path("/api/orders/search").HandlerFunc(searchOrders).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {
	handleRequest()
}
