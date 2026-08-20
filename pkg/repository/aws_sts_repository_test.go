package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestSTS = errors.New("sts error")

func TestNewAWSSTSRepository(t *testing.T) {
	t.Parallel()

	mockClient := &sts.Client{}
	repo := NewAWSSTSRepository(mockClient)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
}

func TestAWSSTSRepository_GetCallerIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSSTSClientInterface)
		want       *sts.GetCallerIdentityOutput
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				expected := &sts.GetCallerIdentityOutput{Account: aws.String("123456789012")}
				m.EXPECT().
					GetCallerIdentity(gomock.Any(), gomock.AssignableToTypeOf(&sts.GetCallerIdentityInput{})).
					Return(expected, nil)
			},
			want: &sts.GetCallerIdentityOutput{Account: aws.String("123456789012")},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				m.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Return(nil, errTestSTS)
			},
			wantErrMsg: "sts GetCallerIdentity",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSTSClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSTSRepositoryWithInterface(mockClient)
			got, err := repo.GetCallerIdentity(context.Background())

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

func TestAWSSTSRepository_GetAccessKeyInfo(t *testing.T) {
	t.Parallel()

	accessKeyID := "AKIAIOSFODNN7EXAMPLE"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSSTSClientInterface)
		want       *sts.GetAccessKeyInfoOutput
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				expected := &sts.GetAccessKeyInfoOutput{Account: aws.String("123456789012")}
				m.EXPECT().
					GetAccessKeyInfo(gomock.Any(), gomock.AssignableToTypeOf(&sts.GetAccessKeyInfoInput{})).
					DoAndReturn(func(_ context.Context, input *sts.GetAccessKeyInfoInput, _ ...func(*sts.Options)) (*sts.GetAccessKeyInfoOutput, error) {
						require.NotNil(t, input.AccessKeyId)
						require.Equal(t, accessKeyID, *input.AccessKeyId)
						return expected, nil
					})
			},
			want: &sts.GetAccessKeyInfoOutput{Account: aws.String("123456789012")},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				m.EXPECT().GetAccessKeyInfo(gomock.Any(), gomock.Any()).Return(nil, errTestSTS)
			},
			wantErrMsg: "sts GetAccessKeyInfo",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSTSClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSTSRepositoryWithInterface(mockClient)
			got, err := repo.GetAccessKeyInfo(context.Background(), accessKeyID)

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

func TestAWSSTSRepository_GetDelegatedAccessToken(t *testing.T) {
	t.Parallel()

	tradeInToken := "trade-token"

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockAWSSTSClientInterface)
		want       *sts.GetDelegatedAccessTokenOutput
		wantErrMsg string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				expected := &sts.GetDelegatedAccessTokenOutput{}
				m.EXPECT().
					GetDelegatedAccessToken(gomock.Any(), gomock.AssignableToTypeOf(&sts.GetDelegatedAccessTokenInput{})).
					DoAndReturn(func(_ context.Context, input *sts.GetDelegatedAccessTokenInput, _ ...func(*sts.Options)) (*sts.GetDelegatedAccessTokenOutput, error) {
						require.NotNil(t, input.TradeInToken)
						require.Equal(t, tradeInToken, *input.TradeInToken)
						return expected, nil
					})
			},
			want: &sts.GetDelegatedAccessTokenOutput{},
		},
		{
			name: "client error",
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				m.EXPECT().GetDelegatedAccessToken(gomock.Any(), gomock.Any()).Return(nil, errTestSTS)
			},
			wantErrMsg: "sts GetDelegatedAccessToken",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSTSClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSTSRepositoryWithInterface(mockClient)
			got, err := repo.GetDelegatedAccessToken(context.Background(), tradeInToken)

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

func TestAWSSTSRepository_GetFederationToken(t *testing.T) {
	t.Parallel()

	in := &sts.GetFederationTokenInput{Name: aws.String("federated-user")}

	tests := []struct {
		name       string
		input      *sts.GetFederationTokenInput
		setupMock  func(m *mocks.MockAWSSTSClientInterface)
		want       *sts.GetFederationTokenOutput
		wantErrMsg string
	}{
		{
			name:  "success",
			input: in,
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				expected := &sts.GetFederationTokenOutput{}
				m.EXPECT().GetFederationToken(gomock.Any(), in).Return(expected, nil)
			},
			want: &sts.GetFederationTokenOutput{},
		},
		{
			name:  "nil input uses empty request",
			input: nil,
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				m.EXPECT().
					GetFederationToken(gomock.Any(), gomock.AssignableToTypeOf(&sts.GetFederationTokenInput{})).
					Return(nil, errTestSTS)
			},
			wantErrMsg: "sts GetFederationToken",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSTSClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSTSRepositoryWithInterface(mockClient)
			got, err := repo.GetFederationToken(context.Background(), tc.input)

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

func TestAWSSTSRepository_GetSessionToken(t *testing.T) {
	t.Parallel()

	in := &sts.GetSessionTokenInput{DurationSeconds: aws.Int32(3600)}

	tests := []struct {
		name       string
		input      *sts.GetSessionTokenInput
		setupMock  func(m *mocks.MockAWSSTSClientInterface)
		want       *sts.GetSessionTokenOutput
		wantErrMsg string
	}{
		{
			name:  "success",
			input: in,
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				expected := &sts.GetSessionTokenOutput{}
				m.EXPECT().GetSessionToken(gomock.Any(), in).Return(expected, nil)
			},
			want: &sts.GetSessionTokenOutput{},
		},
		{
			name:  "nil input uses empty request",
			input: nil,
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				m.EXPECT().
					GetSessionToken(gomock.Any(), gomock.AssignableToTypeOf(&sts.GetSessionTokenInput{})).
					Return(nil, errTestSTS)
			},
			wantErrMsg: "sts GetSessionToken",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSTSClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSTSRepositoryWithInterface(mockClient)
			got, err := repo.GetSessionToken(context.Background(), tc.input)

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

func TestAWSSTSRepository_GetWebIdentityToken(t *testing.T) {
	t.Parallel()

	in := &sts.GetWebIdentityTokenInput{
		Audience:         []string{"service-a"},
		SigningAlgorithm: aws.String("RS256"),
	}

	tests := []struct {
		name       string
		input      *sts.GetWebIdentityTokenInput
		setupMock  func(m *mocks.MockAWSSTSClientInterface)
		want       *sts.GetWebIdentityTokenOutput
		wantErrMsg string
	}{
		{
			name:  "success",
			input: in,
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				expected := &sts.GetWebIdentityTokenOutput{}
				m.EXPECT().GetWebIdentityToken(gomock.Any(), in).Return(expected, nil)
			},
			want: &sts.GetWebIdentityTokenOutput{},
		},
		{
			name:  "nil input uses empty request",
			input: nil,
			setupMock: func(m *mocks.MockAWSSTSClientInterface) {
				m.EXPECT().
					GetWebIdentityToken(gomock.Any(), gomock.AssignableToTypeOf(&sts.GetWebIdentityTokenInput{})).
					Return(nil, errTestSTS)
			},
			wantErrMsg: "sts GetWebIdentityToken",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSTSClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSSTSRepositoryWithInterface(mockClient)
			got, err := repo.GetWebIdentityToken(context.Background(), tc.input)

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
