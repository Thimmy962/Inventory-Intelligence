package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internal/database"
	"net/http"
	"time"
)


type Product struct {
	ID           string `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Price        float64 `json:"price"`
	ReorderLevel int32 `json:"reorder_level"`
	ProductName  string `json:"name"`
}

type Purchase struct {
	ID string `json:"id`
	ProductID string `json:"productId`
	QuantityAdded int32 `json:"quantity"`
}

func (db * Server)CreateProduct(writer http.ResponseWriter, req *http.Request) {
    var product Product

    err := json.NewDecoder(req.Body).Decode(&product)
    if err != nil {
        ProcessingError(writer, 400, err)
        return
    }

	created_product, err := db.queries.CreateProduct(req.Context(), database.CreateProductParams{
		ProductName: product.ProductName, Price: product.Price, ReorderLevel: product.ReorderLevel,
	})
	if err != nil {
			log.Printf("DB error value: %s", err.Error())
			ProcessingError(writer, 400, err)
			return
		}
	respondWithJSON(writer, 201, created_product)
}

func (db *Server)GetProducts(writer http.ResponseWriter, req *http.Request) {
	products, err := db.queries.GetProducts(req.Context())
	if err != nil {
			log.Printf("DB error value: %s", err.Error())
			ProcessingError(writer, 400, err)
			return
		}
	respondWithJSON(writer, 201, products)
}


func (db *Server) BulkCreateProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product

	err := json.NewDecoder(r.Body).Decode(&products)
	if err != nil {
		ProcessingError(w, 400, err)
		return
	}

	for _, p := range products {
		_, err := db.queries.CreateProduct(
			r.Context(),
			database.CreateProductParams{
				ProductName:  p.ProductName,
				Price:        p.Price,
				ReorderLevel: p.ReorderLevel,
			},
		)

		if err != nil {
			ProcessingError(w, 400, err)
			return
		}
	}

	respondWithJSON(w, 201, map[string]string{
		"status": "products inserted",
	})
}

func (db *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := db.queries.GetProduct(r.Context(), id)
	if err != nil {
		ProcessingError(w, http.StatusNotFound, err)
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}


func (db *Server)NewBulkPurchase(w http.ResponseWriter, r *http.Request) {
	var purchases []Purchase
	err := json.NewDecoder(r.Body).Decode(&purchases)
	if err != nil{
		log.Println(err)
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	for _, purchase := range purchases {
		err = db.queries.CreatePurchase(r.Context(),
			database.CreatePurchaseParams{ProductID: purchase.ProductID,
			QuantityAdded: purchase.QuantityAdded})
		if err != nil{
			log.Println(err)
			ProcessingError(w, http.StatusBadRequest, err)
			return
		}
		if db.Inventory(&purchase, w, r) != nil {
			err = db.queries.DeleteProduct(r.Context(), purchase.ID)
			ProcessingError(w, http.StatusBadRequest, fmt.Errorf("Could not delete product after failed update of inventory; Delete manualy"))
			return
		}
	}
	
	respondWithJSON(w, http.StatusCreated, map[string]string{
		"status": "purchase inserted",
	})	
}

func (db *Server)NewPurchase(w http.ResponseWriter, r *http.Request) {
	var purchase Purchase
	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil{
		log.Println(err)
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	err = db.queries.CreatePurchase(r.Context(),
			database.CreatePurchaseParams{ProductID: purchase.ProductID,
			QuantityAdded: purchase.QuantityAdded})
		if err != nil{
			log.Println(err)
			ProcessingError(w, http.StatusBadRequest, err)
			return
		}
		if db.Inventory(&purchase, w, r) != nil {
			err = db.queries.DeleteProduct(r.Context(), purchase.ID)
			ProcessingError(w, http.StatusBadRequest, fmt.Errorf("Could not delete product after failed update of inventory; Delete manualy"))
			return
		}
	
	respondWithJSON(w, http.StatusCreated, map[string]string{
		"status": "purchase inserted",
	})	
}


func (db *Server)Inventory(purchase *Purchase, w http.ResponseWriter, req *http.Request) error {
	id := purchase.ProductID
	quantityAdded := purchase.QuantityAdded
	data, err := db.queries.GetInventory(req.Context(), id)
	if err != nil {
		newErr := db.queries.NewInventory(req.Context(), database.NewInventoryParams{ProductID: id, QuantityOnHand: quantityAdded})
		if newErr != nil {
			ProcessingError(w, http.StatusBadRequest, newErr)
			return newErr
		}
	}
	return db.queries.UpdatedInventory(req.Context(), database.UpdatedInventoryParams{ProductID: id, 
		QuantityOnHand: quantityAdded + data.QuantityOnHand})

}
