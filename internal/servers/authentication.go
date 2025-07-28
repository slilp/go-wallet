package servers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/slilp/go-wallet/internal/servers/api_gen"
)

// (POST /public/register)
func (h *HttpServer) RegisterUser(ctx *gin.Context) {
	var req api_gen.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: err.Error()})
		return
	}

	if err := h.Utils.Validate.Struct(&req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" "+err.Tag()+" "+err.Param())
		}
		ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: strings.TrimSpace(strings.Join(errors, ","))})
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}

// (POST /public/login)
func (h *HttpServer) LoginUser(ctx *gin.Context) {
	var req api_gen.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: err.Error()})
		return
	}

	if err := h.Utils.Validate.Struct(&req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" "+err.Tag()+" "+err.Param())
		}
		ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: strings.Join(errors, ",")})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
