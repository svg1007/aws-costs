package terraform

import (
	"encoding/json"
	"fmt"
	"os"
)

// TerraformPlan represents the structure of Terraform's JSON output
type TerraformPlan struct {
	Configuration struct {
		ProviderConfig struct {
			AWS struct {
				Expressions struct {
					Region struct {
						ConstantValue string `json:"constant_value"`
					} `json:"region"`
				} `json:"expressions"`
			} `json:"aws"`
		} `json:"provider_config"`
	} `json:"configuration"`

	ResourceChanges []ResourceChange `json:"resource_changes"`
}

// ResourceChange represents a single resource change in Terraform
type ResourceChange struct {
	Type   string `json:"type"` // e.g., "aws_instance", "aws_ebs_volume"
	Change struct {
		Actions []string        `json:"actions"`
		After   json.RawMessage `json:"after"` // Raw JSON for conditional parsing
	} `json:"change"`
}

// EC2Instance represents EC2 instances in Terraform
type EC2Instance struct {
	InstanceType string `json:"instance_type"`
}

// EBSVolume represents EBS volumes in Terraform
type EBSVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

// ParseTerraformPlan reads and parses Terraform JSON output
func ParseTerraformPlan(filePath string) (*TerraformPlan, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var plan TerraformPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}

// ExtractResources extracts EC2 instances and EBS volumes dynamically
func (t *TerraformPlan) ExtractResources() ([]EC2Instance, []EBSVolume, string) {
	var ec2Instances []EC2Instance
	var ebsVolumes []EBSVolume

	// Extract AWS Region
	awsRegion := t.Configuration.ProviderConfig.AWS.Expressions.Region.ConstantValue

	for _, resource := range t.ResourceChanges {
		if contains(resource.Change.Actions, "create") {
			switch resource.Type {
			case "aws_instance":
				var instance EC2Instance
				if err := json.Unmarshal(resource.Change.After, &instance); err == nil {
					ec2Instances = append(ec2Instances, instance)
				} else {
					fmt.Println("❌ Error parsing EC2 instance:", err)
				}
			case "aws_ebs_volume":
				var volume EBSVolume
				if err := json.Unmarshal(resource.Change.After, &volume); err == nil {
          // if volume.Type == "" {
          //   volume.Type = "gp2"
          // }
					ebsVolumes = append(ebsVolumes, volume)
				} else {
					fmt.Println("❌ Error parsing EBS volume:", err)
				}
			}
		}
	}

	return ec2Instances, ebsVolumes, awsRegion
}

// Helper function to check if an action exists in the change set
func contains(actions []string, action string) bool {
	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}