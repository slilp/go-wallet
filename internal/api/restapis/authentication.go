package restapis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
)

// (POST /public/register)
func (h *HttpServer) RegisterUser(ctx *gin.Context) {
	var req api_gen.RegisterRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	if err := h.App.Commands.RegisterService.Handle(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, api_gen.ErrorResponse{ErrorCode: "500", ErrorMessage: "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}

// (POST /public/login)
func (h *HttpServer) LoginUser(ctx *gin.Context) {
	var req api_gen.LoginRequest
	if !utils.BindAndValidateRequestBody(ctx, &req, h.App.Utils.Validate) {
		return
	}

	resp, err := h.App.Queries.LoginService.Handle(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, api_gen.ErrorResponse{ErrorCode: "401", ErrorMessage: "Invalid email or password"})
		return
	}

	ctx.JSON(http.StatusOK, api_gen.LoginResponse{
		Data: resp,
	})
}
