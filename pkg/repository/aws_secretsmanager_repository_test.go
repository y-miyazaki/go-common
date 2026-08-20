package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestSecretsManager = errors.New("secretsmanager error")

func TestNewAWSSecretsManagerRepository(t *testing.T) {
	t.Parallel()

	mockClient := &secretsmanager.Client{}
	repo := NewAWSSecretsManagerRepository(mockClient)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
}

func TestAWSSecretsManagerRepository_GetSecretString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		secretName string
		setupMock  func(m *mocks.MockAWSSecretsManagerClientInterface)
		want       string
		wantErrMsg string
	}{
		{
			name:       "secret string",
			secretName: "test-secret",
			setupMock: func(m *mocks.MockAWSSecretsManagerClientInterface) {
				m.EXPECT().
					GetSecretValue(gomock.Any(), gomock.AssignableToTypeOf(&secretsmanager.GetSecretValueInput{})).
					DoAndReturn(func(_ context.Context, input *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
						require.NotNil(t, input.SecretId)
						require.Equal(t, "test-secret", *input.SecretId)
						require.NotNil(t, input.VersionStage)
						require.Equal(t, "AWSCURRENT", *input.VersionStage)
						return &secretsmanager.GetSecretValueOutput{
							SecretString: aws.String("secret-value"),
						}, nil
					})
			},
			want: "secret-value",
		},
		{
			name:       "secret binary",
			secretName: "test-secret-binary",
			setupMock: func(m *mocks.MockAWSSecretsManagerClientInterface) {
				secretValue := "binary-secret-value"
				encodedSecret := base64.StdEncoding.EncodeToString([]byte(secretValue))
				m.EXPECT().
					GetSecretValue(gomock.Any(), gomock.AssignableToTypeOf(&secretsmanager.GetSecretValueInput{})).
					DoAndReturn(func(_ context.Context, input *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
						require.NotNil(t, input.SecretId)
						require.Equal(t, "test-secret-binary", *input.SecretId)
						require.NotNil(t, input.VersionStage)
						require.Equal(t, "AWSCURRENT", *input.VersionStage)
						return &secretsmanager.GetSecretValueOutput{
							SecretString: nil,
							SecretBinary: []byte(encodedSecret),
						}, nil
					})
			},
			want: "binary-secret-value",
		},
		{
			name:       "client error",
			secretName: "test-secret-error",
			setupMock: func(m *mocks.MockAWSSecretsManagerClientInterface) {
				m.EXPECT().
					GetSecretValue(gomock.Any(), gomock.AssignableToTypeOf(&secretsmanager.GetSecretValueInput{})).
					DoAndReturn(func(_ context.Context, input *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
						require.NotNil(t, input.SecretId)
						require.Equal(t, "test-secret-error", *input.SecretId)
						require.NotNil(t, input.VersionStage)
						require.Equal(t, "AWSCURRENT", *input.VersionStage)
						return nil, errTestSecretsManager
					})
			},
			wantErrMsg: "secretsmanager GetSecretValue",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSSecretsManagerClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSSecretsManagerRepositoryWithInterface(mockClient)
			got, err := repo.GetSecretString(context.Background(), tc.secretName)

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				require.Empty(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
