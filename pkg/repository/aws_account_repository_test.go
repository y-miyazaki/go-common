package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
)

var errTestAccountService = errors.New("account service error")

func TestNewAWSAccountRepository(t *testing.T) {
	t.Parallel()

	mockClient := &account.Client{}
	repo := NewAWSAccountRepository(mockClient)

	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
}

func TestAWSAccountRepository_GetAccountInformation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		accountID  string
		setupMock  func(m *mocks.MockAWSAccountClientInterface)
		want       *account.GetAccountInformationOutput
		wantErr    error
		wantErrMsg string
	}{
		{
			name:      "with account ID",
			accountID: "123456789012",
			setupMock: func(m *mocks.MockAWSAccountClientInterface) {
				expected := &account.GetAccountInformationOutput{
					AccountId:   aws.String("123456789012"),
					AccountName: aws.String("sample-account"),
				}
				m.EXPECT().
					GetAccountInformation(gomock.Any(), gomock.AssignableToTypeOf(&account.GetAccountInformationInput{})).
					DoAndReturn(func(_ context.Context, input *account.GetAccountInformationInput, _ ...func(*account.Options)) (*account.GetAccountInformationOutput, error) {
						require.NotNil(t, input.AccountId)
						require.Equal(t, "123456789012", *input.AccountId)
						return expected, nil
					})
			},
			want: &account.GetAccountInformationOutput{
				AccountId:   aws.String("123456789012"),
				AccountName: aws.String("sample-account"),
			},
		},
		{
			name:      "without account ID",
			accountID: "",
			setupMock: func(m *mocks.MockAWSAccountClientInterface) {
				expected := &account.GetAccountInformationOutput{}
				m.EXPECT().
					GetAccountInformation(gomock.Any(), gomock.AssignableToTypeOf(&account.GetAccountInformationInput{})).
					DoAndReturn(func(_ context.Context, input *account.GetAccountInformationInput, _ ...func(*account.Options)) (*account.GetAccountInformationOutput, error) {
						require.Nil(t, input.AccountId)
						return expected, nil
					})
			},
			want: &account.GetAccountInformationOutput{},
		},
		{
			name:      "client error",
			accountID: "123456789012",
			setupMock: func(m *mocks.MockAWSAccountClientInterface) {
				m.EXPECT().
					GetAccountInformation(gomock.Any(), gomock.Any()).
					Return(nil, errTestAccountService)
			},
			wantErrMsg: "account GetAccountInformation",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockAWSAccountClientInterface(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			repo := NewAWSAccountRepositoryWithInterface(mockClient)
			got, err := repo.GetAccountInformation(context.Background(), tc.accountID)

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
