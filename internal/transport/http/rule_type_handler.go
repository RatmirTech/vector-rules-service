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

// CreateRuleType handles POST /rule-types
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

// GetRuleType handles GET /rule-types/:id
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

// UpdateRuleType handles PUT /rule-types/:id
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

// DeleteRuleType handles DELETE /rule-types/:id
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

// ListRuleTypes handles GET /rule-types
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