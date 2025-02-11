package terraform

import (
	"encoding/json"
	"os"
)

// TerraformPlan represents the Terraform JSON output structure
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
	ResourceChanges []struct {
		Type   string `json:"type"`
		Change struct {
			Actions []string `json:"actions"`
			After   struct {
				InstanceType string `json:"instance_type"`
			} `json:"after"`
		} `json:"change"`
	} `json:"resource_changes"`
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

// GetEC2Instances extracts EC2 instances from Terraform plan
func (t *TerraformPlan) GetEC2Instances() []struct {
	InstanceType string
} {
	var instances []struct {
		InstanceType string
	}

	for _, resource := range t.ResourceChanges {
		if resource.Type == "aws_instance" && contains(resource.Change.Actions, "create") {
			instances = append(instances, struct {
				InstanceType string
			}{
				InstanceType: resource.Change.After.InstanceType,
			})
		}
	}

	return instances
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