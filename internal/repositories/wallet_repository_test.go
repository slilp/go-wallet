package repositories_test

import (
	"errors"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/repositories/entity"
)

func (suite *WalletRepositoryTestSuite) TestCreate() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		input       entity.Wallet
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenNewWallet_WhenCreateSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "wallets"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			input: entity.Wallet{
				ID:        "<ID>",
				UserID:    "<UserID>",
				Name:      "<Name>",
				Balance:   100,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenNewWallet_WhenCreateFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "wallets"`).
					WillReturnError(errors.New("insert failed"))
				mock.ExpectRollback()
			},
			input: entity.Wallet{
				ID:        "<ID>",
				UserID:    "<UserID>",
				Name:      "<Name>",
				Balance:   100,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr:     true,
			expectedErr: "insert failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			err := suite.walletRepo.Create(tc.input)

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}

			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *WalletRepositoryTestSuite) TestDelete() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		walletId    string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenWalletId_WhenDeleteSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "wallets"`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			walletId:    "<ID>",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenWalletId_WhenDeleteFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "wallets"`).WillReturnError(errors.New("delete failed"))
				mock.ExpectRollback()
			},
			walletId:    "<ID>",
			wantErr:     true,
			expectedErr: "delete failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			err := suite.walletRepo.Delete(tc.walletId)

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}

			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *WalletRepositoryTestSuite) TestListAll() {
	var (
		userIdArg = "<UserID>"
	)
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		want        []entity.Wallet
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenUserId_WhenWalletsFound_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "balance", "created_at", "updated_at"}).
					AddRow("1", userIdArg, "<Name>", 100, time.Now(), time.Now())
				suite.sqlMock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."user_id" = \$1`).WithArgs(userIdArg).WillReturnRows(rows)
			},
			want: []entity.Wallet{
				{
					ID:      "1",
					UserID:  userIdArg,
					Name:    "<Name>",
					Balance: 100,
				},
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenUserId_WhenQueryFails_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				suite.sqlMock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."user_id" = \$1`).WithArgs(userIdArg).WillReturnError(errors.New("query failed"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "query failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			result, err := suite.walletRepo.ListAll(userIdArg)

			if tc.wantErr {
				suite.Nil(result)
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
				suite.NotNil(result)
				suite.Equal(result[0].UserID, tc.want[0].UserID)
			}
		})
	}
}

func (suite *WalletRepositoryTestSuite) TestUpdateInfo() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		id          string
		nameArg     string
		descArg     *string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenWallet_WhenUpdateInfoSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "wallets"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			id:          "<ID>",
			nameArg:     "<Name>",
			descArg:     null.StringFrom("<Description>").Ptr(),
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenWallet_WhenUpdateInfoFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "wallets"`).
					WillReturnError(errors.New("update info failed"))
				mock.ExpectRollback()
			},
			id:          "<ID>",
			nameArg:     "<Name>",
			descArg:     null.StringFrom("<Description>").Ptr(),
			wantErr:     true,
			expectedErr: "update info failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			err := suite.walletRepo.UpdateInfo(tc.id, tc.nameArg, tc.descArg)

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}

			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *WalletRepositoryTestSuite) TestQueryByIdAndUser() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		userId      string
		walletId    string
		want        *entity.Wallet
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenWalletId_WhenFound_ThenReturnWallet",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "balance", "created_at", "updated_at"}).
					AddRow("<ID>", "<UserID>", "<Name>", 100, nil, nil)
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 AND "wallets"\."user_id" = \$2 ORDER BY "wallets"\."id" LIMIT \$3`).
					WithArgs("<ID>", "<UserID>", 1).
					WillReturnRows(rows)
			},
			userId:   "<UserID>",
			walletId: "<ID>",
			want: &entity.Wallet{
				ID:      "<ID>",
				UserID:  "<UserID>",
				Name:    "<Name>",
				Balance: 100,
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenWalletId_WhenNotFound_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "balance", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 AND "wallets"\."user_id" = \$2 ORDER BY "wallets"\."id" LIMIT \$3`).
					WithArgs("<ID>", "<UserID>", 1).
					WillReturnRows(rows)
			},
			userId:      "<UserID>",
			walletId:    "<ID>",
			want:        nil,
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name: "GivenWalletId_WhenQueryFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE "wallets"\."id" = \$1 AND "wallets"\."user_id" = \$2 ORDER BY "wallets"\."id" LIMIT \$3`).
					WithArgs("<ID>", "<UserID>", 1).
					WillReturnError(errors.New("query failed"))
			},
			userId:      "<UserID>",
			walletId:    "<ID>",
			want:        nil,
			wantErr:     true,
			expectedErr: "query failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			result, err := suite.walletRepo.QueryByIdAndUser(tc.userId, tc.walletId)

			if tc.wantErr {
				suite.Nil(result)
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
				suite.NotNil(result)
				suite.Equal(tc.want.ID, result.ID)
				suite.Equal(tc.want.UserID, result.UserID)
				suite.Equal(tc.want.Name, result.Name)
				suite.Equal(tc.want.Balance, result.Balance)
			}

			suite.sqlMock.ExpectationsWereMet()
		})
	}
}
