package handlers

import (
	"net/http"

	"bankmore/internal/account/service"
	"bankmore/internal/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	service service.AccountService
	logger  *logrus.Logger
}

func NewAccountHandler(service service.AccountService, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Cadastra uma nova conta corrente
// @Description Cadastra uma nova conta corrente no sistema
// @Tags Account
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Dados para criação da conta"
// @Success 200 {object} service.RegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/account/register [post]
func (h *AccountHandler) Register(c *gin.Context) {
	var request service.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: "Dados inválidos",
		})
		return
	}

	response, err := h.service.Register(request)
	if err != nil {
		h.logger.WithError(err).Error("Error registering account")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Realiza login na conta corrente
// @Description Realiza login e retorna token JWT
// @Tags Account
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Dados de login"
// @Success 200 {object} service.LoginResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/account/login [post]
func (h *AccountHandler) Login(c *gin.Context) {
	var request service.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: "Dados inválidos",
		})
		return
	}

	response, err := h.service.Login(request)
	if err != nil {
		h.logger.WithError(err).Error("Error logging in")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Type:    models.ErrorUserUnauthorized,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Inativa a conta corrente
// @Description Inativa a conta corrente do usuário logado
// @Tags Account
// @Accept json
// @Produce json
// @Param request body DeactivateRequest true "Senha para confirmação"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /api/account/deactivate [put]
func (h *AccountHandler) Deactivate(c *gin.Context) {
	type DeactivateRequest struct {
		Password string `json:"password" binding:"required"`
	}

	var request DeactivateRequest
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

	err := h.service.Deactivate(accountID.(string), request.Password)
	if err != nil {
		h.logger.WithError(err).Error("Error deactivating account")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Realiza movimentação na conta corrente
// @Description Realiza depósito ou saque na conta corrente
// @Tags Account
// @Accept json
// @Produce json
// @Param request body service.MovementRequest true "Dados da movimentação"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /api/account/movement [post]
func (h *AccountHandler) CreateMovement(c *gin.Context) {
	var request service.MovementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: "Dados inválidos",
		})
		return
	}

	err := h.service.CreateMovement(request)
	if err != nil {
		h.logger.WithError(err).Error("Error creating movement")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Consulta o saldo da conta corrente
// @Description Consulta o saldo da conta corrente do usuário logado
// @Tags Account
// @Produce json
// @Success 200 {object} service.BalanceResponse
// @Failure 400 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /api/account/balance [get]
func (h *AccountHandler) GetBalance(c *gin.Context) {
	accountID, exists := c.Get("accountId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Type:    models.ErrorUserUnauthorized,
			Message: "Token inválido",
		})
		return
	}

	response, err := h.service.GetBalance(accountID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Error getting balance")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Type:    models.ErrorInvalidData,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Verifica se uma conta existe pelo número
// @Description Verifica se uma conta existe pelo número
// @Tags Account
// @Produce json
// @Param accountNumber path string true "Número da conta"
// @Success 200 {boolean} bool
// @Router /api/account/exists/{accountNumber} [get]
func (h *AccountHandler) AccountExists(c *gin.Context) {
	accountNumber := c.Param("accountNumber")

	exists, err := h.service.AccountExists(accountNumber)
	if err != nil {
		h.logger.WithError(err).Error("Error checking account existence")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Type:    models.ErrorInternalError,
			Message: "Erro interno do servidor",
		})
		return
	}

	c.JSON(http.StatusOK, exists)
}

// @Summary Consulta o saldo de uma conta pelo número
// @Description Consulta o saldo de uma conta pelo número
// @Tags Account
// @Produce json
// @Param accountNumber path string true "Número da conta"
// @Success 200 {object} service.BalanceResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/account/balance/{accountNumber} [get]
func (h *AccountHandler) GetBalanceByAccountNumber(c *gin.Context) {
	accountNumber := c.Param("accountNumber")

	response, err := h.service.GetBalanceByAccountNumber(accountNumber)
	if err != nil {
		h.logger.WithError(err).Error("Error getting balance by account number")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Type:    models.ErrorInternalError,
			Message: "Erro interno do servidor",
		})
		return
	}

	if response == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Type:    models.ErrorAccountNotFound,
			Message: "Conta não encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
