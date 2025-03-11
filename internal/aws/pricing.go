package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
)

// GetEC2Price retrieves the hourly price of an EC2 instance type in a specific region
func GetEC2Price(instanceType, region string) (float64, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return 0, err
	}

	client := pricing.NewFromConfig(cfg)

	// Convert region code (e.g., us-east-2) to full AWS region name
	regionMapping, err := GetRegionMapping()
	if err != nil {
		return 0, err
	}

	location, exists := regionMapping[region]
	if !exists {
		return 0, fmt.Errorf("no pricing data available for region: %s", region)
	}

	// AWS Pricing API query for EC2 instances
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []types.Filter{
			{Type: types.FilterTypeTermMatch, Field: aws.String("instanceType"), Value: aws.String(instanceType)},
			{Type: types.FilterTypeTermMatch, Field: aws.String("location"), Value: aws.String(location)},
			{Type: types.FilterTypeTermMatch, Field: aws.String("tenancy"), Value: aws.String("Shared")},        // TODO
			{Type: types.FilterTypeTermMatch, Field: aws.String("operatingSystem"), Value: aws.String("Linux")}, // TODO
		},
	}

	resp, err := client.GetProducts(context.TODO(), input)
	if err != nil {
		return 0, err
	}

	// fmt.Println("Response:", resp)
	if len(resp.PriceList) == 0 {
		log.Println("❌ No pricing data found for EC2:", instanceType, "in", location)
		return 0, fmt.Errorf("no pricing data found for EC2 instance: %s", instanceType)
	}

	price := extractPrice(resp.PriceList)
	return price, nil
}

// GetEBSPrice retrieves the price per GB/month for an EBS volume type in a specific region
func GetEBSPrice(volumeType, region string) (float64, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return 0, err
	}

	client := pricing.NewFromConfig(cfg)

	// Convert region code (e.g., us-east-2) to full AWS region name
	regionMapping, err := GetRegionMapping()
	if err != nil {
		return 0, err
	}

	location, exists := regionMapping[region]
	if !exists {
		return 0, fmt.Errorf("no pricing data available for region: %s", region)
	}

	// AWS Pricing API query for EBS volumes
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []types.Filter{
			{Type: types.FilterTypeTermMatch, Field: aws.String("productFamily"), Value: aws.String("Storage")},
			{Type: types.FilterTypeTermMatch, Field: aws.String("volumeApiName"), Value: aws.String(volumeType)},
			{Type: types.FilterTypeTermMatch, Field: aws.String("location"), Value: aws.String(location)},
		},
	}

	resp, err := client.GetProducts(context.TODO(), input)
	if err != nil {
		return 0, err
	}

	if len(resp.PriceList) == 0 {
		log.Println("❌ No pricing data found for EBS volume type:", volumeType, "in", location)
		return 0, fmt.Errorf("no pricing data found for EBS volume type: %s", volumeType)
	}

	price := extractPrice(resp.PriceList)
	return price, nil
}

// extractPrice parses the AWS Pricing API response and retrieves the hourly price
func extractPrice(priceList []string) float64 {
	for _, item := range priceList {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(item), &result); err != nil {
			log.Println("❌ Error parsing price JSON:", err)
			continue
		}

		if terms, ok := result["terms"].(map[string]interface{}); ok {
			if onDemand, ok := terms["OnDemand"].(map[string]interface{}); ok {
				for _, offer := range onDemand {
					priceData := offer.(map[string]interface{})["priceDimensions"].(map[string]interface{})
					for _, pd := range priceData {
						priceStr := pd.(map[string]interface{})["pricePerUnit"].(map[string]interface{})["USD"].(string)
						var price float64
						fmt.Sscanf(priceStr, "%f", &price)
						return price
					}
				}
			}
		}
	}
	return 0.0
}

// GetRegionMapping fetches AWS Pricing API region names dynamically
func GetRegionMapping() (map[string]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	client := pricing.NewFromConfig(cfg)

	input := &pricing.GetAttributeValuesInput{
		ServiceCode:   aws.String("AmazonEC2"),
		AttributeName: aws.String("location"),
	}

	resp, err := client.GetAttributeValues(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	regionMap := make(map[string]string)
	for _, attr := range resp.AttributeValues {
		locationName := *attr.Value
		regionCode := inferRegionFromLocation(locationName)

		if regionCode != "" {
			regionMap[regionCode] = locationName
		}
	}

	// log.Println("✅ Successfully retrieved AWS region mapping:", regionMap)
	return regionMap, nil
}

// inferRegionFromLocation extracts region codes dynamically from AWS Pricing API
func inferRegionFromLocation(location string) string {
	regionMapping := map[string]string{
		"US East (N. Virginia)":    "us-east-1",
		"US East (Ohio)":           "us-east-2",
		"US West (N. California)":  "us-west-1",
		"US West (Oregon)":         "us-west-2",
		"EU (Frankfurt)":           "eu-central-1",
		"EU (Ireland)":             "eu-west-1",
		"EU (London)":              "eu-west-2",
		"EU (Paris)":               "eu-west-3",
		"Asia Pacific (Singapore)": "ap-southeast-1",
		"Asia Pacific (Sydney)":    "ap-southeast-2",
		"Asia Pacific (Tokyo)":     "ap-northeast-1",
		"Asia Pacific (Seoul)":     "ap-northeast-2",
	}

	if regionCode, exists := regionMapping[location]; exists {
		return regionCode
	}
	return ""
}
