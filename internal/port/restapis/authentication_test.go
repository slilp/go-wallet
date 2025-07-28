package restapis_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"go.uber.org/mock/gomock"
)

func (suite *RestApisTestSuite) TestRegisterUser() {
	testCases := []struct {
		name        string
		reqBody     api_gen.RegisterRequest
		mock        func()
		wantStatus  int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingCorrectRequest_WhenRegisterSuccess_ThenReturnCreated",
			reqBody: api_gen.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password",
				DisplayName: "Test User",
			},
			mock: func() {
				suite.mockRegisterService.EXPECT().Handle(gomock.Any()).Return(nil)
			},
			wantStatus:  http.StatusCreated,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivingCorrectRequest_WhenRegisterFail_ThenReturnInternalServerError",
			reqBody: api_gen.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password",
				DisplayName: "Test User",
			},
			mock: func() {
				suite.mockRegisterService.EXPECT().Handle(gomock.Any()).Return(errors.New("something wrong"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantErr:     true,
			expectedErr: "Failed to register user",
		},
		{
			name:        "GivingIncorrectRequest_WhenRegister_ThenReturnBadRequest",
			reqBody:     api_gen.RegisterRequest{},
			mock:        func() {},
			wantStatus:  http.StatusBadRequest,
			wantErr:     true,
			expectedErr: "DisplayName required ,Email required ,Password required",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()

			w := httptest.NewRecorder()
			reqBodyBytes, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest("POST", "/public/register", bytes.NewBuffer(reqBodyBytes))
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

func (suite *RestApisTestSuite) TestLoginUser() {
	var (
		email = "test@example.com"
	)
	testCases := []struct {
		name        string
		reqBody     api_gen.LoginRequest
		mock        func()
		wantStatus  int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "GivingCorrectRequest_WhenLoginSuccess_ThenReturnOk",
			reqBody: api_gen.LoginRequest{
				Email:    email,
				Password: "password",
			},
			mock: func() {
				suite.mockLoginService.EXPECT().Handle(email, "password").Return(&api_gen.LoginResponseData{
					Email: email,
				}, nil)
			},
			wantStatus:  http.StatusOK,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "GivingCorrectRequest_WhenLoginFail_ThenReturnUnauthorized",
			reqBody: api_gen.LoginRequest{
				Email:    email,
				Password: "password",
			},
			mock: func() {
				suite.mockLoginService.EXPECT().Handle(email, "password").Return(nil, errors.New("something wrong"))
			},
			wantStatus:  http.StatusUnauthorized,
			wantErr:     true,
			expectedErr: "Invalid email or password",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mock()

			w := httptest.NewRecorder()
			reqBodyBytes, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest("POST", "/public/login", bytes.NewBuffer(reqBodyBytes))
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
