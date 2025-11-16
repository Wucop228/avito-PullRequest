package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Wucop228/avito-PullRequest/internal/models"
	"github.com/Wucop228/avito-PullRequest/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SetIsActive(c echo.Context) error {
	var req models.RequestSetIsActive
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "user_id is required",
			},
		})
	}

	user, err := h.svc.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": echo.Map{
					"code":    "NOT_FOUND",
					"message": "user not found",
				},
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": echo.Map{
				"code":    "INTERNAL",
				"message": err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": models.User{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	})
}
