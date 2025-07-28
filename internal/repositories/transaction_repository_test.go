package repositories_test

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
)

func (suite *TransactionRepositoryTestSuite) TestUpdateBalanceTransaction() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		walletId    string
		amount      float64
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenPositiveAmount_WhenUpdateBalanceSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<WalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<WalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<WalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "from_wallet_id", "to_wallet_id", "amount", "type", "created_at"}).
						AddRow("<TransactionID>", nil, "<WalletID>", 100.0, "deposit", nil))
				mock.ExpectCommit()
			},
			walletId:    "<WalletID>",
			amount:      100.0,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenNegativeAmount_WhenUpdateBalanceSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<WalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<WalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<WalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "from_wallet_id", "to_wallet_id", "amount", "type", "created_at"}).
						AddRow("<TransactionID>", "<WalletID>", nil, -50.0, "withdraw", nil))
				mock.ExpectCommit()
			},
			walletId:    "<WalletID>",
			amount:      -50.0,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenAmount_WhenUpdateBalanceFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<WalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<WalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<WalletID>").
					WillReturnError(errors.New("update balance failed"))
				mock.ExpectRollback()
			},
			walletId:    "<WalletID>",
			amount:      10.0,
			wantErr:     true,
			expectedErr: "update balance failed",
		},
		{
			name: "GivenAmount_WhenCreateTransactionFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<WalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<WalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<WalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnError(errors.New("create transaction failed"))
				mock.ExpectRollback()
			},
			walletId:    "<WalletID>",
			amount:      10.0,
			wantErr:     true,
			expectedErr: "create transaction failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)
			err := suite.transactionRepo.UpdateBalanceTransaction(tc.walletId, tc.amount)
			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}
			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *TransactionRepositoryTestSuite) TestUpdateTransferTransaction() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		from        string
		to          string
		amount      float64
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenWallets_WhenUpdateTransferSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<FromWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<FromWalletID>", 100.0))
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<ToWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<ToWalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<FromWalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<ToWalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "from_wallet_id", "to_wallet_id", "amount", "type", "created_at"}).
						AddRow("<TransactionID>", "<FromWalletID>", "<ToWalletID>", 50.0, "transfer", nil))
				mock.ExpectCommit()
			},
			from:        "<FromWalletID>",
			to:          "<ToWalletID>",
			amount:      50.0,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenWallets_WhenInsufficientBalance_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<FromWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<FromWalletID>", 10.0))
				mock.ExpectRollback()
			},
			from:        "<FromWalletID>",
			to:          "<ToWalletID>",
			amount:      50.0,
			wantErr:     true,
			expectedErr: "insufficient balance",
		},
		{
			name: "GivenWallets_WhenUpdateTransferFromFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<FromWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<FromWalletID>", 100.0))
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<ToWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<ToWalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<FromWalletID>").
					WillReturnError(errors.New("from update failed"))
				mock.ExpectRollback()
			},
			from:        "<FromWalletID>",
			to:          "<ToWalletID>",
			amount:      50.0,
			wantErr:     true,
			expectedErr: "from update failed",
		},
		{
			name: "GivenWallets_WhenUpdateTransferToFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<FromWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<FromWalletID>", 100.0))
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<ToWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<ToWalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<FromWalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<ToWalletID>").
					WillReturnError(errors.New("to update failed"))
				mock.ExpectRollback()
			},
			from:        "<FromWalletID>",
			to:          "<ToWalletID>",
			amount:      50.0,
			wantErr:     true,
			expectedErr: "to update failed",
		},
		{
			name: "GivenWallets_WhenCreateTransactionFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<FromWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<FromWalletID>", 100.0))
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 ORDER BY "wallets"\."id" LIMIT \$2 FOR UPDATE`).
					WithArgs("<ToWalletID>", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow("<ToWalletID>", 100.0))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<FromWalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "wallets"`).
					WithArgs(sqlmock.AnyArg(), "<ToWalletID>").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnError(errors.New("create transaction failed"))
				mock.ExpectRollback()
			},
			from:        "<FromWalletID>",
			to:          "<ToWalletID>",
			amount:      50.0,
			wantErr:     true,
			expectedErr: "create transaction failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)
			err := suite.transactionRepo.UpdateTransferTransaction(tc.from, tc.to, tc.amount)
			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}
			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *TransactionRepositoryTestSuite) TestList() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		walletId    string
		page        int
		limit       int
		wantErr     bool
		expectedLen int
	}{
		{
			name: "GivenWalletId_WhenListSuccess_ThenReturnTransactions",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "from", "to", "amount", "type", "created_at"}).
					AddRow("<TransactionID1>", "<WalletID1>", "<ToWalletID1>", 10.0, "transfer", nil).
					AddRow("<TransactionID2>", "<FromWalletID2>", "<WalletID1>", 20.0, "transfer", nil)
				mock.ExpectQuery(`SELECT \* FROM "transactions"`).
					WithArgs("<WalletID1>", "<WalletID1>", 2).
					WillReturnRows(rows)
			},
			walletId:    "<WalletID1>",
			page:        1,
			limit:       2,
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name: "GivenWalletId_WhenListFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "transactions"`).
					WithArgs("<WalletID2>", "<WalletID2>", 2).
					WillReturnError(errors.New("query error"))
			},
			walletId:    "<WalletID2>",
			page:        1,
			limit:       2,
			wantErr:     true,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)
			transactions, err := suite.transactionRepo.List(tc.walletId, tc.page, tc.limit)
			if tc.wantErr {
				suite.Error(err)
				suite.Len(transactions, 0)
			} else {
				suite.NoError(err)
				suite.Len(transactions, tc.expectedLen)
			}
			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *TransactionRepositoryTestSuite) TestCountByWalletId() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		walletId    string
		wantCount   int64
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenWalletId_WhenCountSuccess_ThenReturnCount",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "transactions" WHERE from = \$1 OR to = \$2`).
					WithArgs("<WalletID>", "<WalletID>").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
			},
			walletId:  "<WalletID>",
			wantCount: int64(5),
			wantErr:   false,
		},
		{
			name: "GivenWalletId_WhenCountFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "transactions" WHERE from = \$1 OR to = \$2`).
					WithArgs("<WalletID>", "<WalletID>").
					WillReturnError(errors.New("count error"))
			},
			walletId:    "<WalletID>",
			wantCount:   int64(0),
			wantErr:     true,
			expectedErr: "count error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)
			count, err := suite.transactionRepo.CountByWalletId(tc.walletId)
			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
				suite.Equal(int64(0), count)
			} else {
				suite.NoError(err)
				suite.Equal(tc.wantCount, count)
			}
			suite.sqlMock.ExpectationsWereMet()
		})
	}
}
