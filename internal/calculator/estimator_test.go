package calculator

import (
	"testing"
)

func TestEstimateEC2MonthlyCost(t *testing.T) {
	// Test cases for the Estimator function
	testCases := []struct {
		hourlyCost   float64
		monthlyHours int
		expectedCost float64
	}{
		{hourlyCost: 0.17, monthlyHours: 100, expectedCost: 17.0},
		{hourlyCost: 0.25, monthlyHours: 50, expectedCost: 12.5},
		{hourlyCost: 0.35, monthlyHours: 100, expectedCost: 35.0},
	}
	for _, tc := range testCases {
		actualCost := EstimateEC2MonthlyCost(tc.hourlyCost, tc.monthlyHours)
		if actualCost != tc.expectedCost {
			t.Errorf("Estimator(%f, %d) = %f, want %f", tc.hourlyCost, tc.monthlyHours, actualCost, tc.expectedCost)
		}
	}
}

func TestEstimateEBSMonthlyCost(t *testing.T) {
	// Test cases for the Estimator function
	testCases := []struct {
		hourlyCost   float64
		monthlyHours int
		expectedCost float64
	}{
		{hourlyCost: 0.17, monthlyHours: 100, expectedCost: 17.0},
		{hourlyCost: 0.25, monthlyHours: 50, expectedCost: 12.5},
		{hourlyCost: 0.35, monthlyHours: 100, expectedCost: 35.0},
	}
	for _, tc := range testCases {
		actualCost := EstimateEBSMonthlyCost(tc.hourlyCost, tc.monthlyHours)
		if actualCost != tc.expectedCost {
			t.Errorf("Estimator(%f, %d) = %f, want %f", tc.hourlyCost, tc.monthlyHours, actualCost, tc.expectedCost)
		}
	}
}
