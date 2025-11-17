#!/usr/bin/env bats

# Tests for scripts/lib/aws.sh (pure functions)

setup() {
    # Source library (tests run from repo root)
    source "scripts/lib/aws.sh"
}

@test "extract_jq_value returns default for empty json" {
    run extract_jq_value "" '.Name' "DEFAULT"
    [ "$status" -eq 0 ]
    [ "$output" = "DEFAULT" ]
}

@test "extract_jq_value extracts value from json" {
    local js
    js='{"Name":"Alice","Age":30}'
    run extract_jq_value "$js" '.Name' "DEFAULT"
    [ "$status" -eq 0 ]
    [ "$output" = "Alice" ]
}

@test "extract_jq_array returns joined quoted list for array" {
    local js
    js='{"Tags":["a","b","c"]}'
    run extract_jq_array "$js" '.Tags'
    [ "$status" -eq 0 ]
    [ "$output" = '"a,b,c"' ]
}

@test "format_aws_timestamp converts seconds" {
    run format_aws_timestamp 1609459200
    [ "$status" -eq 0 ]
    [[ "$output" == 2021* ]]
}

@test "format_aws_timestamp handles milliseconds" {
    run format_aws_timestamp 1609459200000
    [ "$status" -eq 0 ]
    [[ "$output" == 2021* ]]
}

@test "parse_arn returns json with components" {
    run parse_arn "arn:aws:ec2:us-east-1:123456789012:instance/i-0123456789abcdef0"
    [ "$status" -eq 0 ]
    [[ "$output" =~ '"service": "ec2"' ]]
    [[ "$output" =~ '"region": "us-east-1"' ]]
}

@test "get_resource_name_from_arn extracts resource name" {
    run get_resource_name_from_arn "arn:aws:ec2:us-east-1:123456789012:instance/i-0123456789abcdef0"
    [ "$status" -eq 0 ]
    [ "$output" = "i-0123456789abcdef0" ]
}

@test "get_waf_name extracts name from ARN path" {
    local arn="arn:aws:wafv2:us-east-1:123456789012:regional/webacl/MyWebACL/uuid"
    run get_waf_name "$arn"
    [ "$status" -eq 0 ]
    [ "$output" = "MyWebACL" ]
}
