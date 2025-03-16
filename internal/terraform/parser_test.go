package terraform

import (
	"testing"
)

func TestParseTerraformPlan(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		errorAllowed bool
	}{
		{"Valid plan", "testdata/plan_valid.json", false},
		{"Invalid JSON", "testdata/plan_invalid.json", true},
		{"Missing plan", "testdata/plan_missing.json", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := ParseTerraformPlan(tc.path)

			if tc.errorAllowed {
				if err == nil {
					t.Errorf("Expected an error for %s, but got none", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error parsing Terraform plan: %v", err)
				}
				if plan.TerraformVersion != "1.5.0" {
					t.Errorf("Expected TerraformVersion '1.5.0', got: %s", plan.TerraformVersion)
				}
			}
		})
	}
}

func TestExtractResources(t *testing.T) {
	plan, err := ParseTerraformPlan("testdata/plan_valid.json")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	ec2Instances, ebsVolumes, awsRegion := plan.ExtractResources()

	if len(ec2Instances) != 1 {
		t.Errorf("Expected 1 EC2 instance, got: %d", len(ec2Instances))
	}

	expectedInstance := EC2Instance{
		InstanceType: "t3.micro",
		Ami:          "ami-12345678",
		RootBlockDevice: []struct {
			Size int    `json:"volume_size"`
			Type string `json:"volume_type"`
		}{
			{Size: 8, Type: "gp2"},
		},
	}
	if ec2Instances[0].InstanceType != expectedInstance.InstanceType {
		t.Errorf("Expected InstanceType '%s', got '%s'", expectedInstance.InstanceType, ec2Instances[0].InstanceType)
	}
	if ec2Instances[0].Ami != expectedInstance.Ami {
		t.Errorf("Expected AMI '%s', got '%s'", expectedInstance.Ami, ec2Instances[0].Ami)
	}
	if ec2Instances[0].RootBlockDevice[0].Size != expectedInstance.RootBlockDevice[0].Size {
		t.Errorf("Expected RootBlockDevice size '%d', got '%d'", expectedInstance.RootBlockDevice[0].Size, ec2Instances[0].RootBlockDevice[0].Size)
	}

	if len(ebsVolumes) != 2 {
		t.Errorf("Expected 1 EBS volume, got: %d", len(ebsVolumes))
	}
	if ebsVolumes[1].Size != 10 {
		t.Errorf("Expected EBS volume size '%d', got '%d'", 10, ebsVolumes[1].Size)
	}
	if ebsVolumes[1].Type != "gp2" {
		t.Errorf("Expected EBS volume type '%s', got '%s'", "gp2", ebsVolumes[1].Type)
	}

	if awsRegion != "us-east-1" {
		t.Errorf("Expected AWS Region 'us-east-1', got: %s", awsRegion)
	}
}
