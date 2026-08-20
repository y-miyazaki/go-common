package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestCloudWatchLogs = errors.New("cloudwatch logs error")

func TestNewAWSCloudWatchLogsRepository(t *testing.T) {
	t.Parallel()

	mockClient := &cloudwatchlogs.Client{}
	repo := NewAWSCloudWatchLogsRepository(mockClient)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
}

func TestAWSCloudWatchLogsRepository_CreateLogGroup(t *testing.T) {
	t.Parallel()

	logGroupName := "test-log-group"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSCloudWatchLogsClientInterface)
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSCloudWatchLogsClientInterface) {
				m.EXPECT().
					CreateLogGroup(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.CreateLogGroupInput{})).
					DoAndReturn(func(_ context.Context, input *cloudwatchlogs.CreateLogGroupInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogGroupOutput, error) {
						require.NotNil(t, input.LogGroupName)
						require.Equal(t, logGroupName, *input.LogGroupName)
						return &cloudwatchlogs.CreateLogGroupOutput{}, nil
					})
			},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSCloudWatchLogsClientInterface) {
				m.EXPECT().
					CreateLogGroup(gomock.Any(), gomock.Any()).
					Return(nil, errTestCloudWatchLogs)
			},
			wantErrMsg: "cloudwatchlogs CreateLogGroup",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
			_, err := repo.CreateLogGroup(context.Background(), logGroupName)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Contains(t, err.Error(), errTestCloudWatchLogs.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestAWSCloudWatchLogsRepository_CreateLogStream(t *testing.T) {
	t.Parallel()

	logGroupName := "test-log-group"
	logStreamName := "test-log-stream"

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
	mockClient.EXPECT().
		CreateLogStream(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.CreateLogStreamInput{})).
		DoAndReturn(func(_ context.Context, input *cloudwatchlogs.CreateLogStreamInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error) {
			require.NotNil(t, input.LogGroupName)
			require.NotNil(t, input.LogStreamName)
			require.Equal(t, logGroupName, *input.LogGroupName)
			require.Equal(t, logStreamName, *input.LogStreamName)
			return &cloudwatchlogs.CreateLogStreamOutput{}, nil
		})

	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	_, err := repo.CreateLogStream(context.Background(), logGroupName, logStreamName)

	require.NoError(t, err)
}

func TestAWSCloudWatchLogsRepository_PutLogEvents(t *testing.T) {
	t.Parallel()

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

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
	mockClient.EXPECT().
		PutLogEvents(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.PutLogEventsInput{})).
		DoAndReturn(func(_ context.Context, input *cloudwatchlogs.PutLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error) {
			require.NotNil(t, input.LogGroupName)
			require.NotNil(t, input.LogStreamName)
			require.Equal(t, logGroupName, *input.LogGroupName)
			require.Equal(t, logStreamName, *input.LogStreamName)
			require.Len(t, input.LogEvents, 1)
			return expectedOutput, nil
		})

	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	got, err := repo.PutLogEvents(context.Background(), logGroupName, logStreamName, logEvents, nil)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, got)
}

func TestAWSCloudWatchLogsRepository_DescribeLogGroups(t *testing.T) {
	t.Parallel()

	in := &cloudwatchlogs.DescribeLogGroupsInput{LogGroupNamePrefix: aws.String("test-")}
	expectedOutput := &cloudwatchlogs.DescribeLogGroupsOutput{
		LogGroups: []types.LogGroup{
			{
				LogGroupName:    aws.String("test-log-group-1"),
				CreationTime:    aws.Int64(time.Now().Unix()),
				RetentionInDays: aws.Int32(30),
			},
		},
	}

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
	mockClient.EXPECT().
		DescribeLogGroups(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.DescribeLogGroupsInput{})).
		DoAndReturn(func(_ context.Context, input *cloudwatchlogs.DescribeLogGroupsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
			require.NotNil(t, input.LogGroupNamePrefix)
			require.Equal(t, *in.LogGroupNamePrefix, *input.LogGroupNamePrefix)
			return expectedOutput, nil
		})

	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	got, err := repo.DescribeLogGroups(context.Background(), in)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, got)
	require.Len(t, got.LogGroups, 1)
}

func TestAWSCloudWatchLogsRepository_PutRetentionPolicy(t *testing.T) {
	t.Parallel()

	logGroupName := "test-log-group"
	retentionInDays := int32(30)

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
	mockClient.EXPECT().
		PutRetentionPolicy(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.PutRetentionPolicyInput{})).
		DoAndReturn(func(_ context.Context, input *cloudwatchlogs.PutRetentionPolicyInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
			require.NotNil(t, input.LogGroupName)
			require.NotNil(t, input.RetentionInDays)
			require.Equal(t, logGroupName, *input.LogGroupName)
			require.Equal(t, retentionInDays, *input.RetentionInDays)
			return &cloudwatchlogs.PutRetentionPolicyOutput{}, nil
		})

	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	_, err := repo.PutRetentionPolicy(context.Background(), logGroupName, retentionInDays)

	require.NoError(t, err)
}

func TestAWSCloudWatchLogsRepository_DescribeLogStreams(t *testing.T) {
	t.Parallel()

	in := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String("test-log-group"),
		LogStreamNamePrefix: aws.String("test-stream-"),
	}
	expectedOutput := &cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{
			{
				LogStreamName:       aws.String("test-stream-1"),
				CreationTime:        aws.Int64(time.Now().Unix()),
				UploadSequenceToken: aws.String("token-123"),
			},
		},
	}

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
	mockClient.EXPECT().
		DescribeLogStreams(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.DescribeLogStreamsInput{})).
		DoAndReturn(func(_ context.Context, input *cloudwatchlogs.DescribeLogStreamsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
			require.NotNil(t, input.LogGroupName)
			require.NotNil(t, input.LogStreamNamePrefix)
			require.Equal(t, *in.LogGroupName, *input.LogGroupName)
			require.Equal(t, *in.LogStreamNamePrefix, *input.LogStreamNamePrefix)
			return expectedOutput, nil
		})

	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	got, err := repo.DescribeLogStreams(context.Background(), in)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, got)
	require.Len(t, got.LogStreams, 1)
}

func TestAWSCloudWatchLogsRepository_FilterLogEvents(t *testing.T) {
	t.Parallel()

	startTime := int64(1640995200) // 2022-01-01 00:00:00 UTC
	endTime := int64(1641081600)   // 2022-01-02 00:00:00 UTC
	in := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String("test-log-group"),
		StartTime:     aws.Int64(startTime),
		EndTime:       aws.Int64(endTime),
		FilterPattern: aws.String("ERROR"),
	}
	expectedOutput := &cloudwatchlogs.FilterLogEventsOutput{
		Events: []types.FilteredLogEvent{
			{
				LogStreamName: aws.String("test-stream"),
				Timestamp:     aws.Int64(startTime + 3600),
				Message:       aws.String("ERROR: Something went wrong"),
			},
		},
	}

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
	mockClient.EXPECT().
		FilterLogEvents(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.FilterLogEventsInput{})).
		DoAndReturn(func(_ context.Context, input *cloudwatchlogs.FilterLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.FilterLogEventsOutput, error) {
			require.NotNil(t, input.LogGroupName)
			require.NotNil(t, input.StartTime)
			require.NotNil(t, input.EndTime)
			require.NotNil(t, input.FilterPattern)
			require.Equal(t, *in.LogGroupName, *input.LogGroupName)
			require.Equal(t, *in.StartTime, *input.StartTime)
			require.Equal(t, *in.EndTime, *input.EndTime)
			require.Equal(t, *in.FilterPattern, *input.FilterPattern)
			return expectedOutput, nil
		})

	repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
	got, err := repo.FilterLogEvents(context.Background(), in)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, got)
	require.Len(t, got.Events, 1)
}

func TestAWSCloudWatchLogsRepository_GetNextSequenceToken(t *testing.T) {
	t.Parallel()

	logGroupName := "test-log-group"
	logStreamName := "test-stream"
	expectedToken := "sequence-token-123"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSCloudWatchLogsClientInterface)
		want       *string
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSCloudWatchLogsClientInterface) {
				streamsOutput := &cloudwatchlogs.DescribeLogStreamsOutput{
					LogStreams: []types.LogStream{
						{
							LogStreamName:       aws.String(logStreamName),
							UploadSequenceToken: aws.String(expectedToken),
						},
					},
				}
				m.EXPECT().
					DescribeLogStreams(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.DescribeLogStreamsInput{})).
					DoAndReturn(func(_ context.Context, input *cloudwatchlogs.DescribeLogStreamsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
						require.NotNil(t, input.LogGroupName)
						require.NotNil(t, input.LogStreamNamePrefix)
						require.Equal(t, logGroupName, *input.LogGroupName)
						require.Equal(t, logStreamName, *input.LogStreamNamePrefix)
						return streamsOutput, nil
					})
			},
			want: aws.String(expectedToken),
		},
		{
			name: "stream not found",
			setupMock: func(m *mocks.MockAWSCloudWatchLogsClientInterface) {
				m.EXPECT().
					DescribeLogStreams(gomock.Any(), gomock.Any()).
					Return(&cloudwatchlogs.DescribeLogStreamsOutput{LogStreams: []types.LogStream{}}, nil)
			},
			wantErrMsg: "log stream not found",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
			got, err := repo.GetNextSequenceToken(context.Background(), logGroupName, logStreamName)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestAWSCloudWatchLogsRepository_DescribeMetricFilters(t *testing.T) {
	t.Parallel()

	in := &cloudwatchlogs.DescribeMetricFiltersInput{LogGroupName: aws.String("test-log-group")}
	expectedOutput := &cloudwatchlogs.DescribeMetricFiltersOutput{
		MetricFilters: []types.MetricFilter{
			{
				FilterName:    aws.String("test-filter"),
				FilterPattern: aws.String("[...]"),
				LogGroupName:  aws.String(*in.LogGroupName),
				CreationTime:  aws.Int64(time.Now().UnixMilli()),
				MetricTransformations: []types.MetricTransformation{
					{
						MetricName:  aws.String("test-metric-name"),
						MetricValue: aws.String("1"),
					},
				},
			},
		},
	}

	tests := []struct {
		name       string
		input      *cloudwatchlogs.DescribeMetricFiltersInput
		setupMock  func(m *mocks.MockAWSCloudWatchLogsClientInterface)
		want       *cloudwatchlogs.DescribeMetricFiltersOutput
		wantErrMsg string
	}{
		{
			name:  "success",
			input: in,
			setupMock: func(m *mocks.MockAWSCloudWatchLogsClientInterface) {
				m.EXPECT().
					DescribeMetricFilters(gomock.Any(), gomock.AssignableToTypeOf(&cloudwatchlogs.DescribeMetricFiltersInput{})).
					DoAndReturn(func(_ context.Context, input *cloudwatchlogs.DescribeMetricFiltersInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeMetricFiltersOutput, error) {
						require.NotNil(t, input.LogGroupName)
						require.Equal(t, *in.LogGroupName, *input.LogGroupName)
						return expectedOutput, nil
					})
			},
			want: expectedOutput,
		},
		{
			name:  "client error",
			input: &cloudwatchlogs.DescribeMetricFiltersInput{LogGroupName: aws.String("test-log-group")},
			setupMock: func(m *mocks.MockAWSCloudWatchLogsClientInterface) {
				m.EXPECT().
					DescribeMetricFilters(gomock.Any(), gomock.Any()).
					Return(nil, errTestCloudWatchLogs)
			},
			wantErrMsg: "cloudwatchlogs DescribeMetricFilters",
		},
		{
			name:  "nil input",
			input: nil,
			wantErrMsg: "DescribeMetricFiltersInput cannot be nil",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSCloudWatchLogsClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSCloudWatchLogsRepositoryWithInterface(mockClient)
			got, err := repo.DescribeMetricFilters(context.Background(), tc.input)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Nil(t, got)
				if tc.input == nil {
					return
				}
				require.Contains(t, err.Error(), errTestCloudWatchLogs.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
			require.Len(t, got.MetricFilters, 1)
		})
	}
}

func TestAWSCloudWatchLogsRepositoryReal_PutRetentionPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping CloudWatch Logs integration test - requires real AWS credentials")
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
