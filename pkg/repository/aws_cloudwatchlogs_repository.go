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

// AWSCloudWatchLogsClientInterface defines the interface for CloudWatch Logs client operations
type AWSCloudWatchLogsClientInterface interface {
	// CreateLogGroup creates a new CloudWatch Logs log group
	CreateLogGroup(_ context.Context, _ *cloudwatchlogs.CreateLogGroupInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogGroupOutput, error)
	// CreateLogStream creates a new log stream within the specified log group
	CreateLogStream(_ context.Context, _ *cloudwatchlogs.CreateLogStreamInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error)
	// PutRetentionPolicy sets the retention policy for a log group
	PutRetentionPolicy(_ context.Context, _ *cloudwatchlogs.PutRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error)
	// DescribeLogGroups returns log groups optionally filtered by a name prefix
	DescribeLogGroups(_ context.Context, _ *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	// DescribeLogStreams returns log streams in a group optionally filtered by a name prefix
	DescribeLogStreams(_ context.Context, _ *cloudwatchlogs.DescribeLogStreamsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
	// PutLogEvents uploads log events to the specified log stream
	PutLogEvents(_ context.Context, _ *cloudwatchlogs.PutLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error)
	// FilterLogEvents searches log events in a log group with optional time range and filter pattern
	FilterLogEvents(_ context.Context, _ *cloudwatchlogs.FilterLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.FilterLogEventsOutput, error)
}

// AWSCloudWatchLogsRepository implements AWSCloudWatchLogsRepositoryInterface backed by AWS SDK v2 client.
type AWSCloudWatchLogsRepository struct {
	Client AWSCloudWatchLogsClientInterface
}

// NewAWSCloudWatchLogsRepository returns a new repository instance using the provided CloudWatch Logs client.
func NewAWSCloudWatchLogsRepository(c *cloudwatchlogs.Client) *AWSCloudWatchLogsRepository {
	return &AWSCloudWatchLogsRepository{Client: c}
}

// NewAWSCloudWatchLogsRepositoryWithInterface returns a new repository instance using the provided client interface (for testing).
func NewAWSCloudWatchLogsRepositoryWithInterface(c AWSCloudWatchLogsClientInterface) *AWSCloudWatchLogsRepository {
	return &AWSCloudWatchLogsRepository{Client: c}
}

// CreateLogGroup creates a new CloudWatch Logs log group.
func (r *AWSCloudWatchLogsRepository) CreateLogGroup(ctx context.Context, group string) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	out, err := r.Client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(group),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs CreateLogGroup: %w", err)
	}
	return out, nil
}

// CreateLogStream creates a new log stream within the specified log group.
func (r *AWSCloudWatchLogsRepository) CreateLogStream(ctx context.Context, group, stream string) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	out, err := r.Client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs CreateLogStream: %w", err)
	}
	return out, nil
}

// PutRetentionPolicy sets the retention policy for a log group.
func (r *AWSCloudWatchLogsRepository) PutRetentionPolicy(ctx context.Context, group string, logRetentionInDays int32) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
	out, err := r.Client.PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(group),
		RetentionInDays: aws.Int32(logRetentionInDays),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs PutRetentionPolicy: %w", err)
	}
	return out, nil
}

// DescribeLogGroups returns log groups optionally filtered by a name prefix.
func (r *AWSCloudWatchLogsRepository) DescribeLogGroups(ctx context.Context, prefix string) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	in := &cloudwatchlogs.DescribeLogGroupsInput{}
	if prefix != "" {
		in.LogGroupNamePrefix = aws.String(prefix)
	}
	out, err := r.Client.DescribeLogGroups(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs DescribeLogGroups: %w", err)
	}
	return out, nil
}

// DescribeLogStreams returns log streams in a group optionally filtered by a name prefix.
func (r *AWSCloudWatchLogsRepository) DescribeLogStreams(ctx context.Context, group, prefix string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	in := &cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: aws.String(group)}
	if prefix != "" {
		in.LogStreamNamePrefix = aws.String(prefix)
	}
	out, err := r.Client.DescribeLogStreams(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs DescribeLogStreams: %w", err)
	}
	return out, nil
}

// PutLogEvents uploads log events to the specified log stream.
// The sequenceToken should be provided for subsequent calls after the first PutLogEvents.
func (r *AWSCloudWatchLogsRepository) PutLogEvents(ctx context.Context, group, stream string, events []types.InputLogEvent, sequenceToken *string) (*cloudwatchlogs.PutLogEventsOutput, error) {
	in := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		LogEvents:     events,
	}
	if sequenceToken != nil {
		in.SequenceToken = sequenceToken
	}
	out, err := r.Client.PutLogEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs PutLogEvents: %w", err)
	}
	return out, nil
}

// FilterLogEvents searches log events in a log group with optional time range and filter pattern.
func (r *AWSCloudWatchLogsRepository) FilterLogEvents(ctx context.Context, group string, startTime, endTime int64, filterPattern string, nextToken *string, limit int32) (*cloudwatchlogs.FilterLogEventsOutput, error) {
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
	out, err := r.Client.FilterLogEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs FilterLogEvents: %w", err)
	}
	return out, nil
}

// GetNextSequenceToken retrieves the next sequence token for a log stream.
func (r *AWSCloudWatchLogsRepository) GetNextSequenceToken(ctx context.Context, group, stream string) (*string, error) {
	streams, err := r.DescribeLogStreams(ctx, group, stream)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs DescribeLogStreams: %w", err)
	}
	if len(streams.LogStreams) == 0 {
		return nil, ErrLogStreamNotFound
	}
	return streams.LogStreams[0].UploadSequenceToken, nil
}
