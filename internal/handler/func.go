package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/database"
	"net/http"
)

// called on bad request
func ProcessingError(w http.ResponseWriter, code int, Err error) {
	w.Header().Set("Content-Type", "apllication/json")
	w.WriteHeader(code)
	response := map[string]string{"error": Err.Error()}
	json.NewEncoder(w).Encode(response)
}



// handles both writing a single struct or a list of structs to responseWriter
// all thanks to NewEncoder
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		log.Println(err.Error())
		ProcessingError(w, http.StatusInternalServerError, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(code)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println(err.Error())
	}
}

/*
 * if product id of the current purchase is not in inventory
 * create an inventory for the product id set quantity_on_hand to purchase quantity
 * else update the inventory  quantity_on_hand to  quantity_on_hand + purchase quantity
 *
*/

func (db *Handler)NewProduct(purchase *Purchase, w http.ResponseWriter, req *http.Request) {
	id := purchase.ProductID
	quantityAdded := purchase.QuantityAdded
	data, err := db.server.Queries.GetInventory(req.Context(), id)
	if err != nil {
		newErr := db.server.Queries.NewInventory(req.Context(), database.NewInventoryParams{ProductID: id, QuantityOnHand: quantityAdded})
		if newErr != nil {
			ProcessingError(w, http.StatusBadRequest, newErr)
			return
		}
	}
	err = db.server.Queries.UpdatedInventory(req.Context(), database.UpdatedInventoryParams{ProductID: id, 
		QuantityOnHand: quantityAdded + data.QuantityOnHand})

}