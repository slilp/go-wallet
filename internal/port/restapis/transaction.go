package restapis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/consts"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
)

// (GET /secure/wallet/{walletId}/transactions)
func (h *HttpServer) ListWalletTransactions(ctx *gin.Context, walletId string, params api_gen.ListWalletTransactionsParams) {

	page, limit := utils.GetPaginationParams(params.Page, params.Limit)

	userId := utils.GetMiddlewareUserId(ctx)

	totalCount, listData, err := h.App.Queries.ListTransactionsService.Handle(userId, walletId, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Failed to list wallets"})
		return
	}

	ctx.JSON(http.StatusOK, api_gen.ListWalletTransactionsResponse{
		Data: &listData,
		Pagination: &api_gen.PageLimitResponseData{
			Page:         page,
			Limit:        limit,
			TotalRecords: int(totalCount),
		},
	})
}

// (POST /secure/transfer)
func (h *HttpServer) TransferBalance(ctx *gin.Context) {
	var req api_gen.TransferRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	userId := utils.GetMiddlewareUserId(ctx)

	if err := h.App.Commands.TransactionService.HandleTransferBalance(userId, req.FromWalletId, req.ToWalletId, req.Amount); err != nil {
		if errors.Is(err, consts.ErrInsufficientBalance) {
			ctx.JSON(http.StatusBadRequest, api_gen.ErrorResponse{ErrorCode: "400", ErrorMessage: "Insufficient balance"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Fail to transfer"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// (POST /secure/deposit)
func (h *HttpServer) DepositPoints(ctx *gin.Context) {
	var req api_gen.DepositRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	userId := utils.GetMiddlewareUserId(ctx)

	if err := h.App.Commands.TransactionService.HandleDepositWithDrawBalance(userId, req.WalletId, req.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Fail to deposit"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// (POST /secure/withdraw)
func (h *HttpServer) WithdrawPoints(ctx *gin.Context) {
	var req api_gen.WithdrawRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	userId := utils.GetMiddlewareUserId(ctx)

	if err := h.App.Commands.TransactionService.HandleDepositWithDrawBalance(userId, req.WalletId, -req.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Fail to withdraw"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
