package handlers

import (
	"net/http"

	"bankmore/internal/shared/models"
	"bankmore/internal/transfer/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransferHandler struct {
	service service.TransferService
	logger  *logrus.Logger
}

func NewTransferHandler(service service.TransferService, logger *logrus.Logger) *TransferHandler {
	return &TransferHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Realiza transferência entre contas
// @Description Realiza transferência entre contas da mesma instituição
// @Tags Transfer
// @Accept json
// @Produce json
// @Param request body service.CreateTransferRequest true "Dados da transferência"
// @Success 200 {object} service.TransferResponse
// @Failure 400 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /api/transfer [post]
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var request service.CreateTransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: "Dados inválidos",
		})
		return
	}

	accountID, exists := c.Get("accountId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Type:    models.ErrorUserUnauthorized,
			Message: "Token inválido",
		})
		return
	}

	result, err := h.service.CreateTransfer(request, accountID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Error creating transfer")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Type:    models.ErrorInternalError,
			Message: "Erro interno do servidor",
		})
		return
	}

	if !result.IsSuccess {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    result.ErrorType,
			Message: result.ErrorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, result.Data)
}
