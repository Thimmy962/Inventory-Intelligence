package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/analytics"
	"main/internal/app"
	"main/internal/database"
	"net/http"
	"strconv"
)

// called on bad request
func ProcessingError(w http.ResponseWriter, code int, Err error) {    
    w.Header().Set("Content-Type", "application/json") // Fixed typo
    w.WriteHeader(code)
    response := map[string]string{"error": Err.Error()}
    _ = json.NewEncoder(w).Encode(response) 
}



// handles both writing a single struct or a list of structs to responseWriter
// all thanks to NewEncoder

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
    var buf bytes.Buffer

    json.NewEncoder(&buf).Encode(payload)



    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
    
    // 3. If the code above didn't return, this is the FIRST 
    // and ONLY time WriteHeader should be called.
    w.WriteHeader(code) 

    w.Write(buf.Bytes())
}

func (db *Handler)NewProduct(purchase *PurAdj, w http.ResponseWriter, req *http.Request) {
	id := purchase.ProductID
	quantityAdded := purchase.Quantity
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
			return nil, 0, fmt.Errorf("Quantity on hand for '%s' is less than request to buy", product.ProductName)
		}

		//price * quantity
		total = total + (product.Price * float64(item.Quantity_sold))
		
		// append to the list
		products = append(products, product)
	}
	return products, total, nil
}


func (handler *Handler)StartWorker(ctx context.Context, workers int) {
	log.Println("starting Workers")
	log.Printf("%d Workers Starting\n", workers)

	for i:=0; i < workers; i++ {
		for salesID := range handler.channel {
			handler.wg.Add(1)
			sales_id, _ := strconv.Atoi(salesID)
			analytics.Analytics(ctx, handler.server.Queries, int32(sales_id))
			handler.wg.Done()
		}
	}
}

