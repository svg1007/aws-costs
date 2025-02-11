package calculator

// EstimateMonthlyCost calculates the monthly cost based on hourly rate
func EstimateMonthlyCost(hourlyPrice float64, hoursPerMonth int) float64 {
	return hourlyPrice * float64(hoursPerMonth)
}