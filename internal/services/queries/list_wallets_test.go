package queries_test

import (
	"errors"
	"time"

	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	mock_repositories "github.com/slilp/go-wallet/internal/repositories/mocks"
)

func (suite *QueriesTestSuite) TestListWalletsService_Handle() {
	testCases := []struct {
		name        string
		mock        func(*mock_repositories.MockWalletRepository)
		want        []api_gen.WalletResponseData
		wantErr     bool
		expectedErr string
		userId      string
	}{
		{
			name: "GivenValidUserId_WhenWalletsExist_ThenReturnWallets",
			mock: func(mockWalletRepo *mock_repositories.MockWalletRepository) {
				wallets := []entity.Wallet{
					{
						ID:          "<WalletID>",
						Balance:     1000,
						Name:        "<WalletName>",
						Description: null.StringFrom("<WalletDescription>").Ptr(),
						UpdatedAt:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				mockWalletRepo.EXPECT().ListAll("user1").Return(wallets, nil)
			},
			want: []api_gen.WalletResponseData{
				{
					Id:          "<WalletID>",
					Balance:     1000,
					Name:        "<WalletName>",
					Description: null.StringFrom("<WalletDescription>").Ptr(),
					UpdatedAt:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr:     false,
			expectedErr: "",
			userId:      "user1",
		},
		{
			name: "GivenValidUserId_WhenNoWallets_ThenReturnEmptyList",
			mock: func(mockWalletRepo *mock_repositories.MockWalletRepository) {
				mockWalletRepo.EXPECT().ListAll("<UserID>").Return([]entity.Wallet{}, nil)
			},
			want:        []api_gen.WalletResponseData{},
			wantErr:     false,
			expectedErr: "",
			userId:      "<UserID>",
		},
		{
			name: "GivenUserId_WhenRepoReturnsError_ThenReturnError",
			mock: func(mockWalletRepo *mock_repositories.MockWalletRepository) {
				mockWalletRepo.EXPECT().ListAll("<UserID>").Return(nil, errors.New("repo error"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "repo error",
			userId:      "<UserID>",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.mockWalletRepo)

			result, err := suite.listWalletsService.Handle(tc.userId)

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
				suite.Nil(result)
			} else {
				suite.NoError(err)
				suite.Equal(tc.want, result)
			}
		})
	}
}
