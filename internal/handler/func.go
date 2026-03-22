package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/app"
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
	
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(buf)
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

// interface context.Context
func validateSale(req *http.Request, items []Sales_Item, query app.Server) ([]database.GetProductInventoryRow, float64, error) {
	var products []database.GetProductInventoryRow
	var total  = 0.0 // sum of all the sales item
	for _, item := range items {
		product, err := query.Queries.GetProductInventory(req.Context(), item.Product_id)
		if err != nil {
			return nil, 0, err
		} else if (product.QuantityOnHand < item.Quantity_sold){
			return nil, 0, fmt.Errorf("Quantity on hand of %s is less than request to buy", product.ProductName)
		}

		//price * quantity
		total = total + (product.Price * float64(item.Quantity_sold))
		
		// append to the list
		products = append(products, product)
	}
	return products, total, nil
}