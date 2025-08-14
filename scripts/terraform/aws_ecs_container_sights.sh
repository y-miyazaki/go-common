#!/bin/bash
set -e # Exit script if any command fails

# Display usage information
function show_usage {
  echo "Usage: $0 [options]"
  echo ""
  echo "This script retrieves ECS task definition families and outputs them in a format"
  echo "suitable for Terraform configuration."
  echo ""
  echo "Options:"
  echo "  -h, --help    Show this help message"
  echo ""
  echo "Example: $0"
  exit 1
}

# Show section header
function echo_section {
  echo "#--------------------------------------------------------------"
  echo "# $1"
  echo "#--------------------------------------------------------------"
}

# Show usage if -h or --help is provided
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  show_usage
fi

# Check if required tools are installed
if ! command -v aws &> /dev/null; then
  echo "Error: AWS CLI is not installed or not in PATH" >&2
  exit 1
fi

if ! command -v jq &> /dev/null; then
  echo "Error: jq is not installed or not in PATH" >&2
  exit 1
fi

# Get list of task definitions
echo_section "Retrieving ECS task definitions"
if ! task_definitions=$(aws ecs list-task-definitions --output json | jq -r '.taskDefinitionArns[]'); then
  echo "Error: Failed to retrieve ECS task definitions" >&2
  exit 1
fi

# Store unique task definition families
echo_section "Processing task definition families"
unique_families=()

# Build output items
items=""
for task_definition in $task_definitions; do
  # Extract task definition family
  family=$(echo "$task_definition" | awk -F'/' '{print $NF}' | awk -F':' '{print $1}')

  # Add family only if it doesn't already exist in the array
  family_exists=false
  for existing_family in "${unique_families[@]}"; do
    if [[ "$existing_family" == "$family" ]]; then
      family_exists=true
      break
    fi
  done

  if [[ "$family_exists" == "false" ]]; then
    unique_families+=("$family")

    # Append to items string
    items+="  {\n    ClusterName = \"\"\n    TaskDefinitionFamily = \"$family\"\n  },\n"
  fi
done

# Remove trailing comma if items is not empty
if [ -n "$items" ]; then
  items=${items%,*}
fi

# Output formatted result
echo_section "Generating output"
echo "["
echo -e "$items"
echo "]"

echo_section "Completed"
