package main

import (
	"discount_engine/discount_engine"
	"encoding/json"
	"fmt"
	"net/http"
)

// used to bind the incoming json
type OrderRequest struct {
	OrderAmount  float64 `json:"order_amount"`
	CustomerType string  `json:"customer_type"`
}

// used to format the response with the total,discount,applied ruels
type DiscountResponse struct {
	DiscountApplied float64  `json:"discount_applied"`
	FinalAmount     float64  `json:"final_amount"`
	AppliedRules    []string `json:"applied_rules"`
}

// handle incoming HTTP post request , applies the discount and return the result
func HandleDiscountRequest(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	//extracts the order details from the incoming json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return //if failed to extratct the order details, return an error
	}

	//initializing discount engine
	discountEngine := &discount_engine.DiscountEngine{}
	//passing path discount_rules.json
	if err := discountEngine.LoadRules("discount_rules.json"); err != nil {
		http.Error(w, fmt.Sprintf("Error loading discount rules: %s", err), http.StatusInternalServerError)
		return // return an error if failed to load the discount rules
	}

	//calculate the best discount
	discount, finalAmount, appliedRules := discountEngine.CalculateBestDiscount(req.OrderAmount, req.CustomerType)

	//prepare the response
	response := DiscountResponse{
		DiscountApplied: discount,
		FinalAmount:     finalAmount,
		AppliedRules:    appliedRules,
	}
	//return the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error loading response: %s", err), http.StatusInternalServerError)
		return
	}

}

// start the HTTP server and listen for the request
func main() {
	http.HandleFunc("/discount", HandleDiscountRequest)
	http.ListenAndServe(":8080", nil)
}
