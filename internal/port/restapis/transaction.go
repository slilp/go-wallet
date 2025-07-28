package restapis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
)

// (GET /secure/wallet/{walletId}/transactions)
func (h *HttpServer) ListWalletTransactions(ctx *gin.Context, walletId string, params api_gen.ListWalletTransactionsParams) {

	// page, limit := utils.GetPaginationParams(params.Page, params.Limit)

	ctx.JSON(http.StatusOK, nil)
}

// (POST /secure/transfer)
func (h *HttpServer) TransferBalance(ctx *gin.Context) {
	var req api_gen.TransferRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// (POST /secure/deposit)
func (h *HttpServer) DepositPoints(ctx *gin.Context) {
	var req api_gen.DepositRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// (POST /secure/withdraw)
func (h *HttpServer) WithdrawPoints(ctx *gin.Context) {
	var req api_gen.WithdrawRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
