# AWS Costs CLI

AWS Costs CLI is a command-line tool for estimating and analyzing AWS EC2 costs using the AWS API. This tool helps you track your AWS spending efficiently and provides insights into cost optimization.

## Features

- Fetches EC2 instance costs from AWS Pricing API

- Provides a detailed breakdown of instance types and their costs

- Supports cost filtering by instance type and region

- Allows exporting cost data for further analysis

## Installation

Ensure you have Go installed on your system, then clone the repository and build the binary:

### Clone the repository
```
git clone https://github.com/svg1007/aws-costs.git
cd aws-costs
```

### Build the binary
`go build -o aws-costs`

## Usage

### Basic Usage

To estimate AWS EC2 costs, run:

`./aws-costs <plan-file> [options]`

Where:

- `<plan-file>` is the mandatory Terraform plan JSON file.

- `[options]` are optional flags for customization.

### Options

| Flag | Alias | Description |
|------|-------|-------------|
| `-v` | `--verbose` | Show detailed information for each resource |
| `-h` | `--help` | Display help information |

### Example

`./aws-costs my-plan.json`

#### Example output:
```
📊 Summary:
  📦 Total EC2 Instances Cost: $630.72
  💾 Total EBS Volumes Cost: $76.00
  💰 Total Monthly Cost: $706.72
```

## Requirements

- Go 1.18+

- AWS CLI configured with appropriate IAM permissions

## License

This project is licensed under the MIT License.
