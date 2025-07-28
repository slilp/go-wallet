package restapis_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/slilp/go-wallet/internal/consts"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
)

func (suite *RestApisTestSuite) TestTransferBalance() {
	testCases := []struct {
		name        string
		reqBody     api_gen.TransferRequest
		mock        func()
		wantStatus  int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingValidRequest_WhenTransferBalanceSuccess_ThenReturnOk",
			reqBody: api_gen.TransferRequest{
				FromWalletId: "<Wallet1>",
				ToWalletId:   "<Wallet2>",
				Amount:       100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleTransferBalance("<Wallet1>", "<Wallet2>", float64(100)).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "GivingInsufficientBalance_WhenTransferBalance_ThenReturnBadRequest",
			reqBody: api_gen.TransferRequest{
				FromWalletId: "<Wallet1>",
				ToWalletId:   "<Wallet2>",
				Amount:       100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleTransferBalance("<Wallet1>", "<Wallet2>", float64(100)).
					Return(consts.ErrInsufficientBalance)
			},
			wantStatus:  http.StatusBadRequest,
			wantErr:     true,
			expectedErr: "Insufficient balance",
		},
		{
			name: "GivingValidRequest_WhenTransferBalanceFail_ThenReturnInternalServerError",
			reqBody: api_gen.TransferRequest{
				FromWalletId: "<Wallet1>",
				ToWalletId:   "<Wallet2>",
				Amount:       100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleTransferBalance("<Wallet1>", "<Wallet2>", float64(100)).
					Return(errors.New("some error"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantErr:     true,
			expectedErr: "Fail to transfer",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest("POST", "/secure/transfer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			suite.server.ServeHTTP(w, req)
			suite.Equal(tc.wantStatus, w.Code)
			if tc.wantErr {
				var resp api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				suite.Equal(strconv.Itoa(w.Code), resp.ErrorCode)
				suite.Equal(tc.expectedErr, resp.ErrorMessage)
			}
		})
	}
}

func (suite *RestApisTestSuite) TestDepositPoints() {
	testCases := []struct {
		name        string
		reqBody     api_gen.DepositRequest
		mock        func()
		wantStatus  int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingValidRequest_WhenDepositPointsSuccess_ThenReturnOk",
			reqBody: api_gen.DepositRequest{
				WalletId: "<Wallet1>",
				Amount:   100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleDepositWithDrawBalance("<Wallet1>", float64(100)).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "GivingValidRequest_WhenDepositPointsFail_ThenReturnInternalServerError",
			reqBody: api_gen.DepositRequest{
				WalletId: "<Wallet1>",
				Amount:   100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleDepositWithDrawBalance("<Wallet1>", float64(100)).
					Return(errors.New("fail"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantErr:     true,
			expectedErr: "Fail to deposit",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest("POST", "/secure/deposit", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			suite.server.ServeHTTP(w, req)
			suite.Equal(tc.wantStatus, w.Code)
			if tc.wantErr {
				var resp api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				suite.Equal(strconv.Itoa(w.Code), resp.ErrorCode)
				suite.Equal(tc.expectedErr, resp.ErrorMessage)
			}
		})
	}
}

func (suite *RestApisTestSuite) TestWithdrawPoints() {
	testCases := []struct {
		name        string
		reqBody     api_gen.WithdrawRequest
		mock        func()
		wantStatus  int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingValidRequest_WhenWithdrawPointsSuccess_ThenReturnOk",
			reqBody: api_gen.WithdrawRequest{
				WalletId: "<Wallet1>",
				Amount:   100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleDepositWithDrawBalance("<Wallet1>", float64(100)).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "GivingValidRequest_WhenWithdrawPointsFail_ThenReturnInternalServerError",
			reqBody: api_gen.WithdrawRequest{
				WalletId: "<Wallet1>",
				Amount:   100,
			},
			mock: func() {
				suite.mockTransactionService.EXPECT().
					HandleDepositWithDrawBalance("<Wallet1>", float64(100)).
					Return(errors.New("fail"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantErr:     true,
			expectedErr: "Fail to withdraw",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest("POST", "/secure/withdraw", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			suite.server.ServeHTTP(w, req)
			suite.Equal(tc.wantStatus, w.Code)
			if tc.wantErr {
				var resp api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				suite.Equal(strconv.Itoa(w.Code), resp.ErrorCode)
				suite.Equal(tc.expectedErr, resp.ErrorMessage)
			}
		})
	}
}

func (suite *RestApisTestSuite) TestListWalletTransactions() {
	testCases := []struct {
		name           string
		walletId       string
		mock           func()
		wantStatus     int
		wantErr        bool
		expectedErr    string
		expectedTxsLen int
	}{
		{
			name:     "GivingValidRequest_WhenListWalletTransactionsSuccess_ThenReturnOk",
			walletId: "<Wallet1>",
			mock: func() {
				suite.mockListTransactionsService.EXPECT().
					Handle("<UserID>", "<Wallet1>", 1, 30).
					Return(int64(100), []api_gen.TransactionResponseData{
						{
							FromWalletId: "<Wallet1>",
							ToWalletId:   "<Wallet2>",
							Amount:       50,
							Type:         "transfer",
						},
						{
							FromWalletId: "<Wallet2>",
							ToWalletId:   "<Wallet1>",
							Amount:       30,
							Type:         "transfer",
						},
					}, nil)
			},
			wantStatus:     http.StatusOK,
			wantErr:        false,
			expectedTxsLen: 2,
		},
		{
			name:     "GivingValidRequest_WhenListWalletTransactionsFail_ThenReturnInternalServerError",
			walletId: "<Wallet1>",
			mock: func() {
				suite.mockListTransactionsService.EXPECT().
					Handle("<UserID>", "<Wallet1>", 1, 30).
					Return(int64(0), nil, errors.New("fail"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantErr:     true,
			expectedErr: "Failed to list wallets",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/secure/wallet/"+tc.walletId+"/transactions?page=1&limit=30", nil)
			suite.server.ServeHTTP(w, req)
			suite.Equal(tc.wantStatus, w.Code)
			if tc.wantErr {
				var resp api_gen.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				suite.Equal(strconv.Itoa(w.Code), resp.ErrorCode)
				suite.Equal(tc.expectedErr, resp.ErrorMessage)
			} else {
				var resp struct {
					Data *[]api_gen.TransactionResponseData `json:"data"`
				}
				json.Unmarshal(w.Body.Bytes(), &resp)
				suite.NotNil(resp.Data)
				suite.Equal(tc.expectedTxsLen, len(*resp.Data))
			}
		})
	}
}
