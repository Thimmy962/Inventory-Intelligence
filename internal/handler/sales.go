package handler

import (
	"encoding/json"
	"log"
	"main/internal/database"
	"net/http"
)

// any given sales
type Sales struct {
	id int64
	TotalAmount float32 `json:"amount"`
}

// item on a sales
type Sales_Item struct {
	id int64
	sales_id int64
	Product_id string `json:"product_id"`
	Quantity_sold int32 `json:"quantity"`
} 

type Adjustment struct {
	id int16
	ProductID string `json:"productID"`
	Quantity int16 `json:"quantity"`
	Reason string `json:"reason"`
}

func (db *Handler) CreateSales(wr http.ResponseWriter, req *http.Request) {
	var items []Sales_Item // list of items under a sales
	err := json.NewDecoder(req.Body).Decode(&items);
	if err != nil {
		ProcessingError(wr, http.StatusBadRequest, err)
		return
	}

	/*
		Sort all the sales item into a list of format recognized by the system
		if theres an error return
		all get the sum total of the products as sales
	*/
	products, total, err := validateSale(req, items, *db.server)
	if err != nil {
		ProcessingError(wr, http.StatusBadRequest, err)
		return
	}

	// creates the sales that contains all the sold items; returns an error(if available) and the id of the newly created sales
	sales_id, err := db.server.Queries.CreateSales(req.Context(), total)
	if err != nil {
		ProcessingError(wr, http.StatusBadRequest, err)
		return
	}

	for index, product := range products {
		analyticsData, err := db.server.Queries.CreateSalesItems(req.Context(), database.CreateSalesItemsParams{
			SalesID: sales_id, ProductID: product.ID, PriceAtSale: product.Price, QuantitySold: items[index].Quantity_sold,
		})
		// if there was error in processing any item delete every other processed items and the sales it self
		if err != nil {
			ProcessingError(wr, http.StatusBadRequest, err)
			// delete every other sales_item under this sales
			db.server.Queries.DeleteSalesItems(req.Context(), sales_id)

			// delete the sales itself
			db.server.Queries.DeleteSales(req.Context(), sales_id)
			return
		}

		db.channel <- analyticsData
	}

	for index, product := range products {
		//quantity left = quantity before the sale - quantity sold
		quantity_left := product.QuantityOnHand - items[index].Quantity_sold
		db.server.Queries.UpdatedInventory(req.Context(), 
		database.UpdatedInventoryParams{ProductID: product.ID, QuantityOnHand: quantity_left})
	}
	respondWithJSON(wr, http.StatusCreated, map[string]string{
		"status": "Sales Successful",
	})
}


func (db *Handler) CreateAdjustment(wr http.ResponseWriter, req *http.Request) {
	var adjustment Adjustment
	err := json.NewDecoder(req.Body).Decode(&adjustment)
	if err != nil {
		log.Println(err)
		ProcessingError(wr, http.StatusBadRequest, err)
		return
	}

	// create a row for this adjustment
	id, err := db.server.Queries.CreateAdjustment(req.Context(), database.CreateAdjustmentParams{
		ProductID: adjustment.ProductID, QuantityChanged: int32(adjustment.Quantity), Reason: adjustment.Reason,})
	
	// if there is an error while creating process error and return
	if err != nil {
		log.Println(err)
		ProcessingError(wr, http.StatusBadRequest, err)
		return
	}

	// if there is an error while retrieving inventory to be updated delete the adjustment
	inventory, err := db.server.Queries.GetInventory(req.Context(), adjustment.ProductID)
	if err != nil {
		log.Println(err)
		ProcessingError(wr, http.StatusBadRequest, err)

		db.server.Queries.DeleteAdjustment(req.Context(),id)
		return

	}
	var newVal int32 = inventory.QuantityOnHand - int32(adjustment.Quantity)
	if newVal < 0 {
		newVal = 0
	}

	db.server.Queries.UpdatedInventory(req.Context(), 
	database.UpdatedInventoryParams{ProductID: inventory.ProductID, QuantityOnHand: newVal})

	respondWithJSON(wr, http.StatusCreated, map[string]string {
		"status": "Adjustment Created",
	})
}
