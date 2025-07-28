package restapis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
)

// (POST /public/register)
func (h *HttpServer) RegisterUser(ctx *gin.Context) {
	var req api_gen.RegisterRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

// (POST /public/login)
func (h *HttpServer) LoginUser(ctx *gin.Context) {
	var req api_gen.LoginRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.app.Utils.Validate) {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
