package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

// RuleHandler handles HTTP requests for rules
type RuleHandler struct {
	ruleService domain.RuleService
}

// NewRuleHandler creates a new rule handler
func NewRuleHandler(ruleService domain.RuleService) *RuleHandler {
	return &RuleHandler{
		ruleService: ruleService,
	}
}

// CreateRule creates a new rule
// @Summary Create a new rule
// @Description Create a new rule with embedding generation
// @Tags rules
// @Accept json
// @Produce json
// @Param rule body SwaggerCreateRuleRequest true "Rule creation request"
// @Success 201 {object} SwaggerRule
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rules [post]
func (h *RuleHandler) CreateRule(c echo.Context) error {
	var req domain.CreateRuleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	rule, err := h.ruleService.CreateRule(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, rule)
}

// GetRule retrieves a rule by ID
// @Summary Get rule by ID
// @Description Get a specific rule by its ID
// @Tags rules
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} SwaggerRule
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rules/{id} [get]
func (h *RuleHandler) GetRule(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid rule id"})
	}

	rule, err := h.ruleService.GetRule(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "rule not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rule)
}

// UpdateRule updates an existing rule
// @Summary Update a rule
// @Description Update an existing rule and regenerate embedding
// @Tags rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Param rule body SwaggerUpdateRuleRequest true "Rule update request"
// @Success 200 {object} SwaggerRule
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rules/{id} [put]
func (h *RuleHandler) UpdateRule(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid rule id"})
	}

	var req domain.UpdateRuleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	req.ID = id

	rule, err := h.ruleService.UpdateRule(c.Request().Context(), &req)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "rule not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rule)
}

// DeleteRule deletes a rule
// @Summary Delete a rule
// @Description Delete a rule by ID
// @Tags rules
// @Param id path int true "Rule ID"
// @Success 204 "No content"
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rules/{id} [delete]
func (h *RuleHandler) DeleteRule(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid rule id"})
	}

	err = h.ruleService.DeleteRule(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "rule not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ListRules lists rules with pagination
// @Summary List rules
// @Description List rules with optional pagination
// @Tags rules
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} SwaggerListResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /rules [get]
func (h *RuleHandler) ListRules(c echo.Context) error {
	// Parse query parameters
	ruleType := c.QueryParam("type")
	var ruleTypePtr *string
	if ruleType != "" {
		ruleTypePtr = &ruleType
	}

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

	rules, err := h.ruleService.ListRules(c.Request().Context(), ruleTypePtr, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"rules":  rules,
		"limit":  limit,
		"offset": offset,
	})
}
