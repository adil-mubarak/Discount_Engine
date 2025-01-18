package discount_engine

import (
	"encoding/json"
	"os"
	"sync"
)

// codition which under the discount rule is applicable
type Condition struct {
	MinOrderValue float64 `json:"min_order_value,omitempty"` // Minimum order value
	CustomerType  string  `json:"customer_type,omitempty"`   //Cuatomer type (regular,premium)
}

// single discount rule
type DiscountRule struct {
	ID                 string    `json:"id"`                            // Unique identifier
	Description        string    `json:"description"`                   // Description of the discount rule
	Condition          Condition `json:"condition"`                     // Condition to apply the discount
	DiscountPercentage float64   `json:"discount_percentage,omitempty"` // Discount percentage
	DiscountFixed      float64   `json:"discount_fixed,omitempty"`      // Discount fixed amount
	Priority           int       `json:"priiority"`                     // Priority of the discount rule
}

// Discount engine holds the discount rules
type DiscountEngine struct {
	Rules []DiscountRule //list of discount rule
	mu    sync.Mutex     //protect shared state
}

// loads the discount rules from the configuration file
func (de *DiscountEngine) LoadRules(filePath string) error {
	//read the configuration file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err //return error ,file not found
	}

	//unmarshal the json
	err = json.Unmarshal(data, &de.Rules)
	if err != nil {
		return err //return error if unmarshaling fails
	}

	return nil
}

// evaluate the discount rules is applicable based on the order's details
func (de *DiscountEngine) EvaluateRule(rule DiscountRule, orderAmount float64, customerType string) (float64, bool) {
	//check if the order amount meets the minimum value
	if rule.Condition.MinOrderValue > 0 && orderAmount < rule.Condition.MinOrderValue {
		return 0, false // return false if the condition isn't met
	}

	//check the customer type matches the ruels condition
	if rule.Condition.CustomerType != "" && rule.Condition.CustomerType != customerType {
		return 0, false //return false if the condition isn't met
	}

	//calculate the applicable discount
	var discount float64
	if rule.DiscountPercentage > 0 {
		//calculate the percentage discount
		discount = orderAmount * rule.DiscountPercentage / 100
	} else if rule.DiscountFixed > 0 {
		//use the fixed discount if it applicable
		discount = rule.DiscountFixed
	}
	return discount, true //return the discount and the true for indicating valid
}

// calculates the best discount for the order
func (de *DiscountEngine) CalculateBestDiscount(orderAmount float64, customerType string) (float64, float64, []string) {
	var bestDiscount float64  //best deiscount to applied
	var highestPriority int   //highest priority for the rules
	var appliedRules []string //list of applied rules

	//lock the mutex to ensure the thread safety during the rule evaluation and calculation
	de.mu.Lock()
	defer de.mu.Unlock()

	//iterate over all discount rules
	for _, rule := range de.Rules {
		discount, valid := de.EvaluateRule(rule, orderAmount, customerType)
		if valid {
			//if this rule have higher priority or a best discount, update the discount
			if rule.Priority > highestPriority || (rule.Priority == highestPriority && discount > bestDiscount) {
				highestPriority = rule.Priority           //update the highest priority
				bestDiscount = discount                   // update the discount
				appliedRules = []string{rule.Description} //reser the applied rules
			} else if rule.Priority == highestPriority && discount == bestDiscount {
				//if the rules have the same priority and discount, append the rule
				appliedRules = append(appliedRules, rule.Description)
			}
		}
	}

	//calculate the final discount after applying the best discount
	finalOrderTotal := orderAmount - bestDiscount
	//return the discount,final amount, applied rules
	return bestDiscount, finalOrderTotal, appliedRules
}
