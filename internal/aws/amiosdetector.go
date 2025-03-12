package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func DetectAMIOperatingSystem(region, amiID string) (string, string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", "", fmt.Errorf("failed to load AWS config: %s, %w", "", err)
	}
	// Create EC2 client using the provided config
	client := ec2.NewFromConfig(cfg)

	// Describe the AMI to get its details
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{amiID},
	}

	result, err := client.DescribeImages(context.TODO(), input)
	if err != nil {
		return "", "", fmt.Errorf("failed to describe AMI: %w", err)
	}

	if len(result.Images) == 0 {
		return "", "", fmt.Errorf("AMI not found: %s", amiID)
	}

	// Extract OS and software information from the image
	image := result.Images[0]
	description := strings.ToLower(aws.ToString(image.Description))
	name := strings.ToLower(aws.ToString(image.Name))

	// Default preInstalledSw is "NA" (Not Applicable)
	preInstalledSw := "NA"
	operatingSystem := ""

	// Detect pre-installed software
	if strings.Contains(description, "sql server") || strings.Contains(name, "sql server") {
		if strings.Contains(description, "enterprise") || strings.Contains(name, "enterprise") {
			preInstalledSw = "SQL Server Enterprise"
		} else if strings.Contains(description, "standard") || strings.Contains(name, "standard") {
			preInstalledSw = "SQL Server Standard"
		} else if strings.Contains(description, "web") || strings.Contains(name, "web") {
			preInstalledSw = "SQL Server Web"
		} else {
			preInstalledSw = "SQL Server"
		}
	}

	// If platform is specified as "windows", it's Windows
	if image.Platform == "windows" {
		// For Windows with SQL Server, the operating system should be "Windows"
		// and preInstalledSw should be the SQL Server edition
		operatingSystem = "Windows"
		return operatingSystem, preInstalledSw, nil
	}

	// For Linux, determine the distribution
	switch {
	case strings.Contains(description, "rhel") || strings.Contains(name, "rhel"):
		operatingSystem = "Red Hat Enterprise Linux"
	case strings.Contains(description, "suse") || strings.Contains(name, "suse"):
		operatingSystem = "SUSE Linux"
	default:
		// Default to Linux for all other Linux distributions
		operatingSystem = "Linux"
	}

	return operatingSystem, preInstalledSw, nil
}
