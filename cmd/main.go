package main

import (
	"fmt"
	"log"

	"github.com/svg1007/aws-costs/internal/aws"
	"github.com/svg1007/aws-costs/internal/calculator"
	"github.com/svg1007/aws-costs/internal/terraform"
)

func main() {
	// Parse Terraform plan JSON
	plan, err := terraform.ParseTerraformPlan("terraform/tfplan.json")
	if err != nil {
		log.Fatalf("❌ Error parsing Terraform plan: %v", err)
	}

	// Extract EC2 Instances, EBS Volumes, and AWS Region
	ec2Instances, ebsVolumes, awsRegion := plan.ExtractResources()

	// Print AWS Region
	fmt.Println("🌍 AWS Region:", awsRegion)

	// Process EC2 Instances Pricing
	fmt.Println("🚀 EC2 Instances:")
	for _, instance := range ec2Instances {
		price, err := aws.GetEC2Price(instance.InstanceType, awsRegion)
		if err != nil {
			log.Printf("❌ Failed to fetch price for EC2 %s: %v\n", instance.InstanceType, err)
			continue
		}

		monthlyCost := calculator.EstimateEC2MonthlyCost(price, 730)
		fmt.Printf("  ✅ Instance Type: %s, Monthly Cost: $%.2f\n", instance.InstanceType, monthlyCost)
	}

	// Process EBS Volumes Pricing
	fmt.Println("💾 EBS Volumes:")
	for _, volume := range ebsVolumes {
		pricePerGB, err := aws.GetEBSPrice(volume.Type, awsRegion)
		if err != nil {
			log.Printf("❌ Failed to fetch price for EBS %s: %v\n", volume.Type, err)
			continue
		}

		monthlyCost := calculator.EstimateEBSMonthlyCost(pricePerGB, volume.Size)
		fmt.Printf("  📦 Size: %dGB, Type: %s, Monthly Cost: $%.2f\n", volume.Size, volume.Type, monthlyCost)
	}
}