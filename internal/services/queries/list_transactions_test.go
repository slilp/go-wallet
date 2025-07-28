package queries_test

import (
	"errors"
	"time"

	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	"gorm.io/gorm"
)

func (suite *QueriesTestSuite) TestListTransactionsService_Handle() {
	testCases := []struct {
		name        string
		userId      string
		walletId    string
		page        int
		limit       int
		setupMocks  func()
		want        []api_gen.TransactionResponseData
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "GivenValidRequest_WhenSuccess_ThenReturnTransactions",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			page:     1,
			limit:    10,
			setupMocks: func() {
				wallet := &entity.Wallet{
					ID:     "<WalletID>",
					UserID: "<UserID>",
					Name:   "<WalletName>",
				}
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(wallet, nil)

				suite.mockTransactionRepo.EXPECT().
					CountByWalletId("<WalletID>").
					Return(int64(1), nil)

				transactions := []entity.Transaction{
					{
						ID:        "<TransactionID>",
						From:      null.StringFrom("<FromWalletID>").Ptr(),
						To:        null.StringFrom("<ToWalletID>").Ptr(),
						Amount:    100,
						Type:      "transfer",
						CreatedAt: time.Now(),
					},
				}
				suite.mockTransactionRepo.EXPECT().
					List("<WalletID>", 1, 10).
					Return(transactions, nil)
			},
			want: []api_gen.TransactionResponseData{
				{
					Id:           "<TransactionID>",
					FromWalletId: "<FromWalletID>",
					ToWalletId:   "<ToWalletID>",
					Amount:       100,
					Type:         "transfer",
				},
			},
			wantErr: false,
		},
		{
			name:     "GivenValidRequest_WhenWalletNotFound_ThenReturnError",
			userId:   "<UserID>",
			walletId: "<NonExistentWalletID>",
			page:     1,
			limit:    10,
			setupMocks: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<NonExistentWalletID>").
					Return(nil, gorm.ErrRecordNotFound)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name:     "GivenValidRequest_WhenWalletRepoError_ThenReturnError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			page:     1,
			limit:    10,
			setupMocks: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, errors.New("database error"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "database error",
		},
		{
			name:     "GivenValidRequest_WhenTransactionRepoError_ThenReturnError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			page:     1,
			limit:    10,
			setupMocks: func() {
				wallet := &entity.Wallet{
					ID:     "<WalletID>",
					UserID: "<UserID>",
					Name:   "<WalletName>",
				}
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(wallet, nil)

				suite.mockTransactionRepo.EXPECT().
					CountByWalletId("<WalletID>").
					Return(int64(1), nil)

				suite.mockTransactionRepo.EXPECT().
					List("<WalletID>", 1, 10).
					Return(nil, errors.New("transaction query failed"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "transaction query failed",
		},
		{
			name:     "GivenValidRequest_WhenNoTransactions_ThenReturnEmptyArray",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			page:     1,
			limit:    10,
			setupMocks: func() {
				wallet := &entity.Wallet{
					ID:     "<WalletID>",
					UserID: "<UserID>",
					Name:   "<WalletName>",
				}
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(wallet, nil)

				suite.mockTransactionRepo.EXPECT().
					CountByWalletId("<WalletID>").
					Return(int64(0), nil)
			},
			want:    []api_gen.TransactionResponseData{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setupMocks()

			total, result, err := suite.listTransactionsService.Handle(tc.userId, tc.walletId, tc.page, tc.limit)

			if tc.wantErr {
				suite.Error(err)
				suite.Contains(err.Error(), tc.expectedErr)
				suite.Equal(int64(0), total)
				suite.Nil(result)
			} else {
				suite.NoError(err)
				suite.Equal(int64(len(tc.want)), total)
				suite.Equal(len(tc.want), len(result))
				if len(result) > 0 {
					suite.Equal(tc.want[0].Id, result[0].Id)
					suite.Equal(tc.want[0].FromWalletId, result[0].FromWalletId)
					suite.Equal(tc.want[0].ToWalletId, result[0].ToWalletId)
					suite.Equal(tc.want[0].Amount, result[0].Amount)
					suite.Equal(tc.want[0].Type, result[0].Type)
				}
			}
		})
	}
}
