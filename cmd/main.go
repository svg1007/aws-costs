package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"

	"github.com/svg1007/aws-costs/internal/aws"
	"github.com/svg1007/aws-costs/internal/calculator"
	"github.com/svg1007/aws-costs/internal/terraform"
)

func main() {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")) // Pricing API only available in us-east-1
	if err != nil {
		log.Fatalf("❌ Unable to load AWS SDK config: %v", err)
	}

	// Create AWS clients
	pricingClient := pricing.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)

	// Parse Terraform plan JSON
	plan, err := terraform.ParseTerraformPlan("terraform/tfplan.json")
	if err != nil {
		log.Fatalf("❌ Error parsing Terraform plan: %v", err)
	}

	// Extract AWS region from Terraform plan
	regionCode := plan.Configuration.ProviderConfig.AWS.Expressions.Region.ConstantValue
	fmt.Println("🌍 AWS Region Code from Terraform:", regionCode)

	// Retrieve the correct AWS Pricing API region name
	pricingRegion, err := aws.GetRegionName(context.TODO(), pricingClient, regionCode)
	if err != nil {
		log.Fatalf("❌ Failed to map AWS region: %v", err)
	}
	fmt.Println("📍 Mapped AWS Pricing Region:", pricingRegion)

	// Track total cost and instance count
	totalCost := 0.0
	instanceCount := 0

	// Process EC2 instances and estimate costs
	for _, instance := range plan.GetEC2Instances() {
		price, err := aws.GetEC2Price(context.TODO(), ec2Client, pricingClient, instance.InstanceType, regionCode)
		if err != nil {
			log.Printf("❌ Failed to fetch price for %s: %v\n", instance.InstanceType, err)
			continue
		}

		monthlyCost := calculator.EstimateMonthlyCost(price, 730)
		// fmt.Printf("✅ Instance: %s, Region: %s, Monthly Cost: $%.2f\n", instance.InstanceType, regionCode, monthlyCost)

		// Update total cost and count
		totalCost += monthlyCost
		instanceCount++
	}

	// Print final summary
	fmt.Println("\n📊 **Summary**")
	fmt.Printf("🔢 Total Instances Processed: %d\n", instanceCount)
	fmt.Printf("💰 Total Estimated Monthly Cost: $%.2f\n", totalCost)
}