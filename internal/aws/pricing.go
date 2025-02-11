package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
)

var (
	regionMap      = make(map[string]string)
	regionMapMutex sync.Mutex
)

// FetchRegionMap retrieves AWS region names in a human-readable format from AWS Pricing API
func FetchRegionMap(ctx context.Context, pricingClient *pricing.Client) error {
	regionMapMutex.Lock()
	defer regionMapMutex.Unlock()

	// If already populated, return early
	if len(regionMap) > 0 {
		return nil
	}

	// Get EC2 pricing data to retrieve region mappings
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
	}

	resp, err := pricingClient.GetProducts(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to retrieve AWS regions from Pricing API: %w", err)
	}

	// Parse response to extract region mappings
	for _, priceJSON := range resp.PriceList {
		var priceData map[string]interface{}
		if err := json.Unmarshal([]byte(priceJSON), &priceData); err != nil {
			continue
		}

		if product, ok := priceData["product"].(map[string]interface{}); ok {
			if attributes, ok := product["attributes"].(map[string]interface{}); ok {
				regionCode, _ := attributes["regionCode"].(string)
				location, _ := attributes["location"].(string)

				if regionCode != "" && location != "" {
					regionMap[regionCode] = location
				}
			}
		}
	}

	return nil
}

// GetRegionName retrieves the correct AWS Pricing API location name
func GetRegionName(ctx context.Context, pricingClient *pricing.Client, regionCode string) (string, error) {
	err := FetchRegionMap(ctx, pricingClient)
	if err != nil {
		return "", err
	}

	regionMapMutex.Lock()
	defer regionMapMutex.Unlock()

	if region, exists := regionMap[regionCode]; exists {
		return region, nil
	}
	return "", fmt.Errorf("unknown AWS region: %s", regionCode)
}

// extractPrice extracts the price from the AWS pricing response.
func extractPrice(priceList []string) (float64, error) {
	for _, priceJSON := range priceList {
		var priceData map[string]interface{}
		if err := json.Unmarshal([]byte(priceJSON), &priceData); err != nil {
			continue
		}

		if terms, ok := priceData["terms"].(map[string]interface{}); ok {
			if onDemand, ok := terms["OnDemand"].(map[string]interface{}); ok {
				for _, term := range onDemand {
					if priceDimensions, ok := term.(map[string]interface{})["priceDimensions"].(map[string]interface{}); ok {
						for _, dimension := range priceDimensions {
							if pricePerUnit, ok := dimension.(map[string]interface{})["pricePerUnit"].(map[string]interface{}); ok {
								if usd, ok := pricePerUnit["USD"].(string); ok {
									return strconv.ParseFloat(usd, 64)
								}
							}
						}
					}
				}
			}
		}
	}
	return 0, fmt.Errorf("price not found in response")
}

// GetEC2Price fetches EC2 instance pricing dynamically
func GetEC2Price(ctx context.Context, ec2Client *ec2.Client, pricingClient *pricing.Client, instanceType, regionCode string) (float64, error) {
	// Get the correct AWS Pricing API region name
	region, err := GetRegionName(ctx, pricingClient, regionCode) // ✅ Pass the correct client
	if err != nil {
		return 0, fmt.Errorf("failed to map region: %w", err)
	}

	// fmt.Println("📍 AWS Pricing Query for:", instanceType, "in", region)

	// Define filters
	filters := []types.Filter{
		{
			Type:  types.FilterTypeTermMatch,
			Field: aws.String("instanceType"),
			Value: aws.String(instanceType),
		},
		{
			Type:  types.FilterTypeTermMatch,
			Field: aws.String("location"),
			Value: aws.String(region),
		},
	}

	// Create pricing request
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters:     filters,
	}

	resp, err := pricingClient.GetProducts(ctx, input) // ✅ Correctly using Pricing Client
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve pricing data: %w", err)
	}

	if len(resp.PriceList) == 0 {
		return 0, fmt.Errorf("price not found for instance type: %s in region: %s", instanceType, region)
	}

	// Extract price from response
	price, err := extractPrice(resp.PriceList)
	if err != nil {
		return 0, fmt.Errorf("failed to extract price: %w", err)
	}

	return price, nil
}