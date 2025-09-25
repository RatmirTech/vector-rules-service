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

// CreateRule handles POST /rules
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

// GetRule handles GET /rules/:id
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

// UpdateRule handles PUT /rules/:id
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

// DeleteRule handles DELETE /rules/:id
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

// ListRules handles GET /rules
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