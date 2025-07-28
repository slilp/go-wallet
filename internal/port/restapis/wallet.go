package restapis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
)

// (POST /secure/wallet)
func (h *HttpServer) CreateWallet(ctx *gin.Context) {
	var req api_gen.WalletRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

// (GET /secure/wallets)
func (h *HttpServer) ListUserWallets(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, nil)
}

// (PUT /secure/wallet/{walletId})
func (h *HttpServer) UpdateWallet(ctx *gin.Context, walletId string) {
	var req api_gen.WalletRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// (DELETE /secure/wallet/{walletId})
func (h *HttpServer) DeleteWallet(ctx *gin.Context, walletId string) {

	ctx.Status(http.StatusNoContent)
}
