package servers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/slilp/go-wallet/internal/servers/api_gen"
)

// (GET /secure/wallet/{walletId}/transactions)
func (h *HttpServer) ListWalletTransactions(ctx *gin.Context, walletId string, params api_gen.ListWalletTransactionsParams) {

	// page, limit := utils.GetPaginationParams(params.Page, params.Limit)

	ctx.JSON(http.StatusOK, nil)
}

// (POST /secure/transfer)
func (h *HttpServer) TransferBalance(ctx *gin.Context) {
	var req api_gen.TransferRequest
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

// (POST /secure/deposit)
func (h *HttpServer) DepositMoney(ctx *gin.Context) {
	var req api_gen.DepositRequest
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

// (POST /secure/withdraw)
func (h *HttpServer) WithdrawMoney(ctx *gin.Context) {
	var req api_gen.WithdrawRequest
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

// (GET /secure/transaction/{transactionId})
func (h *HttpServer) GetTransactionDetails(ctx *gin.Context, transactionId string) {

	ctx.JSON(http.StatusOK, nil)
}

// (PATCH /secure/transaction/{transactionId})
func (h *HttpServer) UpdateTransactionDescription(ctx *gin.Context, transactionId string) {

	var req api_gen.UpdateTransactionDescRequest
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
