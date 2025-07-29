package restapis_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func (suite *RestApisTestSuite) TestCreateWallet() {
	testCases := []struct {
		name        string
		reqBody     interface{}
		mock        func()
		wantStatus  int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingValidRequest_WhenCreateWalletSuccess_ThenReturnCreated",
			reqBody: api_gen.WalletRequest{
				Name:        "Test Wallet",
				Description: null.StringFrom("Test Description").Ptr(),
			},
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleCreate("<UserID>", gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:        "GivingInvalidRequest_WhenCreateWallet_ThenReturnBadRequest",
			reqBody:     api_gen.WalletRequest{},
			mock:        func() {},
			wantStatus:  http.StatusBadRequest,
			wantErr:     true,
			expectedErr: "Name required",
		},
		{
			name: "GivingValidRequest_WhenCreateWalletFail_ThenReturnInternalServerError",
			reqBody: api_gen.WalletRequest{
				Name:        "Test Wallet",
				Description: null.StringFrom("Test Description").Ptr(),
			},
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleCreate("<UserID>", gomock.Any()).
					Return(fmt.Errorf("service error"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantErr:     true,
			expectedErr: "Failed to create wallet",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()

			w := httptest.NewRecorder()
			reqBodyBytes, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest("POST", "/secure/wallet", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			suite.server.ServeHTTP(w, req)

			suite.Equal(tc.wantStatus, w.Code)

			if tc.wantErr {
				var errorResponse api_gen.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				suite.NoError(err)
				suite.Equal(strconv.Itoa(w.Code), errorResponse.ErrorCode)
				suite.Equal(tc.expectedErr, errorResponse.ErrorMessage)
			}
		})
	}
}

func (suite *RestApisTestSuite) TestListUserWallets() {
	tests := []struct {
		name           string
		userId         string
		mock           func()
		expectedStatus int
		expectedError  *api_gen.ErrorResponse
		expectedData   []api_gen.WalletResponseData
	}{
		{
			name:   "GivingValidUserId_WhenListWalletsSuccess_ThenReturnOk",
			userId: "<UserID>",
			mock: func() {
				expectedWallets := []api_gen.WalletResponseData{
					{Id: "wallet1", Name: "Wallet 1", Description: null.StringFrom("Description 1").Ptr(), Balance: 100.50},
					{Id: "wallet2", Name: "Wallet 2", Description: null.StringFrom("Description 2").Ptr(), Balance: 200.75},
				}
				suite.mockListWalletsService.EXPECT().
					Handle("<UserID>").
					Return(expectedWallets, nil)
			},
			expectedStatus: http.StatusOK,
			expectedData: []api_gen.WalletResponseData{
				{Id: "wallet1", Name: "Wallet 1", Description: null.StringFrom("Description 1").Ptr(), Balance: 100.50},
				{Id: "wallet2", Name: "Wallet 2", Description: null.StringFrom("Description 2").Ptr(), Balance: 200.75},
			},
		},
		{
			name:   "GivingValidUserId_WhenListWalletsFail_ThenReturnInternalServerError",
			userId: "<UserID>",
			mock: func() {
				suite.mockListWalletsService.EXPECT().
					Handle("<UserID>").
					Return(nil, fmt.Errorf("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &api_gen.ErrorResponse{
				ErrorCode:    "500",
				ErrorMessage: "Failed to list wallets",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mock()

			w := httptest.NewRecorder()
			httpReq, _ := http.NewRequest("GET", "/secure/wallets", nil)

			suite.server.ServeHTTP(w, httpReq)

			suite.Equal(tt.expectedStatus, w.Code)

			if tt.expectedError != nil {
				var response api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				suite.Equal(tt.expectedError.ErrorCode, response.ErrorCode)
				suite.Equal(tt.expectedError.ErrorMessage, response.ErrorMessage)
			} else if tt.expectedData != nil {
				var response api_gen.ListUserWalletsResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				suite.NotNil(response.Data)
				suite.Equal(len(tt.expectedData), len(*response.Data))
			}
		})
	}
}

func (suite *RestApisTestSuite) TestUpdateWallet() {
	tests := []struct {
		name           string
		userId         string
		walletId       string
		requestBody    interface{}
		mock           func()
		expectedStatus int
		expectedError  *api_gen.ErrorResponse
	}{
		{
			name:     "GivingValidRequest_WhenUpdateWalletSuccess_ThenReturnOk",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			requestBody: api_gen.WalletRequest{
				Name:        "Updated Wallet",
				Description: null.StringFrom("Updated Description").Ptr(),
			},
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleUpdateInfo("<UserID>", "<WalletID>", gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "GivingInvalidRequest_WhenUpdateWallet_ThenReturnBadRequest",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mock:           func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError: &api_gen.ErrorResponse{
				ErrorCode:    "400",
				ErrorMessage: "Name required",
			},
		},
		{
			name:     "GivingValidRequest_WhenWalletNotFound_ThenReturnNotFound",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			requestBody: api_gen.WalletRequest{
				Name:        "Updated Wallet",
				Description: null.StringFrom("Description").Ptr(),
			},
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleUpdateInfo("<UserID>", "<WalletID>", gomock.Any()).
					Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError: &api_gen.ErrorResponse{
				ErrorCode:    "404",
				ErrorMessage: "Wallet not found",
			},
		},
		{
			name:     "GivingValidRequest_WhenUpdateWalletFail_ThenReturnInternalServerError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			requestBody: api_gen.WalletRequest{
				Name:        "Updated Wallet",
				Description: null.StringFrom("Description").Ptr(),
			},
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleUpdateInfo("<UserID>", "<WalletID>", gomock.Any()).
					Return(fmt.Errorf("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &api_gen.ErrorResponse{
				ErrorCode:    "500",
				ErrorMessage: "Failed to update wallet",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mock()

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			httpReq, _ := http.NewRequest("PUT", "/secure/wallet/"+tt.walletId, bytes.NewBuffer(body))
			httpReq.Header.Set("Content-Type", "application/json")

			suite.server.ServeHTTP(w, httpReq)

			suite.Equal(tt.expectedStatus, w.Code)

			if tt.expectedError != nil {
				var response api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				suite.Equal(tt.expectedError.ErrorCode, response.ErrorCode)
				suite.Equal(tt.expectedError.ErrorMessage, response.ErrorMessage)
			}
		})
	}
}

func (suite *RestApisTestSuite) TestDeleteWallet() {
	tests := []struct {
		name           string
		userId         string
		walletId       string
		mock           func()
		expectedStatus int
		expectedError  *api_gen.ErrorResponse
	}{
		{
			name:     "GivingValidWalletId_WhenDeleteWalletSuccess_ThenReturnNoContent",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleDelete("<UserID>", "<WalletID>").
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "GivingValidRequest_WhenWalletNotFound_ThenReturnNotFound",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleDelete("<UserID>", "<WalletID>").
					Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError: &api_gen.ErrorResponse{
				ErrorCode:    "404",
				ErrorMessage: "Wallet not found",
			},
		},
		{
			name:     "GivingValidWalletId_WhenDeleteWalletFail_ThenReturnInternalServerError",
			userId:   "<UserID>",
			walletId: "<WalletID>",
			mock: func() {
				suite.mockWalletService.EXPECT().
					HandleDelete("<UserID>", "<WalletID>").
					Return(fmt.Errorf("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &api_gen.ErrorResponse{
				ErrorCode:    "500",
				ErrorMessage: "Failed to delete wallet",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mock()

			w := httptest.NewRecorder()
			httpReq, _ := http.NewRequest("DELETE", "/secure/wallet/"+tt.walletId, nil)

			suite.server.ServeHTTP(w, httpReq)

			suite.Equal(tt.expectedStatus, w.Code)

			if tt.expectedError != nil {
				var response api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				suite.Equal(tt.expectedError.ErrorCode, response.ErrorCode)
				suite.Equal(tt.expectedError.ErrorMessage, response.ErrorMessage)
			}
		})
	}
}
