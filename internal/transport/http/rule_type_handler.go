package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

// RuleTypeHandler handles HTTP requests for rule types
type RuleTypeHandler struct {
	ruleTypeService domain.RuleTypeService
}

// NewRuleTypeHandler creates a new rule type handler
func NewRuleTypeHandler(ruleTypeService domain.RuleTypeService) *RuleTypeHandler {
	return &RuleTypeHandler{
		ruleTypeService: ruleTypeService,
	}
}

// CreateRuleType creates a new rule type
// @Summary Create a new rule type
// @Description Create a new rule type
// @Tags rule-types
// @Accept json
// @Produce json
// @Param ruleType body SwaggerCreateRuleTypeRequest true "Rule type creation request"
// @Success 201 {object} SwaggerRuleType
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rule-types [post]
func (h *RuleTypeHandler) CreateRuleType(c echo.Context) error {
	var req domain.CreateRuleTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	ruleType, err := h.ruleTypeService.CreateRuleType(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, ruleType)
}

// GetRuleType retrieves a rule type by ID
// @Summary Get rule type by ID
// @Description Get a specific rule type by its ID
// @Tags rule-types
// @Produce json
// @Param id path int true "Rule type ID"
// @Success 200 {object} SwaggerRuleType
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rule-types/{id} [get]
func (h *RuleTypeHandler) GetRuleType(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid rule type id"})
	}

	ruleType, err := h.ruleTypeService.GetRuleType(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRuleTypeNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "rule type not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ruleType)
}

// UpdateRuleType updates an existing rule type
// @Summary Update a rule type
// @Description Update an existing rule type
// @Tags rule-types
// @Accept json
// @Produce json
// @Param id path int true "Rule type ID"
// @Param ruleType body SwaggerUpdateRuleTypeRequest true "Rule type update request"
// @Success 200 {object} SwaggerRuleType
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rule-types/{id} [put]
func (h *RuleTypeHandler) UpdateRuleType(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid rule type id"})
	}

	var req domain.UpdateRuleTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	req.ID = id

	ruleType, err := h.ruleTypeService.UpdateRuleType(c.Request().Context(), &req)
	if err != nil {
		if err == domain.ErrRuleTypeNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "rule type not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ruleType)
}

// DeleteRuleType deletes a rule type
// @Summary Delete a rule type
// @Description Delete a rule type by ID
// @Tags rule-types
// @Param id path int true "Rule type ID"
// @Success 204 "No content"
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rule-types/{id} [delete]
func (h *RuleTypeHandler) DeleteRuleType(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid rule type id"})
	}

	err = h.ruleTypeService.DeleteRuleType(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRuleTypeNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "rule type not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ListRuleTypes lists rule types with pagination
// @Summary List rule types
// @Description List rule types with optional pagination
// @Tags rule-types
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} SwaggerListResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rule-types [get]
func (h *RuleTypeHandler) ListRuleTypes(c echo.Context) error {
	// Parse query parameters
	limitStr := c.QueryParam("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offsetStr := c.QueryParam("offset")
	offset := 0 // default
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	ruleTypes, err := h.ruleTypeService.ListRuleTypes(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"rule_types": ruleTypes,
		"limit":      limit,
		"offset":     offset,
	})
}
