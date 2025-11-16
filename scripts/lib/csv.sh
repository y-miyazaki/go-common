#!/bin/bash
#######################################
# Description: CSV utility functions for shell scripts
# Usage: source /path/to/scripts/lib/csv.sh
#
# This library provides CSV data processing and normalization functions:
# - CSV value normalization with empty/null handling
# - CSV safety processing for special characters
# - Support for preserving newlines in CSV output
# - Quote escaping and comma handling for proper CSV format
#######################################

#######################################
# CSV Value Normalization Functions
#######################################

#######################################
# normalize_csv_value: Normalize and quote a value for safe CSV output
# Arguments:
#   $1 - value to normalize
# Outputs:
#   CSV-safe value, always quoted if contains comma, quote, or newline
#   Handles PRESERVE_NEWLINES for Excel/Numbers compatibility
#   Escapes double quotes per RFC 4180
#######################################
function normalize_csv_value {
    local value="$1"

    # Treat empty/null as empty string
    if [[ -z "$value" || "$value" == "null" ]]; then
        echo ""
        return
    fi

    if [[ "${PRESERVE_NEWLINES:-false}" == "true" ]]; then
        # Preserve newlines, only escape double quotes
        value="${value//\"/\"\"}"
        echo "\"$value\""
    else
        # Replace newlines with literal \n for compatibility
        value="${value//$'\n'/\\n}"
        value="${value//\"/\"\"}"
        echo "\"$value\""
    fi
}

#######################################
# Helper function to make values CSV-safe by removing/replacing problematic characters
# Uses PRESERVE_NEWLINES environment variable to control newline handling
# Arguments:
#   $1 - value to make CSV-safe
# Outputs:
#   CSV-safe value with proper quoting and escaping
#######################################
function make_csv_safe {
    local value=$1

    if [[ "${PRESERVE_NEWLINES:-false}" == "true" ]]; then
        # Preserve newlines for better Excel/Numbers compatibility
        # Only remove carriage returns, keep newlines
        value=$(echo "$value" | tr -d '\r')

        # Convert tabs to spaces
        value=$(echo "$value" | tr '\t' ' ')

        # Compress multiple spaces to single space (but preserve newlines)
        # Use a more precise pattern that doesn't affect newlines
        value=$(echo "$value" | sed 's/[[:blank:]]\{2,\}/ /g' | sed 's/^[[:blank:]]*//;s/[[:blank:]]*$//')

        # If value contains commas, wrap in quotes and escape internal quotes
        if [[ "$value" == *","* ]]; then
            # Escape existing quotes by doubling them
            value=${value//\"/\"\"}
            value="\"$value\""
        elif [[ "$value" == *"\""* ]] || [[ "$value" == *$'\n'* ]] || [[ "$value" =~ ^[[:blank:]]*$ ]]; then
            # Escape existing quotes by doubling them
            value=${value//\"/\"\"}
            value="\"$value\""
        fi
    else
        # Remove carriage returns and newlines, replace with spaces (default for compatibility)
        value=$(echo "$value" | tr -d '\r' | tr '\n' ' ')

        # Convert tabs to spaces
        value=$(echo "$value" | tr '\t' ' ')

        # Compress multiple spaces to single space
        value=$(echo "$value" | tr -s ' ')

        # Trim leading/trailing whitespace
        value=$(echo "$value" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

        # If value contains commas, wrap in quotes and escape internal quotes
        if [[ "$value" == *","* ]]; then
            # Escape existing quotes by doubling them
            value=${value//\"/\"\"}
            value="\"$value\""
        elif [[ "$value" == *"\""* ]] || [[ "$value" =~ ^[[:space:]]*$ ]]; then
            # Escape existing quotes by doubling them
            value=${value//\"/\"\"}
            value="\"$value\""
        fi
    fi

    echo "$value"
}

#######################################
# csv_sort: Sort CSV data properly handling quoted fields with commas and newlines
# Arguments:
#   $1 - CSV data to sort
# Outputs:
#   Sorted CSV data with proper field handling
# Sorting order: Region(3), Subcategory(1), Subsubcategory(2), Name(4)
#######################################
function csv_sort {
    local input_data="$1"

    # For complex CSV data with quoted fields, use Python to sort properly
    if command -v python3 > /dev/null 2>&1; then
        printf "%b" "$input_data" | python3 -c "
import csv
import sys
from io import StringIO

# Read CSV data from stdin
csv_data = sys.stdin.read()
reader = csv.reader(StringIO(csv_data))
rows = list(reader)

# Sort by multiple columns: Region(3), Subcategory(1), Subsubcategory(2), Name(4)
# Columns are 0-indexed, so subtract 1 from 1-indexed sort keys
try:
    sorted_rows = sorted(rows, key=lambda row: (
        row[3] if len(row) > 3 else '',  # Region
        row[1] if len(row) > 1 else '',  # Subcategory
        row[2] if len(row) > 2 else '',  # Subsubcategory
        row[4] if len(row) > 4 else ''   # Name
    ))

    # Write sorted CSV to stdout
    writer = csv.writer(sys.stdout, quoting=csv.QUOTE_MINIMAL)
    for row in sorted_rows:
        writer.writerow(row)
except Exception as e:
    # Fallback: output unsorted if Python sorting fails
    sys.stdout.write(csv_data)
"
    else
        # Fallback to basic sort if Python is not available
        printf "%b" "$input_data" | sort -t, -k4,4 -k2,2 -k3,3 -k5,5
    fi
}
