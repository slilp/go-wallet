package utils_test

import (
	"time"

	"github.com/slilp/go-wallet/internal/utils"
)

func (suite *UtilsTestSuite) TestGenerateToken() {
	testCases := []struct {
		name        string
		userID      string
		tokenType   string
		tokenTime   int
		expectError bool
	}{
		{
			name:        "ValidAccessToken_ReturnsTokenString",
			userID:      "user123",
			tokenType:   "access",
			tokenTime:   30,
			expectError: false,
		},
		{
			name:        "ValidRefreshToken_ReturnsTokenString",
			userID:      "user456",
			tokenType:   "refresh",
			tokenTime:   10080,
			expectError: false,
		},
		{
			name:        "EmptyUserID_ReturnsTokenString",
			userID:      "",
			tokenType:   "access",
			tokenTime:   15,
			expectError: false,
		},
		{
			name:        "ZeroTokenTime_ReturnsExpiredToken",
			userID:      "user789",
			tokenType:   "access",
			tokenTime:   0,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tokenString, err := utils.GenerateToken(tc.userID, tc.tokenType, tc.tokenTime)

			if tc.expectError {
				suite.Error(err)
				suite.Empty(tokenString)
			} else {
				suite.NoError(err)
				suite.NotEmpty(tokenString)

				// Verify token can be parsed
				claims, parseErr := utils.ValidateToken(tokenString)
				if tc.tokenTime > 0 {
					suite.NoError(parseErr)
					suite.Equal(tc.userID, claims.UserID)
					suite.Equal(tc.tokenType, claims.TokenType)
					suite.True(claims.ExpiresAt.After(time.Now()))
				}
			}
		})
	}
}

func (suite *UtilsTestSuite) TestValidateToken() {
	// Generate a valid token for testing
	validToken, _ := utils.GenerateToken("testuser", "access", 30)

	// Generate an expired token
	expiredToken, _ := utils.GenerateToken("expireduser", "access", -1)

	testCases := []struct {
		name        string
		tokenString string
		expectError bool
		expectedErr string
	}{
		{
			name:        "ValidToken_ReturnsClaims",
			tokenString: validToken,
			expectError: false,
		},
		{
			name:        "ExpiredToken_ReturnsError",
			tokenString: expiredToken,
			expectError: true,
			expectedErr: "token is expired",
		},
		{
			name:        "InvalidTokenFormat_ReturnsError",
			tokenString: "invalid.token.format",
			expectError: true,
		},
		{
			name:        "EmptyToken_ReturnsError",
			tokenString: "",
			expectError: true,
		},
		{
			name:        "MalformedToken_ReturnsError",
			tokenString: "not-a-jwt-token",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			claims, err := utils.ValidateToken(tc.tokenString)

			if tc.expectError {
				suite.Error(err)
				suite.Nil(claims)
				if tc.expectedErr != "" {
					suite.Contains(err.Error(), tc.expectedErr)
				}
			} else {
				suite.NoError(err)
				suite.NotNil(claims)
				suite.Equal("testuser", claims.UserID)
				suite.Equal("access", claims.TokenType)
			}
		})
	}
}
