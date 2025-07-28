package repositories_test

import (
	"errors"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/slilp/go-wallet/internal/repositories/entity"
)

func (suite *UserRepositoryTestSuite) TestCreate() {
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		input       entity.User
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenNewUser_WhenCreateSuccess_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			input: entity.User{
				Email:       "<Email>",
				Password:    "<Password>",
				DisplayName: "<DisplayName>",
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenNewUser_WhenCreateFail_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).
					WillReturnError(errors.New("insert failed"))
				mock.ExpectRollback()
			},
			input: entity.User{
				Email:       "<Email>",
				Password:    "<Password>",
				DisplayName: "<DisplayName>",
			},
			wantErr:     true,
			expectedErr: "insert failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			err := suite.userRepo.Create(tc.input)

			if tc.wantErr {
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
			}

			suite.sqlMock.ExpectationsWereMet()
		})
	}
}

func (suite *UserRepositoryTestSuite) TestQueryByEmail() {
	var (
		emailArg = "<Email>"
	)
	testCases := []struct {
		name        string
		mock        func(sqlmock.Sqlmock)
		want        *entity.User
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivenEmail_WhenUserFound_ThenSuccess",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "display_name", "created_at", "updated_at"}).
					AddRow("6110d5ec-9abb-4cfd-b62d-92d0e5186f77", emailArg, "<Password>", "<DisplayName>", time.Now(), time.Now())
				suite.sqlMock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 LIMIT \$2`).WithArgs(emailArg, 1).WillReturnRows(rows)
			},
			want: &entity.User{
				ID:          "6110d5ec-9abb-4cfd-b62d-92d0e5186f77",
				Email:       "<Email>",
				Password:    "<Password>",
				DisplayName: "<DisplayName>",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivenEmail_WhenUserNotFound_ThenError",
			mock: func(mock sqlmock.Sqlmock) {
				suite.sqlMock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 LIMIT \$2`).WithArgs(emailArg, 1).WillReturnError(errors.New("something wrong"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "something wrong",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock(suite.sqlMock)

			result, err := suite.userRepo.QueryByEmail(emailArg)

			if tc.wantErr {
				suite.Nil(result)
				suite.EqualError(err, tc.expectedErr)
			} else {
				suite.NoError(err)
				suite.NotNil(result)
				suite.Equal(result.Email, tc.want.Email)
			}
		})
	}
}
