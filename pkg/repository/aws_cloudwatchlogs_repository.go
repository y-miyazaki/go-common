package repository

// Package-level repository for AWS CloudWatch Logs using AWS SDK for Go v2.
// This mirrors the repository style used for other AWS services (e.g., S3).

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

var (
	// ErrLogStreamNotFound indicates that a requested log stream could not be found.
	ErrLogStreamNotFound = errors.New("log stream not found")
)

// AWSCloudWatchLogsRepositoryInterface defines the contract for CloudWatch Logs operations.
// nolint:iface,revive,unused
type AWSCloudWatchLogsRepositoryInterface interface {
	CreateLogGroup(group string) (*cloudwatchlogs.CreateLogGroupOutput, error)
	CreateLogStream(group, stream string) (*cloudwatchlogs.CreateLogStreamOutput, error)
	PutRetentionPolicy(group string, logRetentionInDays int32) (*cloudwatchlogs.PutRetentionPolicyOutput, error)
	DescribeLogGroups(prefix string) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	DescribeLogStreams(group, prefix string) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
	PutLogEvents(group, stream string, events []types.InputLogEvent, sequenceToken *string) (*cloudwatchlogs.PutLogEventsOutput, error)
	FilterLogEvents(group string, startTime, endTime int64, filterPattern string, nextToken *string, limit int32) (*cloudwatchlogs.FilterLogEventsOutput, error)
	GetNextSequenceToken(group, stream string) (*string, error)
}

// AWSCloudWatchLogsRepository implements AWSCloudWatchLogsRepositoryInterface backed by AWS SDK v2 client.
type AWSCloudWatchLogsRepository struct {
	c *cloudwatchlogs.Client
}

// NewAWSCloudWatchLogsRepository returns a new repository instance using the provided CloudWatch Logs client.
func NewAWSCloudWatchLogsRepository(c *cloudwatchlogs.Client) *AWSCloudWatchLogsRepository {
	return &AWSCloudWatchLogsRepository{c: c}
}

// CreateLogGroup creates a new CloudWatch Logs log group.
func (r *AWSCloudWatchLogsRepository) CreateLogGroup(group string) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	out, err := r.c.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(group),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs CreateLogGroup: %w", err)
	}
	return out, nil
}

// CreateLogStream creates a new log stream within the specified log group.
func (r *AWSCloudWatchLogsRepository) CreateLogStream(group, stream string) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	out, err := r.c.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs CreateLogStream: %w", err)
	}
	return out, nil
}

// PutRetentionPolicy sets the retention policy for a log group.
func (r *AWSCloudWatchLogsRepository) PutRetentionPolicy(group string, logRetentionInDays int32) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
	out, err := r.c.PutRetentionPolicy(context.TODO(), &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(group),
		RetentionInDays: aws.Int32(logRetentionInDays),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs PutRetentionPolicy: %w", err)
	}
	return out, nil
}

// DescribeLogGroups returns log groups optionally filtered by a name prefix.
func (r *AWSCloudWatchLogsRepository) DescribeLogGroups(prefix string) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	in := &cloudwatchlogs.DescribeLogGroupsInput{}
	if prefix != "" {
		in.LogGroupNamePrefix = aws.String(prefix)
	}
	out, err := r.c.DescribeLogGroups(context.TODO(), in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs DescribeLogGroups: %w", err)
	}
	return out, nil
}

// DescribeLogStreams returns log streams in a group optionally filtered by a name prefix.
func (r *AWSCloudWatchLogsRepository) DescribeLogStreams(group, prefix string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	in := &cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: aws.String(group)}
	if prefix != "" {
		in.LogStreamNamePrefix = aws.String(prefix)
	}
	out, err := r.c.DescribeLogStreams(context.TODO(), in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs DescribeLogStreams: %w", err)
	}
	return out, nil
}

// PutLogEvents uploads log events to the specified log stream.
// The sequenceToken should be provided for subsequent calls after the first PutLogEvents.
func (r *AWSCloudWatchLogsRepository) PutLogEvents(group, stream string, events []types.InputLogEvent, sequenceToken *string) (*cloudwatchlogs.PutLogEventsOutput, error) {
	in := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		LogEvents:     events,
	}
	if sequenceToken != nil {
		in.SequenceToken = sequenceToken
	}
	out, err := r.c.PutLogEvents(context.TODO(), in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs PutLogEvents: %w", err)
	}
	return out, nil
}

// FilterLogEvents searches log events in a log group with optional time range and filter pattern.
func (r *AWSCloudWatchLogsRepository) FilterLogEvents(group string, startTime, endTime int64, filterPattern string, nextToken *string, limit int32) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	in := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(group),
	}
	if startTime > 0 {
		in.StartTime = aws.Int64(startTime)
	}
	if endTime > 0 {
		in.EndTime = aws.Int64(endTime)
	}
	if filterPattern != "" {
		in.FilterPattern = aws.String(filterPattern)
	}
	if nextToken != nil {
		in.NextToken = nextToken
	}
	if limit > 0 {
		in.Limit = aws.Int32(limit)
	}
	out, err := r.c.FilterLogEvents(context.TODO(), in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs FilterLogEvents: %w", err)
	}
	return out, nil
}

// GetNextSequenceToken retrieves the next sequence token for a log stream.
func (r *AWSCloudWatchLogsRepository) GetNextSequenceToken(group, stream string) (*string, error) {
	streams, err := r.DescribeLogStreams(group, stream)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs DescribeLogStreams: %w", err)
	}
	if len(streams.LogStreams) == 0 {
		return nil, ErrLogStreamNotFound
	}
	return streams.LogStreams[0].UploadSequenceToken, nil
}
