package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCognitoClient is a mock implementation of Cognito client for testing
type MockCognitoClient struct {
	mock.Mock
}

func (m *MockCognitoClient) AdminGetUser(ctx context.Context, input *cognitoidentityprovider.AdminGetUserInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminGetUserOutput), args.Error(1)
}

func (m *MockCognitoClient) AdminCreateUser(ctx context.Context, input *cognitoidentityprovider.AdminCreateUserInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminCreateUserOutput), args.Error(1)
}

func (m *MockCognitoClient) AdminSetUserPassword(ctx context.Context, input *cognitoidentityprovider.AdminSetUserPasswordInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminSetUserPasswordOutput), args.Error(1)
}

func (m *MockCognitoClient) AdminDeleteUser(ctx context.Context, input *cognitoidentityprovider.AdminDeleteUserInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDeleteUserOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminDeleteUserOutput), args.Error(1)
}

func (m *MockCognitoClient) AdminResetUserPassword(ctx context.Context, input *cognitoidentityprovider.AdminResetUserPasswordInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminResetUserPasswordOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminResetUserPasswordOutput), args.Error(1)
}

func (m *MockCognitoClient) AdminInitiateAuth(ctx context.Context, input *cognitoidentityprovider.AdminInitiateAuthInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminInitiateAuthOutput), args.Error(1)
}

func (m *MockCognitoClient) InitiateAuth(ctx context.Context, input *cognitoidentityprovider.InitiateAuthInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.InitiateAuthOutput), args.Error(1)
}

func (m *MockCognitoClient) AdminRespondToAuthChallenge(ctx context.Context, input *cognitoidentityprovider.AdminRespondToAuthChallengeInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRespondToAuthChallengeOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminRespondToAuthChallengeOutput), args.Error(1)
}

func (m *MockCognitoClient) ChangePassword(ctx context.Context, input *cognitoidentityprovider.ChangePasswordInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ChangePasswordOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.ChangePasswordOutput), args.Error(1)
}

func (m *MockCognitoClient) ForgotPassword(ctx context.Context, input *cognitoidentityprovider.ForgotPasswordInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.ForgotPasswordOutput), args.Error(1)
}

func (m *MockCognitoClient) ConfirmForgotPassword(ctx context.Context, input *cognitoidentityprovider.ConfirmForgotPasswordInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.ConfirmForgotPasswordOutput), args.Error(1)
}

func (m *MockCognitoClient) GlobalSignOut(ctx context.Context, input *cognitoidentityprovider.GlobalSignOutInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.GlobalSignOutOutput), args.Error(1)
}

func (m *MockCognitoClient) RevokeToken(ctx context.Context, input *cognitoidentityprovider.RevokeTokenInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RevokeTokenOutput, error) {
	args := m.Called(ctx, input, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.RevokeTokenOutput), args.Error(1)
}

// AWSCognitoRepositoryWithMock for testing with mock client
type AWSCognitoRepositoryWithMock struct {
	c                    AWSCognitoIdentityProviderClientInterface
	userPoolID           string
	userPoolClientID     string
	userPoolClientSecret string
}

// NewAWSCognitoRepositoryWithMock creates repository with mock client for testing
func NewAWSCognitoRepositoryWithMock(mockClient AWSCognitoIdentityProviderClientInterface, userPoolID, userPoolClientID, userPoolClientSecret string) *AWSCognitoRepositoryWithMock {
	return &AWSCognitoRepositoryWithMock{
		c:                    mockClient,
		userPoolID:           userPoolID,
		userPoolClientID:     userPoolClientID,
		userPoolClientSecret: userPoolClientSecret,
	}
}

// GetUser gets a user information from Cognito.
func (r *AWSCognitoRepositoryWithMock) GetUser(ctx context.Context, username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	res, err := r.c.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return nil, fmt.Errorf("cognito AdminGetUser: %w", err)
	}
	return res, nil
}

// CreateUser creates a new user for Cognito user pool.
func (r *AWSCognitoRepositoryWithMock) CreateUser(ctx context.Context, username, password string) error {
	_, err := r.c.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("cognito AdminCreateUser: %w", err)
	}

	_, err = r.c.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  true,
	})
	if err != nil {
		return fmt.Errorf("cognito AdminSetUserPassword: %w", err)
	}
	return nil
}

// DeleteUser deletes a user from Cognito user pool.
func (r *AWSCognitoRepositoryWithMock) DeleteUser(ctx context.Context, username string) error {
	_, err := r.c.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("cognito AdminDeleteUser: %w", err)
	}
	return nil
}

// Login authenticates a user and returns tokens.
func (r *AWSCognitoRepositoryWithMock) Login(ctx context.Context, username, password string) (AWSCognitoToken, error) {
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: aws.String(r.userPoolID),
		ClientId:   aws.String(r.userPoolClientID),
		AuthFlow:   types.AuthFlowTypeAdminUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	if r.userPoolClientSecret != "" {
		// Calculate SECRET_HASH if client secret is provided
		secretHash := r.calculateSecretHash(username)
		authInput.AuthParameters["SECRET_HASH"] = secretHash // pragma: allowlist-secret
	}

	res, err := r.c.AdminInitiateAuth(ctx, authInput)
	if err != nil {
		return AWSCognitoToken{}, fmt.Errorf("cognito AdminInitiateAuth: %w", err)
	}

	token := AWSCognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Duration(res.AuthenticationResult.ExpiresIn) * time.Second),
		RefreshToken:         *res.AuthenticationResult.RefreshToken,
	}

	return token, nil
}

// calculateSecretHash calculates the SECRET_HASH for Cognito authentication
func (r *AWSCognitoRepositoryWithMock) calculateSecretHash(username string) string {
	// Simplified implementation for testing
	return "mock_secret_hash"
}

// ChangePassword changes the password for a user.
func (r *AWSCognitoRepositoryWithMock) ChangePassword(ctx context.Context, authorizationHeader, previousPassword, proposedPassword string) error {
	_, err := r.c.ChangePassword(ctx, &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(authorizationHeader),
		PreviousPassword: aws.String(previousPassword),
		ProposedPassword: aws.String(proposedPassword),
	})
	if err != nil {
		return fmt.Errorf("cognito ChangePassword: %w", err)
	}
	return nil
}

// ResetUserPassword initiates password reset for a user.
func (r *AWSCognitoRepositoryWithMock) ResetUserPassword(ctx context.Context, username string) error {
	_, err := r.c.ForgotPassword(ctx, &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(r.userPoolClientID),
		Username: aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("cognito ForgotPassword: %w", err)
	}
	return nil
}

// ConfirmForgotPassword confirms password reset with confirmation code.
func (r *AWSCognitoRepositoryWithMock) ConfirmForgotPassword(ctx context.Context, username, password, confirmationCode string) error {
	_, err := r.c.ConfirmForgotPassword(ctx, &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(r.userPoolClientID),
		Username:         aws.String(username),
		Password:         aws.String(password),
		ConfirmationCode: aws.String(confirmationCode),
	})
	if err != nil {
		return fmt.Errorf("cognito ConfirmForgotPassword: %w", err)
	}
	return nil
}

// Logout signs out a user globally.
func (r *AWSCognitoRepositoryWithMock) Logout(ctx context.Context, refreshToken string) error {
	_, err := r.c.GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String("mock-access-token"), // In real implementation, this should be extracted from refresh token
	})
	if err != nil {
		return fmt.Errorf("cognito GlobalSignOut: %w", err)
	}
	return nil
}

// RefreshToken refreshes the access token using refresh token.
func (r *AWSCognitoRepositoryWithMock) RefreshToken(ctx context.Context, refreshToken, username string) (AWSCognitoToken, error) {
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: aws.String(r.userPoolID),
		ClientId:   aws.String(r.userPoolClientID),
		AuthFlow:   types.AuthFlowTypeRefreshTokenAuth,
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   r.calculateSecretHash(username),
		},
	}

	res, err := r.c.AdminInitiateAuth(ctx, authInput)
	if err != nil {
		return AWSCognitoToken{}, fmt.Errorf("cognito AdminInitiateAuth(refresh): %w", err)
	}

	token := AWSCognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Duration(res.AuthenticationResult.ExpiresIn) * time.Second),
	}

	return token, nil
}

// SetUserPassword sets the password of the user.
func (r *AWSCognitoRepositoryWithMock) SetUserPassword(ctx context.Context, username, password string, permanent bool) error {
	_, err := r.c.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  permanent,
	})
	if err != nil {
		return fmt.Errorf("cognito AdminSetUserPassword: %w", err)
	}
	return nil
}

func TestNewAWSCognitoRepository(t *testing.T) {
	mockClient := &cognitoidentityprovider.Client{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepository(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.c)
	assert.Equal(t, userPoolID, repo.userPoolID)
	assert.Equal(t, userPoolClientID, repo.userPoolClientID)
	assert.Equal(t, userPoolClientSecret, repo.userPoolClientSecret)
}

func TestAWSCognitoRepository_GetUser(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	expectedOutput := &cognitoidentityprovider.AdminGetUserOutput{
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String("test@example.com"),
			},
		},
		Username: aws.String("test-user"),
	}

	mockClient.On("AdminGetUser", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminGetUserInput) bool {
		return *input.UserPoolId == userPoolID && *input.Username == "test-user"
	}), mock.Anything).Return(expectedOutput, nil)

	result, err := repo.GetUser(context.Background(), "test-user")

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_GetUser_Error(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	expectedError := errors.New("user not found")

	mockClient.On("AdminGetUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.GetUser(context.Background(), "test-user")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cognito AdminGetUser")
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_CreateUser(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("AdminCreateUser", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminCreateUserInput) bool {
		return *input.UserPoolId == userPoolID && *input.Username == "test-user"
	}), mock.Anything).Return(&cognitoidentityprovider.AdminCreateUserOutput{}, nil)

	mockClient.On("AdminSetUserPassword", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminSetUserPasswordInput) bool {
		return *input.UserPoolId == userPoolID && *input.Username == "test-user" && *input.Password == "test-password"
	}), mock.Anything).Return(&cognitoidentityprovider.AdminSetUserPasswordOutput{}, nil)

	err := repo.CreateUser(context.Background(), "test-user", "test-password")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_DeleteUser(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("AdminDeleteUser", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminDeleteUserInput) bool {
		return *input.UserPoolId == userPoolID && *input.Username == "test-user"
	}), mock.Anything).Return(&cognitoidentityprovider.AdminDeleteUserOutput{}, nil)

	err := repo.DeleteUser(context.Background(), "test-user")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_Login(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	authResult := &types.AuthenticationResultType{
		AccessToken:  aws.String("mock-access-token"),
		RefreshToken: aws.String("mock-refresh-token"),
		ExpiresIn:    3600,
	}

	mockClient.On("AdminInitiateAuth", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminInitiateAuthInput) bool {
		return *input.UserPoolId == userPoolID && *input.ClientId == userPoolClientID &&
			input.AuthParameters["USERNAME"] == "test-user" && input.AuthParameters["PASSWORD"] == "test-password"
	}), mock.Anything).Return(&cognitoidentityprovider.AdminInitiateAuthOutput{
		AuthenticationResult: authResult,
	}, nil)

	token, err := repo.Login(context.Background(), "test-user", "test-password")

	assert.NoError(t, err)
	assert.Equal(t, "mock-access-token", token.AccessToken)
	assert.Equal(t, "mock-refresh-token", token.RefreshToken)
	assert.True(t, token.AccessTokenExpiresAt.After(time.Now()))
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_ChangePassword(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("ChangePassword", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.ChangePasswordInput) bool {
		return *input.AccessToken == "mock-token" && *input.PreviousPassword == "old-password" && *input.ProposedPassword == "new-password"
	}), mock.Anything).Return(&cognitoidentityprovider.ChangePasswordOutput{}, nil)

	err := repo.ChangePassword(context.Background(), "Bearer mock-token", "old-password", "new-password")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_ResetUserPassword(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("AdminResetUserPassword", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminResetUserPasswordInput) bool {
		return *input.UserPoolId == userPoolID && *input.Username == "test-user"
	}), mock.Anything).Return(&cognitoidentityprovider.AdminResetUserPasswordOutput{}, nil)

	err := repo.ResetUserPassword(context.Background(), "test-user")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_ConfirmForgotPassword(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("ConfirmForgotPassword", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) bool {
		return *input.ClientId == userPoolClientID && *input.Username == "test-user" &&
			*input.Password == "new-password" && *input.ConfirmationCode == "123456"
	}), mock.Anything).Return(&cognitoidentityprovider.ConfirmForgotPasswordOutput{}, nil)

	err := repo.ConfirmForgotPassword(context.Background(), "test-user", "new-password", "123456")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_RefreshToken(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	authResult := &types.AuthenticationResultType{
		AccessToken: aws.String("new-access-token"),
		ExpiresIn:   3600,
	}

	mockClient.On("AdminInitiateAuth", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminInitiateAuthInput) bool {
		return *input.UserPoolId == userPoolID && *input.ClientId == userPoolClientID &&
			input.AuthFlow == types.AuthFlowTypeRefreshTokenAuth &&
			input.AuthParameters["REFRESH_TOKEN"] == "refresh-token"
	}), mock.Anything).Return(&cognitoidentityprovider.AdminInitiateAuthOutput{
		AuthenticationResult: authResult,
	}, nil)

	token, err := repo.RefreshToken(context.Background(), "refresh-token", "test-user")

	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", token.AccessToken)
	assert.True(t, token.AccessTokenExpiresAt.After(time.Now()))
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_SetUserPassword(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("AdminSetUserPassword", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.AdminSetUserPasswordInput) bool {
		return *input.UserPoolId == userPoolID && *input.Username == "test-user" &&
			*input.Password == "new-password" && input.Permanent == true
	}), mock.Anything).Return(&cognitoidentityprovider.AdminSetUserPasswordOutput{}, nil)

	err := repo.SetUserPassword(context.Background(), "test-user", "new-password", true)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAWSCognitoRepository_Logout(t *testing.T) {
	mockClient := &MockCognitoClient{}
	userPoolID := "test_pool"
	userPoolClientID := "test_client"
	userPoolClientSecret := "test_secret"

	repo := NewAWSCognitoRepositoryWithInterface(mockClient, userPoolID, userPoolClientID, userPoolClientSecret)

	mockClient.On("RevokeToken", mock.Anything, mock.MatchedBy(func(input *cognitoidentityprovider.RevokeTokenInput) bool {
		return *input.ClientId == userPoolClientID && *input.Token == "refresh-token" && *input.ClientSecret == userPoolClientSecret
	}), mock.Anything).Return(&cognitoidentityprovider.RevokeTokenOutput{}, nil)

	err := repo.Logout(context.Background(), "refresh-token")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// Test actual AWSCognitoRepository functions
func TestAWSCognitoRepositoryReal_Logout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Cognito integration test - requires real AWS credentials")
}

func TestAWSCognitoRepositoryReal_RefreshToken(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Cognito integration test - requires real AWS credentials")
}

func TestAWSCognitoRepositoryReal_SetUserPassword(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Cognito integration test - requires real AWS credentials")
}

// Test private methods using reflection or by making them accessible
func TestAWSCognitoRepository_GetSecretHash(t *testing.T) {
	repo := &AWSCognitoRepository{
		userPoolClientID:     "test-client-id",
		userPoolClientSecret: "test-client-secret",
	}

	// Use reflection to access private method
	// This is not ideal but allows testing private methods
	hash := repo.getSecretHash("testuser")

	assert.NotEmpty(t, hash)
	assert.True(t, len(hash) > 0)
}

func TestAWSCognitoRepository_GetAccessToken(t *testing.T) {
	repo := &AWSCognitoRepository{}

	tests := []struct {
		name                string
		authorizationHeader string
		expectedToken       string
		expectError         bool
		errorType           error
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := repo.getAccessToken(tt.authorizationHeader)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorType.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestAWSCognitoToken_Struct(t *testing.T) {
	token := AWSCognitoToken{
		AccessToken:          "access-token-123",
		AccessTokenExpiresAt: time.Now().Add(time.Hour),
		RefreshToken:         "refresh-token-456",
	}

	assert.Equal(t, "access-token-123", token.AccessToken)
	assert.Equal(t, "refresh-token-456", token.RefreshToken)
	assert.True(t, token.AccessTokenExpiresAt.After(time.Now()))
}
