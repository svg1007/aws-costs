package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/svg1007/aws-costs/internal/aws"
	"github.com/svg1007/aws-costs/internal/calculator"
	"github.com/svg1007/aws-costs/internal/terraform"
)

func main() {
	verbose := flag.Bool("verbose", false, "Show detailed information for each resource")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Missing required argument <plan-file> in JSON format")
		fmt.Println("Usage: aws-costs <plan-file> [-v|--verbose]")
		os.Exit(1)
	}

	// First positional argument is the required plan-file
	planFile := args[0]

	// Parse Terraform plan JSON
	plan, err := terraform.ParseTerraformPlan(planFile)
	if err != nil {
		log.Fatalf("❌ Error parsing Terraform plan: %v", err)
	}

	// Extract EC2 Instances, EBS Volumes, and AWS Region
	ec2Instances, ebsVolumes, awsRegion := plan.ExtractResources()

	// Print AWS Region
	if *verbose {
		fmt.Println("🌍 AWS Region:", awsRegion)
	}

	// Process EC2 Instances Pricing
	if *verbose {
		fmt.Println("🚀 EC2 Instances:")
	}
	ec2InstanceMonthlyCostDict := map[string]float64{}
	ec2InstancesCostsSummary := float64(0)
	for _, instance := range ec2Instances {
		if _, exists := ec2InstanceMonthlyCostDict[instance.InstanceType]; !exists {
			os, preinstalledSw, err := aws.DetectAMIOperatingSystem(awsRegion, instance.Ami)
			if err != nil {
				log.Printf("❌ Failed to fetch OS for EC2 %s: %v\n", instance.Ami, err)
				continue
			}

			price, err := aws.GetEC2Price(awsRegion, instance.InstanceType, os, preinstalledSw)
			if err != nil {
				log.Printf("❌ Failed to fetch price for EC2 %s: %v\n", instance.InstanceType, err)
				continue
			}

			monthlyCost := calculator.EstimateEC2MonthlyCost(price, 730)
			ec2InstanceMonthlyCostDict[instance.InstanceType] = monthlyCost
		}

		if *verbose {
			fmt.Printf("  ✅ Instance Type: %s, Monthly Cost: $%.2f\n", instance.InstanceType, ec2InstanceMonthlyCostDict[instance.InstanceType])
		}
		ec2InstancesCostsSummary += ec2InstanceMonthlyCostDict[instance.InstanceType]
	}

	// Process EBS Volumes Pricing
	if *verbose {
		fmt.Println("💾 EBS Volumes:")
	}
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

		if *verbose {
			fmt.Printf("  📦 Size: %dGB, Type: %s, Monthly Cost: $%.2f\n", volume.Size, volume.Type, ebsVolumesMonthlyCostDict[volume.Type])
		}
		ebsVolumesCostsSummary += ebsVolumesMonthlyCostDict[volume.Type]
	}

	fmt.Println("📊 Summary:")
	fmt.Printf("  📦 Total EC2 Instances Cost: $%.2f\n", ec2InstancesCostsSummary)
	fmt.Printf("  💾 Total EBS Volumes Cost: $%.2f\n", ebsVolumesCostsSummary)
	fmt.Printf("  💰 Total Monthly Cost: $%.2f\n", ec2InstancesCostsSummary+ebsVolumesCostsSummary)
}
