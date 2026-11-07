package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

//updated product struct to
//fetch products by id

type Product struct {
	//json encoders to use lowercase id
	//instead of upper case ID
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{ID: 1, Name: "Headphones", Price: 79.99},
	{ID: 2, Name: "Keyboard", Price: 49.99},
}

// global counter to products slice
var nextProductID = 3

// handler function for health and product request
// has response writer and request as parameters
func healthhandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Wishlist with me API is running \n"))
}

// for POST request ie add a product
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /products received"))
	var newProduct Product
	//r.Body is data sent with the request
	//decode reads the data
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	//if invalid json return an error and stop
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(newProduct.Name) == "" {
		http.Error(w, " Product name is required", http.StatusBadRequest)
		return
	}
	if newProduct.Price < 0 {
		http.Error(w, " Price cannot be negative", http.StatusBadRequest)
		return
	}
	//append the product in go slice
	newProduct.ID = nextProductID
	nextProductID++
	products = append(products, newProduct)
	//means server succesfully created a product
	//status 201
	w.WriteHeader(http.StatusCreated)
	//adds the new product to json
	json.NewEncoder(w).Encode(newProduct)
}

// will return products as a list
func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createProductHandler(w, r)
		return
	}
	path := strings.TrimSuffix(r.URL.Path, "/")
	//now for fetching a particular product
	if path != "/products" {
		idStr := strings.TrimPrefix(path, "/products/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Product ID", http.StatusBadRequest)
			return
		}
		for _, product := range products {
			if product.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(product)
				return
			}
		}
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	//tells cleint that im sending json
	//creates json encoder that writes directly to the response
	w.Header().Set("Content-Type", "application/json")
	//converts the slice into JSON and sends it
	json.NewEncoder(w).Encode(products)
}

func main() {
	//when request comes for /health and /product rep
	//it calls the resp function
	http.HandleFunc("/health", healthhandler)
	//2 different paths for the exact and /id paths
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productsHandler)
	//runs on 8080 port
	http.ListenAndServe(":8080", nil)
}
