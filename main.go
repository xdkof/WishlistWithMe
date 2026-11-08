package main

import (
	"encoding/json"
	"math"
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

// common validation
func validateProduct(p Product) string {
	if strings.TrimSpace(p.Name) == "" {
		return "Product name is required"
	}
	if len(p.Name) > 100 {
		return "Product name cannot exceed 100 characters"
	}
	if p.Price <= 0 {
		return "Product price must be greater than zero"
	}
	if math.IsNaN(p.Price) || math.IsInf(p.Price, 0) {
		return "Product price must be a valid finite number"
	}
	return ""
}

// global counter to products slice
var nextProductID = 3

// handler function for health and product request
// has response writer and request as parameters
func healthhandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Wishlist with me API is running \n"))
}

// It is /PUT ie add a product
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	//r.Body is data sent with the request
	//decode reads the data
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	//validate the product
	if errMsg := validateProduct(newProduct); errMsg != "" {
		http.Error(w, errMsg, http.StatusBadRequest)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

// IT is /GET for all products
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	//tells cleint that im sending json
	//creates json encoder that writes directly to the response
	w.Header().Set("Content-Type", "application/json")
	//converts the slice into JSON and sends it
	json.NewEncoder(w).Encode(products)
}

// It is  /GET for particular product by id
func getProductHandler(w http.ResponseWriter, r *http.Request, id int) {
	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

// It is /POST ie updates an already present product
func updateProductHandler(w http.ResponseWriter, r *http.Request, id int) {
	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if errMsg := validateProduct(updatedProduct); errMsg != "" {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	//u get product id from the json
	for i, product := range products {
		if product.ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

// It is /DELETE ie delete a product
func deleteProductHandler(w http.ResponseWriter, r *http.Request, id int) {
	for i, product := range products {
		if product.ID == id {
			//concat products till i and after i
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Function that handles all the Product related requests
func productsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	if path == "/products" {
		switch r.Method {
		case http.MethodGet:
			getProductsHandler(w, r)
		case http.MethodPost:
			createProductHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	idStr := strings.TrimPrefix(path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getProductHandler(w, r, id)
	case http.MethodPut:
		updateProductHandler(w, r, id)
	case http.MethodDelete:
		deleteProductHandler(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
