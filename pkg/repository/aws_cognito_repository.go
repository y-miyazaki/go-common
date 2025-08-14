// pragma: allowlist-secret
// NOTE: This file may contain references to configuration secrets (client secret, passwords) that are
// provided via environment or secret manager at runtime. The occurrences are documented with inline
// allowlist comments. To avoid false positives from static secret scanners during validation, mark
// the file as intentional for secret references.

// Package repository provides AWS SDK v2 based repository implementations.
package repository

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var (
	ErrAWSCognitoAccessTokenNotFound           = errors.New("access token not found")
	ErrAWSCognitoAccessTokenFormatNotSupported = errors.New("access token format is not supported")
)

// AWSCognitoRepositoryInterface interface.
// nolint:iface,revive,unused
type AWSCognitoRepositoryInterface interface {
	GetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error)
	CreateUser(username, password string) error
	DeleteUser(username string) error
	Login(username, password string) (AWSCognitoToken, error)
	RefreshToken(refreshToken, username string) (AWSCognitoToken, error)
	ChangePassword(authorizationHeader, previousPassword, proposedPassword string) error
	ResetUserPassword(username string) error
	ConfirmForgotPassword(username, password, confirmationCode string) error
	SetUserPassword(username, password string, permanent bool) error
	Logout(refreshToken string) error
}

// AWSCognitoRepository struct.
type AWSCognitoRepository struct {
	c                    *cognitoidentityprovider.Client
	userPoolID           string
	userPoolClientID     string
	userPoolClientSecret string
}

// AWSCognitoToken struct
// nolint:tagliatelle
type AWSCognitoToken struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken         string    `json:"refresh_token"`
}

// NewAWSCognitoRepository returns AWSCognitoRepository instance.
func NewAWSCognitoRepository(c *cognitoidentityprovider.Client, userPoolID, userPoolClientID, userPoolClientSecret string) *AWSCognitoRepository {
	return &AWSCognitoRepository{
		c:                    c,
		userPoolID:           userPoolID,
		userPoolClientID:     userPoolClientID,
		userPoolClientSecret: userPoolClientSecret, // pragma: allowlist secret
	}
}

// GetUser gets a user information from Cognito.
func (r *AWSCognitoRepository) GetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	res, err := r.c.AdminGetUser(context.TODO(), &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return nil, fmt.Errorf("cognito AdminGetUser: %w", err)
	}
	return res, nil
}

// CreateUser creates a new user for Cognito user pool.
func (r *AWSCognitoRepository) CreateUser(username, password string) error {
	_, err := r.c.AdminCreateUser(context.TODO(), &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("cognito AdminCreateUser: %w", err)
	}

	_, err = r.c.AdminSetUserPassword(context.TODO(), &cognitoidentityprovider.AdminSetUserPasswordInput{
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
func (r *AWSCognitoRepository) DeleteUser(username string) error {
	_, err := r.c.AdminDeleteUser(context.TODO(), &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("cognito AdminDeleteUser: %w", err)
	}
	return nil
}

// Login logs to Cognito.
func (r *AWSCognitoRepository) Login(username, password string) (AWSCognitoToken, error) {
	res, err := r.c.AdminInitiateAuth(context.TODO(), &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   types.AuthFlowTypeAdminUserPasswordAuth,
		ClientId:   aws.String(r.userPoolClientID),
		UserPoolId: aws.String(r.userPoolID),
		AuthParameters: map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password, // pragma: allowlist secret
			"SECRET_HASH": r.getSecretHash(username),
		},
	})
	if err != nil {
		return AWSCognitoToken{}, fmt.Errorf("cognito AdminInitiateAuth: %w", err)
	}

	return AWSCognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Second * time.Duration(res.AuthenticationResult.ExpiresIn)),
		RefreshToken:         *res.AuthenticationResult.RefreshToken,
	}, nil
}

// Logout logs out of Cognito.
func (r *AWSCognitoRepository) Logout(refreshToken string) error {
	_, err := r.c.RevokeToken(context.TODO(), &cognitoidentityprovider.RevokeTokenInput{
		ClientId:     aws.String(r.userPoolClientID),
		Token:        aws.String(refreshToken),
		ClientSecret: aws.String(r.userPoolClientSecret),
	})
	if err != nil {
		return fmt.Errorf("cognito RevokeToken: %w", err)
	}
	return nil
}

func (r *AWSCognitoRepository) RefreshToken(refreshToken, username string) (AWSCognitoToken, error) {
	res, err := r.c.AdminInitiateAuth(context.TODO(), &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   types.AuthFlowTypeRefreshTokenAuth,
		ClientId:   aws.String(r.userPoolClientID),
		UserPoolId: aws.String(r.userPoolID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   r.getSecretHash(username),
		},
	})
	if err != nil {
		return AWSCognitoToken{}, fmt.Errorf("cognito AdminInitiateAuth(refresh): %w", err)
	}

	return AWSCognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Second * time.Duration(res.AuthenticationResult.ExpiresIn)),
	}, nil
}

// SetUserPassword sets the password of the user.
func (r *AWSCognitoRepository) SetUserPassword(username, password string, permanent bool) error {
	_, err := r.c.AdminSetUserPassword(context.TODO(), &cognitoidentityprovider.AdminSetUserPasswordInput{
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

// ChangePassword changes the password of the user.
func (r *AWSCognitoRepository) ChangePassword(authorizationHeader, previousPassword, proposedPassword string) error {
	accessToken, err := r.getAccessToken(authorizationHeader)
	if err != nil {
		return fmt.Errorf("cognito getAccessToken: %w", err)
	}

	_, err = r.c.ChangePassword(context.TODO(), &cognitoidentityprovider.ChangePasswordInput{
		PreviousPassword: aws.String(previousPassword),
		ProposedPassword: aws.String(proposedPassword),
		AccessToken:      aws.String(accessToken),
	})
	if err != nil {
		return fmt.Errorf("cognito ChangePassword: %w", err)
	}
	return nil
}

// ResetPassword resets the specified user's password in a user pool as an administrator. Works on any user.
func (r *AWSCognitoRepository) ResetUserPassword(username string) error {
	_, err := r.c.AdminResetUserPassword(context.TODO(), &cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("cognito AdminResetUserPassword: %w", err)
	}
	return nil
}

// ConfirmForgotPassword allows a user to enter a confirmation code to reset a forgotten password.
func (r *AWSCognitoRepository) ConfirmForgotPassword(username, password, confirmationCode string) error {
	_, err := r.c.ConfirmForgotPassword(context.TODO(), &cognitoidentityprovider.ConfirmForgotPasswordInput{
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

// getSecretHash gets the secret hash.
func (r *AWSCognitoRepository) getSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(r.userPoolClientSecret))
	if _, err := mac.Write([]byte(username + r.userPoolClientID)); err != nil {
		// If HMAC write fails, return empty string so callers will fail fast.
		return ""
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// getSecretHash gets the access token.
func (*AWSCognitoRepository) getAccessToken(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", ErrAWSCognitoAccessTokenNotFound
	}

	numberOfElementsInArray := 2
	prefixBearer := "Bearer "
	stringArr := strings.Split(authorizationHeader, prefixBearer)
	if len(stringArr) == numberOfElementsInArray {
		return stringArr[1], nil
	}
	return "", fmt.Errorf("%w", ErrAWSCognitoAccessTokenFormatNotSupported)
}
