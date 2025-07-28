package commands_test

import (
	"errors"
)

func (suite *CommandsTestSuite) TestTransactionService_HandleTransferBalance() {
	testCases := []struct {
		name        string
		from        string
		to          string
		amount      float64
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:   "GivingValidFromToAmount_WhenUpdateBalanceSuccess_ThenSuccess",
			from:   "<FromWalletID>",
			to:     "<ToWalletID>",
			amount: 100.0,
			mock: func() {
				suite.mockTransactionRepo.EXPECT().UpdateTransferTransaction("<FromWalletID>", "<ToWalletID>", 100.0).Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:   "GivingValidFromToAmount_WhenUpdateBalanceFails_ThenError",
			from:   "<FromWalletID>",
			to:     "<ToWalletID>",
			amount: 100.0,
			mock: func() {
				suite.mockTransactionRepo.EXPECT().UpdateTransferTransaction("<FromWalletID>", "<ToWalletID>", 100.0).Return(errors.New("update balance error"))
			},
			wantErr:     true,
			expectedErr: "update balance error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			err := suite.transactionService.HandleTransferBalance(tc.from, tc.to, tc.amount)
			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *CommandsTestSuite) TestWalletService_HandleDepositWithDrawBalance() {

	testCases := []struct {
		name        string
		walletId    string
		amount      float64
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "GivenValidWalletIdAndPositiveAmount_WhenDepositSuccess_ThenSuccess",
			walletId: "<WalletID>",
			amount:   100.0,
			mock: func() {
				suite.mockTransactionRepo.EXPECT().UpdateBalanceTransaction("<WalletID>", 100.0).Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:     "GivenValidWalletIdAndNegativeAmount_WhenWithdrawSuccess_ThenSuccess",
			walletId: "<WalletID>",
			amount:   -50.0,
			mock: func() {
				suite.mockTransactionRepo.EXPECT().UpdateBalanceTransaction("<WalletID>", -50.0).Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:     "GivenValidWalletId_WhenUpdateBalanceFails_ThenError",
			walletId: "<WalletID>",
			amount:   10.0,
			mock: func() {
				suite.mockTransactionRepo.EXPECT().UpdateBalanceTransaction("<WalletID>", 10.0).Return(errors.New("update balance error"))
			},
			wantErr:     true,
			expectedErr: "update balance error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			err := suite.transactionService.HandleDepositWithDrawBalance(tc.walletId, tc.amount)
			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}
		})
	}
}
