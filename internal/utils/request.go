package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
)

func GetPaginationParams(pageQuery, limitQuery *int) (int, int) {
	page := 1
	limit := 20

	if pageQuery != nil {
		page = *pageQuery
	}

	if limitQuery != nil {
		limit = *limitQuery
	}

	return page, limit
}

func BindAndValidateRequestBody(ctx *gin.Context, req interface{}, validate *validator.Validate) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: err.Error()})
		return false
	}
	if err := validate.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" "+err.Tag()+" "+err.Param())
		}
		ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: strings.TrimSpace(strings.Join(errors, ","))})
		return false
	}
	return true
}
