package commands_test

import (
	"errors"
	"sync"

	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	mock_repositories "github.com/slilp/go-wallet/internal/repositories/mocks"
	"go.uber.org/mock/gomock"
)

func (suite *CommandsTestSuite) TestRegisterService_Handle() {

	var (
		wg sync.WaitGroup
	)

	testCases := []struct {
		name        string
		mock        func(*mock_repositories.MockUserRepository)
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenValidRequest_WhenCreateSuccess_ThenSuccessIsReturned",
			mock: func(mockUserRepo *mock_repositories.MockUserRepository) {
				mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivienValidRequest_WhenCreateFails_ThenErrorIsReturned",
			mock: func(mockUserRepo *mock_repositories.MockUserRepository) {
				mockUserRepo.EXPECT().Create(gomock.Any()).Return(errors.New("something wrong"))
			},
			wantErr:     true,
			expectedErr: "something wrong",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			tc.mock(suite.mockUserRepo)

			err := suite.registerService.Handle(api_gen.RegisterRequest{
				Email:       "<Email>",
				Password:    "<Password>",
				DisplayName: "<DisplayName>",
			})

			wg.Wait()

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}
		})
	}
}
