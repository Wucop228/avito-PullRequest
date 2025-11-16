package http

import (
	"errors"
	"github.com/Wucop228/avito-PullRequest/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/Wucop228/avito-PullRequest/internal/models"
)

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) TeamAdd(c echo.Context) error {
	var req models.RequestTeamAdd
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
	}

	if req.TeamName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "team_name is required",
			},
		})
	}

	if err := h.svc.CreateTeamWithMembers(&req); err != nil {
		if errors.Is(err, service.ErrTeamExists) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": echo.Map{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
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

	return c.JSON(http.StatusCreated, echo.Map{
		"team": req,
	})
}

func (h *TeamHandler) TeamGet(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "team_name is required",
			},
		})
	}

	team, err := h.svc.GetTeam(teamName)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": echo.Map{
					"code":    "NOT_FOUND",
					"message": "team not found",
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
	
	return c.JSON(http.StatusOK, team)
}
