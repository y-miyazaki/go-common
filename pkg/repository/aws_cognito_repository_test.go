package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestCognito = errors.New("cognito service error")

const (
	testUserPoolID           = "test_pool"
	testUserPoolClientID     = "test_client"
	testUserPoolClientSecret = "test_secret"
)

func TestNewAWSCognitoRepository(t *testing.T) {
	t.Parallel()

	mockClient := &cognitoidentityprovider.Client{}
	repo := NewAWSCognitoRepository(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
	require.Equal(t, testUserPoolID, repo.userPoolID)
	require.Equal(t, testUserPoolClientID, repo.userPoolClientID)
	require.Equal(t, testUserPoolClientSecret, repo.userPoolClientSecret)
}

func TestAWSCognitoRepository_GetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		username   string
		setupMock  func(m *mocks.MockAWSCognitoIdentityProviderClientInterface)
		want       *cognitoidentityprovider.AdminGetUserOutput
		wantErrMsg string
	}{
		{
			name:     "success",
			username: "test-user",
			setupMock: func(m *mocks.MockAWSCognitoIdentityProviderClientInterface) {
				expected := &cognitoidentityprovider.AdminGetUserOutput{
					UserAttributes: []types.AttributeType{
						{
							Name:  aws.String("email"),
							Value: aws.String("test@example.com"),
						},
					},
					Username: aws.String("test-user"),
				}
				m.EXPECT().
					AdminGetUser(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminGetUserInput{})).
					DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
						require.Equal(t, testUserPoolID, *input.UserPoolId)
						require.Equal(t, "test-user", *input.Username)
						return expected, nil
					})
			},
			want: &cognitoidentityprovider.AdminGetUserOutput{
				UserAttributes: []types.AttributeType{
					{
						Name:  aws.String("email"),
						Value: aws.String("test@example.com"),
					},
				},
				Username: aws.String("test-user"),
			},
		},
		{
			name:     "client error",
			username: "test-user",
			setupMock: func(m *mocks.MockAWSCognitoIdentityProviderClientInterface) {
				m.EXPECT().AdminGetUser(gomock.Any(), gomock.Any()).Return(nil, errTestCognito)
			},
			wantErrMsg: "cognito AdminGetUser",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)
			tc.setupMock(mockClient)

			repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
			got, err := repo.GetUser(context.Background(), tc.username)

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

func TestAWSCognitoRepository_CreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		AdminCreateUser(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminCreateUserInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminCreateUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, "test-user", *input.Username)
			return &cognitoidentityprovider.AdminCreateUserOutput{}, nil
		})

	mockClient.EXPECT().
		AdminSetUserPassword(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminSetUserPasswordInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminSetUserPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, "test-user", *input.Username)
			require.Equal(t, "test-password", *input.Password)
			require.True(t, input.Permanent)
			return &cognitoidentityprovider.AdminSetUserPasswordOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.CreateUser(context.Background(), "test-user", "test-password")

	require.NoError(t, err)
}

func TestAWSCognitoRepository_DeleteUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		AdminDeleteUser(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminDeleteUserInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminDeleteUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDeleteUserOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, "test-user", *input.Username)
			return &cognitoidentityprovider.AdminDeleteUserOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.DeleteUser(context.Background(), "test-user")

	require.NoError(t, err)
}

func TestAWSCognitoRepository_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	expectedSecretHash := repo.getSecretHash("test-user")

	authResult := &types.AuthenticationResultType{
		AccessToken:  aws.String("mock-access-token"),
		RefreshToken: aws.String("mock-refresh-token"),
		ExpiresIn:    3600,
	}

	mockClient.EXPECT().
		AdminInitiateAuth(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminInitiateAuthInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminInitiateAuthInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, testUserPoolClientID, *input.ClientId)
			require.Equal(t, types.AuthFlowTypeAdminUserPasswordAuth, input.AuthFlow)
			require.Equal(t, "test-user", input.AuthParameters["USERNAME"])
			require.Equal(t, "test-password", input.AuthParameters["PASSWORD"])
			require.Equal(t, expectedSecretHash, input.AuthParameters["SECRET_HASH"])
			return &cognitoidentityprovider.AdminInitiateAuthOutput{
				AuthenticationResult: authResult,
			}, nil
		})

	token, err := repo.Login(context.Background(), "test-user", "test-password")

	require.NoError(t, err)
	require.Equal(t, "mock-access-token", token.AccessToken)
	require.Equal(t, "mock-refresh-token", token.RefreshToken)
	require.True(t, token.AccessTokenExpiresAt.After(time.Now()))
}

func TestAWSCognitoRepository_ChangePassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		ChangePassword(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.ChangePasswordInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.ChangePasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ChangePasswordOutput, error) {
			require.Equal(t, "mock-token", *input.AccessToken)
			require.Equal(t, "old-password", *input.PreviousPassword)
			require.Equal(t, "new-password", *input.ProposedPassword)
			return &cognitoidentityprovider.ChangePasswordOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.ChangePassword(context.Background(), "Bearer mock-token", "old-password", "new-password")

	require.NoError(t, err)
}

func TestAWSCognitoRepository_ResetUserPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		AdminResetUserPassword(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminResetUserPasswordInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminResetUserPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminResetUserPasswordOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, "test-user", *input.Username)
			return &cognitoidentityprovider.AdminResetUserPasswordOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.ResetUserPassword(context.Background(), "test-user")

	require.NoError(t, err)
}

func TestAWSCognitoRepository_ConfirmForgotPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		ConfirmForgotPassword(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.ConfirmForgotPasswordInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
			require.Equal(t, testUserPoolClientID, *input.ClientId)
			require.Equal(t, "test-user", *input.Username)
			require.Equal(t, "new-password", *input.Password)
			require.Equal(t, "123456", *input.ConfirmationCode)
			return &cognitoidentityprovider.ConfirmForgotPasswordOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.ConfirmForgotPassword(context.Background(), "test-user", "new-password", "123456")

	require.NoError(t, err)
}

func TestAWSCognitoRepository_RefreshToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	expectedSecretHash := repo.getSecretHash("test-user")

	authResult := &types.AuthenticationResultType{
		AccessToken: aws.String("new-access-token"),
		ExpiresIn:   3600,
	}

	mockClient.EXPECT().
		AdminInitiateAuth(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminInitiateAuthInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminInitiateAuthInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, testUserPoolClientID, *input.ClientId)
			require.Equal(t, types.AuthFlowTypeRefreshTokenAuth, input.AuthFlow)
			require.Equal(t, "refresh-token", input.AuthParameters["REFRESH_TOKEN"])
			require.Equal(t, expectedSecretHash, input.AuthParameters["SECRET_HASH"])
			return &cognitoidentityprovider.AdminInitiateAuthOutput{
				AuthenticationResult: authResult,
			}, nil
		})

	token, err := repo.RefreshToken(context.Background(), "refresh-token", "test-user")

	require.NoError(t, err)
	require.Equal(t, "new-access-token", token.AccessToken)
	require.True(t, token.AccessTokenExpiresAt.After(time.Now()))
}

func TestAWSCognitoRepository_SetUserPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		AdminSetUserPassword(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.AdminSetUserPasswordInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.AdminSetUserPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error) {
			require.Equal(t, testUserPoolID, *input.UserPoolId)
			require.Equal(t, "test-user", *input.Username)
			require.Equal(t, "new-password", *input.Password)
			require.True(t, input.Permanent)
			return &cognitoidentityprovider.AdminSetUserPasswordOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.SetUserPassword(context.Background(), "test-user", "new-password", true)

	require.NoError(t, err)
}

func TestAWSCognitoRepository_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockAWSCognitoIdentityProviderClientInterface(ctrl)

	mockClient.EXPECT().
		RevokeToken(gomock.Any(), gomock.AssignableToTypeOf(&cognitoidentityprovider.RevokeTokenInput{})).
		DoAndReturn(func(_ context.Context, input *cognitoidentityprovider.RevokeTokenInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RevokeTokenOutput, error) {
			require.Equal(t, testUserPoolClientID, *input.ClientId)
			require.Equal(t, "refresh-token", *input.Token)
			require.Equal(t, testUserPoolClientSecret, *input.ClientSecret)
			return &cognitoidentityprovider.RevokeTokenOutput{}, nil
		})

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, testUserPoolID, testUserPoolClientID, testUserPoolClientSecret)
	err := repo.Logout(context.Background(), "refresh-token")

	require.NoError(t, err)
}

func TestAWSCognitoRepositoryReal_Logout(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Cognito integration test - requires real AWS credentials")
}

func TestAWSCognitoRepositoryReal_RefreshToken(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Cognito integration test - requires real AWS credentials")
}

func TestAWSCognitoRepositoryReal_SetUserPassword(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Cognito integration test - requires real AWS credentials")
}

func TestAWSCognitoRepository_GetSecretHash(t *testing.T) {
	t.Parallel()

	repo := &AWSCognitoRepository{
		userPoolClientID:     "test-client-id",
		userPoolClientSecret: "test-client-secret",
	}

	hash := repo.getSecretHash("testuser")

	require.NotEmpty(t, hash)
	require.Positive(t, len(hash))
}

func TestAWSCognitoRepository_GetAccessToken(t *testing.T) {
	t.Parallel()

	repo := &AWSCognitoRepository{}

	tests := []struct {
		errorType           error
		name                string
		authorizationHeader string
		expectedToken       string
		expectError         bool
	}{
		{
			name:                "valid bearer token",
			authorizationHeader: "Bearer test-token-123",
			expectedToken:       "test-token-123",
			expectError:         false,
		},
		{
			name:                "empty header",
			authorizationHeader: "",
			expectedToken:       "",
			expectError:         true,
			errorType:           ErrAWSCognitoAccessTokenNotFound,
		},
		{
			name:                "invalid format",
			authorizationHeader: "InvalidFormat",
			expectedToken:       "",
			expectError:         true,
			errorType:           ErrAWSCognitoAccessTokenFormatNotSupported,
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			token, err := repo.getAccessToken(tc.authorizationHeader)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorType.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestAWSCognitoToken_Struct(t *testing.T) {
	t.Parallel()

	token := AWSCognitoToken{
		AccessToken:          "access-token-123",
		AccessTokenExpiresAt: time.Now().Add(time.Hour),
		RefreshToken:         "refresh-token-456",
	}

	require.Equal(t, "access-token-123", token.AccessToken)
	require.Equal(t, "refresh-token-456", token.RefreshToken)
	require.True(t, token.AccessTokenExpiresAt.After(time.Now()))
}
