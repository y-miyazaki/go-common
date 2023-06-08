package repository

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var (
	AWSCognitoAccessTokenNotFound           = errors.New("access token not found")
	AWSCognitoAccessTokenFormatNotSupported = errors.New("access token format is not supported")
)

// AWSCognitoRepositoryInterface interface.
type AWSCognitoRepositoryInterface interface {
	GetUser(username, password string) (*cognitoidentityprovider.AdminGetUserOutput, error)
	CreateUser(username, password string) error
	DeleteUser(username string) error
	Login(username, password string) (AWSCognitoToken, error)
	Logout(username string) error
	RefreshToken(refreshToken string, username string) (AWSCognitoToken, error)
	SetUserPassword(username string, password string) error
	ChangePassword(authorizationHeader, previousPassword, proposedPassword string) error
	ResetUserPassword(username string) error
	ConfirmForgotPassword(username, password, confirmationCode string) error
}

// AWSCognitoRepository struct.
type AWSCognitoRepository struct {
	c                    *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID           string
	userPoolClientID     string
	userPoolClientSecret string
}

// AWSCognitoToken struct
type AWSCognitoToken struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken         string    `json:"refresh_token"`
}

// NewAWSCognitoRepository returns AWSCognitoRepository instance.
func NewAWSCognitoRepository(c *cognitoidentityprovider.CognitoIdentityProvider, userPoolID, userPoolClientID, userPoolClientSecret string) *AWSCognitoRepository {
	return &AWSCognitoRepository{
		c:                    c,
		userPoolID:           userPoolID,
		userPoolClientID:     userPoolClientID,
		userPoolClientSecret: userPoolClientSecret, // pragma: allowlist secret
	}
}

// GetUser gets a user information from Cognito.
func (r *AWSCognitoRepository) GetUser(username, password string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	output, err := r.c.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return nil, err
	}
	return output, err
}

// CreateUser creates a new user for Cognito user pool.
func (r *AWSCognitoRepository) CreateUser(username string, password string) error {
	_, err := r.c.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		return err
	}

	_, err = r.c.AdminSetUserPassword(&cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  aws.Bool(true),
	})
	return err
}

// DeleteUser deletes a user from Cognito user pool.
func (r *AWSCognitoRepository) DeleteUser(username string) error {
	_, err := r.c.AdminDeleteUser(&cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	return err
}

// Login logs to Cognito.
func (r *AWSCognitoRepository) Login(username, password string) (AWSCognitoToken, error) {
	res, err := r.c.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   aws.String("ADMIN_USER_PASSWORD_AUTH"),
		ClientId:   aws.String(r.userPoolClientID),
		UserPoolId: aws.String(r.userPoolID),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(username),
			"PASSWORD":    aws.String(password),
			"SECRET_HASH": aws.String(r.getSecretHash(username)),
		},
	})
	if err != nil {
		return AWSCognitoToken{}, err
	}

	return AWSCognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Second * time.Duration(*res.AuthenticationResult.ExpiresIn)),
		RefreshToken:         *res.AuthenticationResult.RefreshToken,
	}, nil
}

// Logout logs out of Cognito.
func (r *AWSCognitoRepository) Logout(refreshToken string) error {
	_, err := r.c.RevokeToken(&cognitoidentityprovider.RevokeTokenInput{
		ClientId:     aws.String(r.userPoolClientID),
		Token:        aws.String(refreshToken),
		ClientSecret: aws.String(r.userPoolClientSecret),
	})
	return err
}

func (r *AWSCognitoRepository) RefreshToken(refreshToken, username string) (AWSCognitoToken, error) {
	res, err := r.c.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   aws.String("REFRESH_TOKEN_AUTH"),
		ClientId:   aws.String(r.userPoolClientID),
		UserPoolId: aws.String(r.userPoolID),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
			"SECRET_HASH":   aws.String(r.getSecretHash(username)),
		},
	})
	if err != nil {
		return AWSCognitoToken{}, err
	}

	return AWSCognitoToken{
		AccessToken:          *res.AuthenticationResult.AccessToken,
		AccessTokenExpiresAt: time.Now().Add(time.Second * time.Duration(*res.AuthenticationResult.ExpiresIn)),
	}, nil
}

// SetUserPassword sets the password of the user.
func (r *AWSCognitoRepository) SetUserPassword(username, password string, permanent bool) error {
	_, err := r.c.AdminSetUserPassword(&cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  aws.Bool(permanent),
	})

	return err
}

// ChangePassword changes the password of the user.
func (r *AWSCognitoRepository) ChangePassword(authorizationHeader, previousPassword, proposedPassword string) error {
	accessToken, err := r.getAccessToken(authorizationHeader)
	if err != nil {
		return err
	}

	_, err = r.c.ChangePassword(&cognitoidentityprovider.ChangePasswordInput{
		PreviousPassword: aws.String(previousPassword),
		ProposedPassword: aws.String(proposedPassword),
		AccessToken:      aws.String(accessToken),
	})
	return err
}

// ResetPassword resets the specified user's password in a user pool as an administrator. Works on any user.
func (r *AWSCognitoRepository) ResetUserPassword(username string) error {
	_, err := r.c.AdminResetUserPassword(&cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: aws.String(r.userPoolID),
		Username:   aws.String(username),
	})
	return err
}

// ConfirmForgotPassword allows a user to enter a confirmation code to reset a forgotten password.
func (r *AWSCognitoRepository) ConfirmForgotPassword(username, password, confirmationCode string) error {
	_, err := r.c.ConfirmForgotPassword(&cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(r.userPoolClientID),
		Username:         aws.String(username),
		Password:         aws.String(password),
		ConfirmationCode: aws.String(confirmationCode),
	})
	return err
}

// getSecretHash gets the secret hash.
func (r *AWSCognitoRepository) getSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(r.userPoolClientSecret))
	mac.Write([]byte(username + r.userPoolClientID))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// getSecretHash gets the access token.
func (r *AWSCognitoRepository) getAccessToken(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", AWSCognitoAccessTokenNotFound
	}

	numberOfElementsInArray := 2
	prefixBearer := "Bearer "
	stringArr := strings.Split(authorizationHeader, prefixBearer)
	if len(stringArr) == numberOfElementsInArray {
		return stringArr[1], nil
	}
	return "", AWSCognitoAccessTokenFormatNotSupported
}
