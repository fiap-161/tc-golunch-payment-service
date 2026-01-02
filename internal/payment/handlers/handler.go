package handlers

import (
	"net/http"

	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/controllers"
	dto "github.com/fiap-161/tc-golunch-payment-service/internal/payment/dto"
	apperror "github.com/fiap-161/tc-golunch-payment-service/internal/shared/errors"
	"github.com/fiap-161/tc-golunch-payment-service/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	controller *controllers.Controller
}

func New(controller *controllers.Controller) *Handler {
	return &Handler{
		controller: controller,
	}
}

// CheckPayment godoc
// @Summary      Check Payment [Mercado Pago Integration]
// @Description  Check the status of a payment by its resource URL
// @Tags         Payment Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CheckPaymentRequestDTO true "Resource URL to check payment status"
// @Success      200
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      500  {object}  errors.ErrorDTO
// @Router       /webhook/payment/check [post]
func (h *Handler) CheckPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var checkPaymentDTO dto.CheckPaymentRequestDTO
	if err := c.ShouldBindJSON(&checkPaymentDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	_, err := h.controller.CheckPayment(ctx, checkPaymentDTO.Resource)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
