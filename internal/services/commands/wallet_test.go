package commands_test

import (
	"errors"
	"testing"

	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	mock_repositories "github.com/slilp/go-wallet/internal/repositories/mocks"
	"github.com/slilp/go-wallet/internal/services/commands"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWalletService_HandleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_repositories.NewMockWalletRepository(ctrl)
	service := commands.NewWalletService(mockWalletRepo)

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
				mockWalletRepo.EXPECT().
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
				mockWalletRepo.EXPECT().
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
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := service.HandleCreate(tc.userId, tc.req)
			if tc.wantErr {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWalletService_HandleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_repositories.NewMockWalletRepository(ctrl)
	service := commands.NewWalletService(mockWalletRepo)

	testCases := []struct {
		name        string
		walletId    string
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "GivenValidWalletId_WhenDeleteSuccess_ThenSuccess",
			walletId: "<WalletID>",
			mock: func() {
				mockWalletRepo.EXPECT().Delete("<WalletID>").Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:     "GivenValidWalletId_WhenDeleteFails_ThenError",
			walletId: "<WalletID>",
			mock: func() {
				mockWalletRepo.EXPECT().Delete("<WalletID>").Return(errors.New("delete error"))
			},
			wantErr:     true,
			expectedErr: "delete error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := service.HandleDelete(tc.walletId)
			if tc.wantErr {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWalletService_HandleUpdateInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mock_repositories.NewMockWalletRepository(ctrl)
	service := commands.NewWalletService(mockWalletRepo)

	testCases := []struct {
		name        string
		walletId    string
		req         api_gen.WalletRequest
		mock        func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "GivingValidWalletId_WhenUpdateSuccess_ThenSuccess",
			walletId: "<WalletID>",
			req: api_gen.WalletRequest{
				Name:        "<EditWalletName>",
				Description: null.StringFrom("<EditDescription>").Ptr(),
			},
			mock: func() {
				mockWalletRepo.EXPECT().UpdateInfo("<WalletID>", "<EditWalletName>", null.StringFrom("<EditDescription>").Ptr()).Return(nil)
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:     "GivingValidWalletId_WhenUpdateFails_ThenError",
			walletId: "<WalletID>",
			req: api_gen.WalletRequest{
				Name:        "<EditWalletName>",
				Description: null.StringFrom("<EditDescription>").Ptr(),
			},
			mock: func() {
				mockWalletRepo.EXPECT().UpdateInfo("<WalletID>", "<EditWalletName>", null.StringFrom("<EditDescription>").Ptr()).Return(errors.New("update error"))
			},
			wantErr:     true,
			expectedErr: "update error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := service.HandleUpdateInfo(tc.walletId, tc.req)
			if tc.wantErr {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
