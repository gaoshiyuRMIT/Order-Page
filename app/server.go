package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"app/models"
	"app/utils"

	"github.com/gorilla/mux"
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

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/heartbeat", heartbeat)
	myRouter.HandleFunc("/api/orders", orders)
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {
	handleRequest()
}
