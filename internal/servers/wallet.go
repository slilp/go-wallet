package servers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/slilp/go-wallet/internal/servers/api_gen"
)

// (POST /secure/wallet)
func (h *HttpServer) CreateWallet(ctx *gin.Context) {
	var req api_gen.WalletRequest
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

	ctx.JSON(http.StatusCreated, nil)
}

// (GET /secure/wallet)
func (h *HttpServer) ListUserWallets(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, nil)
}

// (GET /secure/wallet/{walletId})
func (h *HttpServer) GetWalletById(ctx *gin.Context, walletId string) {

	ctx.JSON(http.StatusOK, nil)
}

// (PUT /secure/wallet/{walletId})
func (h *HttpServer) UpdateWallet(ctx *gin.Context, walletId string) {

	var req api_gen.WalletRequest
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

// (DELETE /secure/wallet/{walletId})
func (h *HttpServer) DeleteWallet(ctx *gin.Context, walletId string) {

	ctx.Status(http.StatusNoContent)
}
