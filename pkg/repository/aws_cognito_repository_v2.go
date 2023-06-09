package repository

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var (
	AWSV2CognitoAccessTokenNotFound           = errors.New("access token not found")
	AWSV2CognitoAccessTokenFormatNotSupported = errors.New("access token format is not supported")
)

// AWSV2CognitoRepositoryInterface interface.
type AWSV2CognitoRepositoryInterface interface {
	GetUser(username, password string) (*cognitoidentityprovider.AdminGetUserOutput, error)
	CreateUser(username, password string) error
	DeleteUser(username string) error
	Login(username, password string) (AWSV2CognitoToken, error)
	Logout(username string) error
	RefreshToken(refreshToken string, username string) (AWSV2CognitoToken, error)
	SetUserPassword(username string, password string) error
	ChangePassword(authorizationHeader, previousPassword, proposedPassword string) error
	ResetUserPassword(username string) error
	ConfirmForgotPassword(username, password, confirmationCode string) error
}

// AWSV2CognitoRepository struct.
type AWSV2CognitoRepository struct {
	c                    *cognitoidentityprovider.Client
	userPoolID           string
	userPoolClientID     string
	userPoolClientSecret string
}

// AWSV2CognitoToken struct
type AWSV2CognitoToken struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken         string    `json:"refresh_token"`
}

// NewAWSV2CognitoRepository returns AWSV2CognitoRepository instance.
func NewAWSV2CognitoRepository(c *cognitoidentityprovider.Client, userPoolID, userPoolClientID, userPoolClientSecret string) *AWSV2CognitoRepository {
	return &AWSV2CognitoRepository{
		c:                    c,
		userPoolID:           userPoolID,
		userPoolClientID:     userPoolClientID,
		userPoolClientSecret: userPoolClientSecret, // pragma: allowlist secret
	}
}

// GetUser gets a user information from Cognito.
func (r *AWSV2CognitoRepository) GetUser(username, password string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	return r.c.AdminGetUser(context.TODO(), &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
}

// CreateUser creates a new user for Cognito user pool.
func (r *AWSV2CognitoRepository) CreateUser(username string, password string) error {
	_, err := r.c.AdminCreateUser(context.TODO(), &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return err
	}

	_, err = r.c.AdminSetUserPassword(context.TODO(), &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  true,
	})
	return err
}

// DeleteUser deletes a user from Cognito user pool.
func (r *AWSV2CognitoRepository) DeleteUser(username string) error {
	_, err := r.c.AdminDeleteUser(context.TODO(), &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	return err
}

// Login logs to Cognito.
func (r *AWSV2CognitoRepository) Login(username, password string) (AWSV2CognitoToken, error) {
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
		return AWSV2CognitoToken{}, err
	}

	return AWSV2CognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Second * time.Duration(res.AuthenticationResult.ExpiresIn)),
		RefreshToken:         *res.AuthenticationResult.RefreshToken,
	}, nil
}

// Logout logs out of Cognito.
func (r *AWSV2CognitoRepository) Logout(refreshToken string) error {
	_, err := r.c.RevokeToken(context.TODO(), &cognitoidentityprovider.RevokeTokenInput{
		ClientId:     aws.String(r.userPoolClientID),
		Token:        aws.String(refreshToken),
		ClientSecret: aws.String(r.userPoolClientSecret),
	})
	return err
}

func (r *AWSV2CognitoRepository) RefreshToken(refreshToken, username string) (AWSV2CognitoToken, error) {
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
		return AWSV2CognitoToken{}, err
	}

	return AWSV2CognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Second * time.Duration(res.AuthenticationResult.ExpiresIn)),
	}, nil
}

// SetUserPassword sets the password of the user.
func (r *AWSV2CognitoRepository) SetUserPassword(username, password string, permanent bool) error {
	_, err := r.c.AdminSetUserPassword(context.TODO(), &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  permanent,
	})

	return err
}

// ChangePassword changes the password of the user.
func (r *AWSV2CognitoRepository) ChangePassword(authorizationHeader, previousPassword, proposedPassword string) error {
	accessToken, err := r.getAccessToken(authorizationHeader)
	if err != nil {
		return err
	}

	_, err = r.c.ChangePassword(context.TODO(), &cognitoidentityprovider.ChangePasswordInput{
		PreviousPassword: aws.String(previousPassword),
		ProposedPassword: aws.String(proposedPassword),
		AccessToken:      aws.String(accessToken),
	})
	return err
}

// ResetPassword resets the specified user's password in a user pool as an administrator. Works on any user.
func (r *AWSV2CognitoRepository) ResetUserPassword(username string) error {
	_, err := r.c.AdminResetUserPassword(context.TODO(), &cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	return err
}

// ConfirmForgotPassword allows a user to enter a confirmation code to reset a forgotten password.
func (r *AWSV2CognitoRepository) ConfirmForgotPassword(username, password, confirmationCode string) error {
	_, err := r.c.ConfirmForgotPassword(context.TODO(), &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(r.userPoolClientID),
		Username:         aws.String(username),
		Password:         aws.String(password),
		ConfirmationCode: aws.String(confirmationCode),
	})
	return err
}

// getSecretHash gets the secret hash.
func (r *AWSV2CognitoRepository) getSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(r.userPoolClientSecret))
	mac.Write([]byte(username + r.userPoolClientID))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// getSecretHash gets the access token.
func (r *AWSV2CognitoRepository) getAccessToken(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", AWSV2CognitoAccessTokenNotFound
	}

	numberOfElementsInArray := 2
	prefixBearer := "Bearer "
	stringArr := strings.Split(authorizationHeader, prefixBearer)
	if len(stringArr) == numberOfElementsInArray {
		return stringArr[1], nil
	}
	return "", AWSV2CognitoAccessTokenFormatNotSupported
}
