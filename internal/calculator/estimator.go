package calculator

// EstimateEC2MonthlyCost calculates the monthly cost for an EC2 instance
func EstimateEC2MonthlyCost(hourlyPrice float64, hoursPerMonth int) float64 {
	return hourlyPrice * float64(hoursPerMonth)
}

// EstimateEBSMonthlyCost calculates the cost based on storage size (per GB/month)
func EstimateEBSMonthlyCost(pricePerGB float64, sizeGB int) float64 {
	return pricePerGB * float64(sizeGB)
}