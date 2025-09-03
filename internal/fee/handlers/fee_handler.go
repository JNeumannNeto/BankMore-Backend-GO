package handlers

import (
	"net/http"
	"strconv"

	"bankmore/internal/fee/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type FeeHandler struct {
	service service.FeeService
	logger  *logrus.Logger
}

func NewFeeHandler(service service.FeeService, logger *logrus.Logger) *FeeHandler {
	return &FeeHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Consulta tarifas por número da conta
// @Description Consulta todas as tarifas de uma conta pelo número
// @Tags Fee
// @Produce json
// @Param accountNumber path string true "Número da conta"
// @Success 200 {array} domain.Fee
// @Failure 500 {string} string
// @Router /api/fee/{accountNumber} [get]
func (h *FeeHandler) GetFeesByAccount(c *gin.Context) {
	accountNumber := c.Param("accountNumber")

	fees, err := h.service.GetFeesByAccountNumber(accountNumber)
	if err != nil {
		h.logger.WithError(err).Error("Error getting fees for account")
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, fees)
}

// @Summary Consulta tarifa específica por ID
// @Description Consulta uma tarifa específica pelo ID
// @Tags Fee
// @Produce json
// @Param id path int true "ID da tarifa"
// @Success 200 {object} domain.Fee
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/fee/fee/{id} [get]
func (h *FeeHandler) GetFeeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid fee ID")
		return
	}

	fee, err := h.service.GetFeeByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Error getting fee by ID")
		c.JSON(http.StatusNotFound, "Fee not found")
		return
	}

	c.JSON(http.StatusOK, fee)
}
