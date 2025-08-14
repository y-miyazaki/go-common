#!/bin/bash
set -e # Exit script if any command fails

# Display usage information
function show_usage {
  echo "Usage: $0 [options]"
  echo ""
  echo "This script retrieves SQS Dead Letter Queue (DLQ) information and outputs them in a format"
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

# Get list of SQS queues and extract DLQ information
echo_section "Retrieving SQS queue information"
if ! queue_urls=$(aws sqs list-queues | jq -r '.QueueUrls[]'); then
  echo "Error: Failed to retrieve SQS queue URLs" >&2
  exit 1
fi

echo_section "Processing Dead Letter Queue information"
# Store unique DLQ names
dlqs=()
items=""

for queue_url in $queue_urls; do
  # Get queue attributes (RedrivePolicy contains DLQ information)
  if redrive_policy=$(aws sqs get-queue-attributes --queue-url "$queue_url" --attribute-names RedrivePolicy 2> /dev/null | jq -r '.Attributes.RedrivePolicy' 2> /dev/null); then
    if [ "$redrive_policy" != "null" ] && [ -n "$redrive_policy" ]; then
      # Extract DLQ ARN from RedrivePolicy
      if dlq_arn=$(echo "$redrive_policy" | jq -r '.deadLetterTargetArn' 2> /dev/null); then
        if [ "$dlq_arn" != "null" ] && [ -n "$dlq_arn" ]; then
          # Extract queue name from ARN
          dlq_name=$(echo "$dlq_arn" | awk -F':' '{print $NF}')

          # Add DLQ name only if it doesn't already exist in the array
          dlq_exists=false
          for existing_dlq in "${dlqs[@]}"; do
            if [[ "$existing_dlq" == "$dlq_name" ]]; then
              dlq_exists=true
              break
            fi
          done

          if [[ "$dlq_exists" == "false" ]]; then
            dlqs+=("$dlq_name")

            # Append to items string
            items+="  {\n    QueueName = \"$dlq_name\"\n  },\n"
          fi
        fi
      fi
    fi
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
