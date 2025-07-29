package restapis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
	"gorm.io/gorm"
)

// (POST /secure/wallet)
func (h *HttpServer) CreateWallet(ctx *gin.Context) {
	var req api_gen.WalletRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	userId := utils.GetMiddlewareUserId(ctx)

	if err := h.App.Commands.WalletService.HandleCreate(userId, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Failed to create wallet"})
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}

// (GET /secure/wallets)
func (h *HttpServer) ListUserWallets(ctx *gin.Context) {

	userId := utils.GetMiddlewareUserId(ctx)

	resp, err := h.App.Queries.ListWalletsService.Handle(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Failed to list wallets"})
		return
	}

	ctx.JSON(http.StatusOK, api_gen.ListUserWalletsResponse{
		Data: &resp,
	})
}

// (PUT /secure/wallet/{walletId})
func (h *HttpServer) UpdateWallet(ctx *gin.Context, walletId string) {
	var req api_gen.WalletRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	userId := utils.GetMiddlewareUserId(ctx)

	if err := h.App.Commands.WalletService.HandleUpdateInfo(userId, walletId, req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, api_gen.ErrorResponse{ErrorCode: "404", ErrorMessage: "Wallet not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Failed to update wallet"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// (DELETE /secure/wallet/{walletId})
func (h *HttpServer) DeleteWallet(ctx *gin.Context, walletId string) {

	userId := utils.GetMiddlewareUserId(ctx)

	if err := h.App.Commands.WalletService.HandleDelete(userId, walletId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, api_gen.ErrorResponse{ErrorCode: "404", ErrorMessage: "Wallet not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Failed to delete wallet"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
