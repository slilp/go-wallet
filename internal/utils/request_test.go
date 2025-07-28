package utils_test

import (
	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/utils"
)

func (suite *UtilsTestSuite) TestGetPaginationParams() {
	testCases := []struct {
		name       string
		pageQuery  *int
		limitQuery *int
		wantPage   int
		wantLimit  int
	}{
		{
			name:       "NilPageAndLimitQuery_ReturnsDefaultValues",
			pageQuery:  nil,
			limitQuery: nil,
			wantPage:   1,
			wantLimit:  20,
		},
		{
			name:       "ValidPageAndNilLimitQuery_ReturnsPageAndDefaultLimit",
			pageQuery:  null.IntFrom(2).Ptr(),
			limitQuery: nil,
			wantPage:   2,
			wantLimit:  20,
		},
		{
			name:       "NilPageAndValidLimitQuery_ReturnsDefaultPageAndLimit",
			pageQuery:  nil,
			limitQuery: null.IntFrom(50).Ptr(),
			wantPage:   1,
			wantLimit:  50,
		},
		{
			name:       "ValidPageAndLimitQuery_ReturnsPageAndLimit",
			pageQuery:  null.IntFrom(3).Ptr(),
			limitQuery: null.IntFrom(30).Ptr(),
			wantPage:   3,
			wantLimit:  30,
		},
	}

	for _, tc := range testCases {

		suite.Run(tc.name, func() {
			page, limit := utils.GetPaginationParams(tc.pageQuery, tc.limitQuery)

			suite.Equal(tc.wantPage, page)
			suite.Equal(tc.wantLimit, limit)
		})

	}
}
