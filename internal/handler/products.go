package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internal/app"
	"main/internal/database"
	"net/http"
	"time"
)

type Handler struct {
	server *app.Server
}

func NewHandler(q *app.Server) *Handler {
	return &Handler{
		server: q,
	}
}

type Product struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Price        float64   `json:"price"`
	ReorderLevel int32     `json:"reorder_level"`
	ProductName  string    `json:"name"`
}

type Purchase struct {
	ID            string `json:"id"`
	ProductID     string `json:"productId"`
	QuantityAdded int32  `json:"quantity"`
}

func (db *Handler) CreateProduct(writer http.ResponseWriter, req *http.Request) {
	var product Product

	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		go log.Println(err)
		go ProcessingError(writer, 400, err)
		return
	}

	created_product, err := db.server.Queries.CreateProduct(req.Context(), database.CreateProductParams{
		ProductName: product.ProductName, Price: product.Price, ReorderLevel: product.ReorderLevel,
	})
	if err != nil {
		go log.Printf("DB error value: %s", err.Error())
		go ProcessingError(writer, 400, err)
		return
	}
	go respondWithJSON(writer, 201, created_product)
}

func (db *Handler) GetProducts(writer http.ResponseWriter, req *http.Request) {
	products, err := db.server.Queries.GetProducts(req.Context())
	if err != nil {
		go log.Printf("DB error value: %s", err.Error())
		go ProcessingError(writer, 400, err)
		return
	}
	go respondWithJSON(writer, 200, products)
}

func (db *Handler) BulkCreateProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product

	err := json.NewDecoder(r.Body).Decode(&products)
	if err != nil {
		go log.Println(err)
		go ProcessingError(w, 400, err)
		return
	}

	for _, p := range products {
		_, err := db.server.Queries.CreateProduct(
			r.Context(),
			database.CreateProductParams{
				ProductName:  p.ProductName,
				Price:        p.Price,
				ReorderLevel: p.ReorderLevel,
			},
		)

		if err != nil {
			go log.Println(err)
			go ProcessingError(w, 400, err)
			return
		}
	}

	go respondWithJSON(w, 201, map[string]string{
		"status": "products inserted",
	})
}

func (db *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := db.server.Queries.GetProduct(r.Context(), id)
	if err != nil {
		go log.Println(err)
		go ProcessingError(w, http.StatusNotFound, err)
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (db *Handler) NewBulkPurchase(w http.ResponseWriter, r *http.Request) {
	var purchases []Purchase
	err := json.NewDecoder(r.Body).Decode(&purchases)
	if err != nil {
		go log.Println(err)
		go ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	for _, purchase := range purchases {
		err = db.server.Queries.CreatePurchase(r.Context(),
			database.CreatePurchaseParams{ProductID: purchase.ProductID,
				QuantityAdded: purchase.QuantityAdded})
		if err != nil {
			go log.Println(err)
			go ProcessingError(w, http.StatusBadRequest, err)
			return
		}
		if db.Inventory(&purchase, w, r) != nil {
			err = db.server.Queries.DeleteProduct(r.Context(), purchase.ID)
			go ProcessingError(w, http.StatusBadRequest, fmt.Errorf("Could not delete product after failed update of inventory; Delete manualy"))
			return
		}
	}

	go respondWithJSON(w, http.StatusCreated, map[string]string{
		"status": "purchase inserted",
	})
}


// creates a new purchase.
// If the new purchase is in inventory update to the quantity available else create a new inventory
func (db *Handler) NewPurchase(w http.ResponseWriter, r *http.Request) {
	var purchase Purchase
	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil {
		go log.Println(err)
		go ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	//creates new purchase
	err = db.server.Queries.CreatePurchase(r.Context(),
		database.CreatePurchaseParams{ProductID: purchase.ProductID,
			QuantityAdded: purchase.QuantityAdded})

	// if there is an error in creating new purchase return
	if err != nil {
		go log.Println(err)
		go ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	// if new inventory or updating inventory fails delete the purchase 
	if db.Inventory(&purchase, w, r) != nil {
		err = db.server.Queries.DeleteProduct(r.Context(), purchase.ID)
		ProcessingError(w, http.StatusBadRequest, fmt.Errorf("Could not delete product after failed update of inventory; Delete manually"))
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"status": "purchase inserted",
	})
}


// creates new or update available inventory 
func (db *Handler) Inventory(purchase *Purchase, w http.ResponseWriter, req *http.Request) error {
	id := purchase.ProductID
	quantityAdded := purchase.QuantityAdded

	data, err := db.server.Queries.GetInventory(req.Context(), id)

	// if inventory dies not exists create a new one
	if err != nil {
		newErr := db.server.Queries.NewInventory(req.Context(), database.NewInventoryParams{ProductID: id, QuantityOnHand: quantityAdded})
		
		// if creatting a new inventory fails return
		if newErr != nil {
			go log.Println(err)
			go ProcessingError(w, http.StatusBadRequest, newErr)

			return newErr
		}
	}
	// try to update inventory if it exists
	return db.server.Queries.UpdatedInventory(req.Context(), database.UpdatedInventoryParams{ProductID: id,
		QuantityOnHand: quantityAdded + data.QuantityOnHand})
}
