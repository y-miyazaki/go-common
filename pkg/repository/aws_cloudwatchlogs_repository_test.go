package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCloudWatchLogsClient is a mock implementation of the CloudWatch Logs client
type MockCloudWatchLogsClient struct {
	mock.Mock
}

func (m *MockCloudWatchLogsClient) CreateLogGroup(ctx context.Context, params *cloudwatchlogs.CreateLogGroupInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.CreateLogGroupOutput), args.Error(1)
}

func (m *MockCloudWatchLogsClient) CreateLogStream(ctx context.Context, params *cloudwatchlogs.CreateLogStreamInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.CreateLogStreamOutput), args.Error(1)
}

func (m *MockCloudWatchLogsClient) PutRetentionPolicy(ctx context.Context, params *cloudwatchlogs.PutRetentionPolicyInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.PutRetentionPolicyOutput), args.Error(1)
}

func (m *MockCloudWatchLogsClient) DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput), args.Error(1)
}

func (m *MockCloudWatchLogsClient) DescribeLogStreams(ctx context.Context, params *cloudwatchlogs.DescribeLogStreamsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.DescribeLogStreamsOutput), args.Error(1)
}

func (m *MockCloudWatchLogsClient) PutLogEvents(ctx context.Context, params *cloudwatchlogs.PutLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.PutLogEventsOutput), args.Error(1)
}

func (m *MockCloudWatchLogsClient) FilterLogEvents(ctx context.Context, params *cloudwatchlogs.FilterLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudwatchlogs.FilterLogEventsOutput), args.Error(1)
}

func TestNewAWSCloudWatchLogsRepository(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.c)
}

func TestAWSCloudWatchLogsRepository_CreateLogGroup(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	mockClient.On("CreateLogGroup", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.CreateLogGroupInput) bool {
		return *input.LogGroupName == logGroupName
	}), mock.Anything).Return(&cloudwatchlogs.CreateLogGroupOutput{}, nil)

	_, err := repo.CreateLogGroup(context.Background(), logGroupName)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_CreateLogGroup_Error(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	expectedError := errors.New("cloudwatch logs error")

	mockClient.On("CreateLogGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	_, err := repo.CreateLogGroup(context.Background(), logGroupName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cloudwatchlogs CreateLogGroup")
	assert.Contains(t, err.Error(), expectedError.Error())
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_CreateLogStream(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	logStreamName := "test-log-stream"

	mockClient.On("CreateLogStream", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.CreateLogStreamInput) bool {
		return *input.LogGroupName == logGroupName && *input.LogStreamName == logStreamName
	}), mock.Anything).Return(&cloudwatchlogs.CreateLogStreamOutput{}, nil)

	_, err := repo.CreateLogStream(context.Background(), logGroupName, logStreamName)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_PutLogEvents(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	logStreamName := "test-log-stream"
	logEvents := []types.InputLogEvent{
		{
			Message:   aws.String("Test log message"),
			Timestamp: aws.Int64(time.Now().UnixMilli()),
		},
	}

	expectedOutput := &cloudwatchlogs.PutLogEventsOutput{
		NextSequenceToken: aws.String("next-token"),
	}

	mockClient.On("PutLogEvents", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.PutLogEventsInput) bool {
		return *input.LogGroupName == logGroupName && *input.LogStreamName == logStreamName && len(input.LogEvents) == 1
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.PutLogEvents(context.Background(), logGroupName, logStreamName, logEvents, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_DescribeLogGroups(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupNamePrefix := "test-"
	expectedOutput := &cloudwatchlogs.DescribeLogGroupsOutput{
		LogGroups: []types.LogGroup{
			{
				LogGroupName:    aws.String("test-log-group-1"),
				CreationTime:    aws.Int64(time.Now().Unix()),
				RetentionInDays: aws.Int32(30),
			},
		},
	}

	mockClient.On("DescribeLogGroups", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.DescribeLogGroupsInput) bool {
		return input.LogGroupNamePrefix != nil && *input.LogGroupNamePrefix == logGroupNamePrefix
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.DescribeLogGroups(context.Background(), logGroupNamePrefix)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	assert.Len(t, result.LogGroups, 1)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_PutRetentionPolicy(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	retentionInDays := int32(30)

	mockClient.On("PutRetentionPolicy", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.PutRetentionPolicyInput) bool {
		return *input.LogGroupName == logGroupName && *input.RetentionInDays == retentionInDays
	}), mock.Anything).Return(&cloudwatchlogs.PutRetentionPolicyOutput{}, nil)

	_, err := repo.PutRetentionPolicy(context.Background(), logGroupName, retentionInDays)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_DescribeLogStreams(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	logStreamNamePrefix := "test-stream-"

	expectedOutput := &cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{
			{
				LogStreamName:       aws.String("test-stream-1"),
				CreationTime:        aws.Int64(time.Now().Unix()),
				UploadSequenceToken: aws.String("token-123"),
			},
		},
	}

	mockClient.On("DescribeLogStreams", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.DescribeLogStreamsInput) bool {
		return *input.LogGroupName == logGroupName && input.LogStreamNamePrefix != nil && *input.LogStreamNamePrefix == logStreamNamePrefix
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.DescribeLogStreams(context.Background(), logGroupName, logStreamNamePrefix)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	assert.Len(t, result.LogStreams, 1)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_FilterLogEvents(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	startTime := int64(1640995200) // 2022-01-01 00:00:00 UTC
	endTime := int64(1641081600)   // 2022-01-02 00:00:00 UTC
	filterPattern := "ERROR"

	expectedOutput := &cloudwatchlogs.FilterLogEventsOutput{
		Events: []types.FilteredLogEvent{
			{
				LogStreamName: aws.String("test-stream"),
				Timestamp:     aws.Int64(startTime + 3600),
				Message:       aws.String("ERROR: Something went wrong"),
			},
		},
	}

	mockClient.On("FilterLogEvents", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.FilterLogEventsInput) bool {
		return *input.LogGroupName == logGroupName && *input.StartTime == startTime && *input.EndTime == endTime && *input.FilterPattern == filterPattern
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.FilterLogEvents(context.Background(), logGroupName, startTime, endTime, filterPattern, nil, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	assert.Len(t, result.Events, 1)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_GetNextSequenceToken(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	logStreamName := "test-stream"
	expectedToken := "sequence-token-123"

	streamsOutput := &cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{
			{
				LogStreamName:       aws.String(logStreamName),
				UploadSequenceToken: aws.String(expectedToken),
			},
		},
	}

	mockClient.On("DescribeLogStreams", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.DescribeLogStreamsInput) bool {
		return *input.LogGroupName == logGroupName && *input.LogStreamNamePrefix == logStreamName
	}), mock.Anything).Return(streamsOutput, nil)

	result, err := repo.GetNextSequenceToken(context.Background(), logGroupName, logStreamName)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, *result)
	mockClient.AssertExpectations(t)
}

func TestAWSCloudWatchLogsRepository_GetNextSequenceToken_StreamNotFound(t *testing.T) {
	mockClient := &MockCloudWatchLogsClient{}
	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)

	logGroupName := "test-log-group"
	logStreamName := "non-existent-stream"

	streamsOutput := &cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{}, // Empty result
	}

	mockClient.On("DescribeLogStreams", mock.Anything, mock.Anything, mock.Anything).Return(streamsOutput, nil)

	result, err := repo.GetNextSequenceToken(context.Background(), logGroupName, logStreamName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "log stream not found")
	mockClient.AssertExpectations(t)
}

// Test actual AWSCloudWatchLogsRepository functions with integration-style tests
func TestAWSCloudWatchLogsRepositoryReal_PutRetentionPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if no real AWS credentials are available
	t.Skip("Skipping CloudWatch Logs integration test - requires real AWS credentials")

	// This would test the actual repository with real AWS client
	// mockClient := &cloudwatchlogs.Client{} // Real client
	// repo := NewAWSCloudWatchLogsRepository(mockClient)
	// ... test actual AWS calls
}

func TestAWSCloudWatchLogsRepositoryReal_DescribeLogStreams(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping CloudWatch Logs integration test - requires real AWS credentials")
}

func TestAWSCloudWatchLogsRepositoryReal_FilterLogEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping CloudWatch Logs integration test - requires real AWS credentials")
}

func TestAWSCloudWatchLogsRepositoryReal_GetNextSequenceToken(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping CloudWatch Logs integration test - requires real AWS credentials")
}

func TestAWSCloudWatchLogsRepositoryReal_GetNextSequenceToken_StreamNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping CloudWatch Logs integration test - requires real AWS credentials")
}
