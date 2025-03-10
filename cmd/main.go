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
	ec2InstanceMonthlyCostDict := map[string]float64{}
	ec2InstancesCostsSummary := float64(0)
	for _, instance := range ec2Instances {
		if _, exists := ec2InstanceMonthlyCostDict[instance.InstanceType]; !exists {
			price, err := aws.GetEC2Price(instance.InstanceType, awsRegion)
			if err != nil {
				log.Printf("❌ Failed to fetch price for EC2 %s: %v\n", instance.InstanceType, err)
				continue
			}

			monthlyCost := calculator.EstimateEC2MonthlyCost(price, 730)
			ec2InstanceMonthlyCostDict[instance.InstanceType] = monthlyCost
		}

		fmt.Printf("  ✅ Instance Type: %s, Monthly Cost: $%.2f\n", instance.InstanceType, ec2InstanceMonthlyCostDict[instance.InstanceType])
		ec2InstancesCostsSummary += ec2InstanceMonthlyCostDict[instance.InstanceType]
	}

	// Process EBS Volumes Pricing
	fmt.Println("💾 EBS Volumes:")
	ebsVolumesMonthlyCostDict := map[string]float64{}
	ebsVolumesCostsSummary := float64(0)
	for _, volume := range ebsVolumes {
		if _, exists := ebsVolumesMonthlyCostDict[volume.Type]; !exists {
			pricePerGB, err := aws.GetEBSPrice(volume.Type, awsRegion)
			if err != nil {
				log.Printf("❌ Failed to fetch price for EBS %s: %v\n", volume.Type, err)
				continue
			}

			monthlyCost := calculator.EstimateEBSMonthlyCost(pricePerGB, volume.Size)
			ebsVolumesMonthlyCostDict[volume.Type] = monthlyCost
		}
		fmt.Printf("  📦 Size: %dGB, Type: %s, Monthly Cost: $%.2f\n", volume.Size, volume.Type, ebsVolumesMonthlyCostDict[volume.Type])
		ebsVolumesCostsSummary += ebsVolumesMonthlyCostDict[volume.Type]
	}

	fmt.Println("📊 Summary:")
	fmt.Printf("📦 Total EC2 Instances Cost: $%.2f\n", ec2InstancesCostsSummary)
	fmt.Printf("💾 Total EBS Volumes Cost: $%.2f\n", ebsVolumesCostsSummary)
	fmt.Printf("💰 Total Monthly Cost: $%.2f\n", ec2InstancesCostsSummary+ebsVolumesCostsSummary)
}
