package unit_test

import (
	"discount_engine/discount_engine"
	"testing"
)

// verify that no discount applied the order does not satisfy any rules
func TestNoDiscount(t *testing.T) {
	engine := &discount_engine.DiscountEngine{}
	engine.Rules = []discount_engine.DiscountRule{
		{
			ID:          "rules_1",
			Description: "10% off for orders over $100",
			Condition: discount_engine.Condition{
				MinOrderValue: 100,
			},
			Priority: 1,
		},
	}

	discount, finalAmount, _ := engine.CalculateBestDiscount(50, "regular")
	if discount != 0 {
		t.Errorf("Expected 0 discount, got %f", discount)
	}
	if finalAmount != 50 {
		t.Errorf("Expected final amount to be 50, got %f", finalAmount)
	}

}

// test the engines ability to handle conflicting rules and apply the rule with higheset priority
func TestConflictRules(t *testing.T) {
	engine := &discount_engine.DiscountEngine{}
	engine.Rules = []discount_engine.DiscountRule{
		{
			ID:          "rules_1",
			Description: "10% off for orders over $100",
			Condition: discount_engine.Condition{
				MinOrderValue: 100,
			},
			DiscountPercentage: 10,
			Priority:           1,
		},
		{
			ID:          "rules_2",
			Description: "15% off for orders over $100",
			Condition: discount_engine.Condition{
				MinOrderValue: 100,
			},
			DiscountPercentage: 15,
			Priority:           2,
		},
	}

	discount, finalAmount, appliedRules := engine.CalculateBestDiscount(150, "regular")
	if discount != 22.5 {
		t.Errorf("Expected 22.5 discount, got %v", discount)
	}
	if finalAmount != 127.5 {
		t.Errorf("Expected final amount to be 127.5, got %v", finalAmount)
	}
	if len(appliedRules) != 1 || appliedRules[0] != "15% off for orders over $100" {
		t.Errorf("Expected 1 rule to be applied, got %v", appliedRules)
	}
}

// test that when two rules have the same priority, the first rule is applied
func TestEqualPriority(t *testing.T) {
	engine := &discount_engine.DiscountEngine{}
	engine.Rules = []discount_engine.DiscountRule{
		{
			ID:          "rules_1",
			Description: "10% off for orders over $100",
			Condition: discount_engine.Condition{
				MinOrderValue: 100,
			},
			DiscountPercentage: 10,
			Priority:           1,
		},
		{
			ID:          "rules_2",
			Description: "10% off for orders over $100",
			Condition: discount_engine.Condition{
				MinOrderValue: 100,
			},
			DiscountPercentage: 5,
			Priority:           1,
		},
	}

	discount, finalAmount, appliedRules := engine.CalculateBestDiscount(150, "regular")

	if discount != 15 {
		t.Errorf("Expected 15 discount, got %v", discount)
	}
	if finalAmount != 135 {
		t.Errorf("Expected final amount to be 135, got %v", finalAmount)
	}
	if len(appliedRules) != 1 || appliedRules[0] != "10% off for orders over $100" {
		t.Errorf("Expected 1 rule to be applied, got %v", appliedRules)
	}
}
