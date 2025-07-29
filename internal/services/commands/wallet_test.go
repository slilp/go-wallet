package commands_test

import (
	"errors"

	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (suite *CommandsTestSuite) TestWalletService_HandleCreate() {

	testCases := []struct {
		name        string
		userId      string
		req         api_gen.WalletRequest
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:   "GivenValidRequest_WhenCreateSuccess_ThenSucces",
			userId: "<UserID>",
			req: api_gen.WalletRequest{
				Name:        "<WalletName>",
				Description: nil,
			},
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					Create(entity.Wallet{
						UserID:      "<UserID>",
						Name:        "<WalletName>",
						Description: nil,
					}).
					Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:   "GivenValidRequest_WhenCreateFails_ThenError",
			userId: "<UserID>",
			req: api_gen.WalletRequest{
				Name:        "<WalletName>",
				Description: nil,
			},
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					Create(entity.Wallet{
						UserID:      "<UserID>",
						Name:        "<WalletName>",
						Description: nil,
					}).
					Return(errors.New("create error"))
			},
			wantErr:     true,
			expectedErr: "create error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			err := suite.walletService.HandleCreate(tc.userId, tc.req)
			if tc.wantErr {
				assert.EqualError(suite.T(), err, tc.expectedErr)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *CommandsTestSuite) TestWalletService_HandleDelete() {

	testCases := []struct {
		name        string
		userId      string
		walletId    string
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "GivenValidWalletId_WhenDeleteSuccess_ThenSuccess",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, nil)
				suite.mockWalletRepo.EXPECT().Delete("<WalletID>").Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:     "GivenInvalidWalletId_WhenNotFound_ThenError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name:     "GivenValidWalletId_WhenDeleteFail_ThenError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, nil)
				suite.mockWalletRepo.EXPECT().Delete("<WalletID>").Return(errors.New("delete error"))
			},
			wantErr:     true,
			expectedErr: "delete error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			err := suite.walletService.HandleDelete(tc.userId, tc.walletId)
			if tc.wantErr {
				assert.EqualError(suite.T(), err, tc.expectedErr)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *CommandsTestSuite) TestWalletService_HandleUpdateInfo() {

	testCases := []struct {
		name        string
		userId      string
		walletId    string
		req         api_gen.WalletRequest
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "GivingValidWalletId_WhenUpdateSuccess_ThenSuccess",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			req: api_gen.WalletRequest{
				Name:        "<EditWalletName>",
				Description: null.StringFrom("<EditDescription>").Ptr(),
			},
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, nil)
				suite.mockWalletRepo.EXPECT().UpdateInfo("<WalletID>", "<EditWalletName>", null.StringFrom("<EditDescription>").Ptr()).Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:     "GivingInvalidWalletId_WhenNotFound_ThenError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			req: api_gen.WalletRequest{
				Name:        "<EditWalletName>",
				Description: null.StringFrom("<EditDescription>").Ptr(),
			},
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name:     "GivingValidWalletId_WhenUpdateFail_ThenError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			req: api_gen.WalletRequest{
				Name:        "<EditWalletName>",
				Description: null.StringFrom("<EditDescription>").Ptr(),
			},
			mock: func() {
				suite.mockWalletRepo.EXPECT().
					QueryByIdAndUser("<UserID>", "<WalletID>").
					Return(nil, nil)
				suite.mockWalletRepo.EXPECT().UpdateInfo("<WalletID>", "<EditWalletName>", null.StringFrom("<EditDescription>").Ptr()).Return(errors.New("update error"))
			},
			wantErr:     true,
			expectedErr: "update error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			err := suite.walletService.HandleUpdateInfo(tc.userId, tc.walletId, tc.req)
			if tc.wantErr {
				assert.EqualError(suite.T(), err, tc.expectedErr)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}
