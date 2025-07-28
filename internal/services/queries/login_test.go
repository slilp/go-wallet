package queries_test

import (
	"errors"

	"github.com/slilp/go-wallet/internal/repositories/entity"
	mock_repositories "github.com/slilp/go-wallet/internal/repositories/mocks"
	"golang.org/x/crypto/bcrypt"
)

func (suite *QueriesTestSuite) TestLoginService_Handle() {

	testCases := []struct {
		name        string
		mock        func(*mock_repositories.MockUserRepository)
		want        *entity.User
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingCorrectEmailPassword_WhenMatch_ThenSuccess",
			mock: func(mockUserRepo *mock_repositories.MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("<Password>"), bcrypt.DefaultCost)

				mockUserRepo.EXPECT().QueryByEmail("<Email>").Return(&entity.User{
					Email:       "<Email>",
					Password:    string(hashedPassword),
					DisplayName: "<DisplayName>",
				}, nil)

			},
			want: &entity.User{
				Email:       "<Email>",
				Password:    "<Password>",
				DisplayName: "<DisplayName>",
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivingIncorrectEmail_WhenNotMatch_ThenError",
			mock: func(mockUserRepo *mock_repositories.MockUserRepository) {
				mockUserRepo.EXPECT().QueryByEmail("<Email>").Return(nil, errors.New("something wrong"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "something wrong",
		},
		{
			name: "GivingIncorrectEmailPassword_WhenNotMatch_ThenError",
			mock: func(mockUserRepo *mock_repositories.MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("<WrongPassword>"), bcrypt.DefaultCost)

				mockUserRepo.EXPECT().QueryByEmail("<Email>").Return(&entity.User{
					Email:       "<Email>",
					Password:    string(hashedPassword),
					DisplayName: "<DisplayName>",
				}, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.mockUserRepo)

			result, err := suite.loginService.Handle("<Email>", "<Password>")

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
				suite.Nil(result)
			} else {
				suite.NoError(err)
				suite.NotNil(result)
				resultValue := *result
				suite.Equal("<Email>", resultValue.Email)
				suite.Equal("<DisplayName>", resultValue.DisplayName)
			}
		})
	}
}
