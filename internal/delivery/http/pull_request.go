package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Wucop228/avito-PullRequest/internal/models"
	"github.com/Wucop228/avito-PullRequest/internal/service"
)

type PullRequestHandler struct {
	svc *service.PullRequestService
}

func NewPullRequestHandler(svc *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{svc: svc}
}

func (h *PullRequestHandler) Create(c echo.Context) error {
	var req models.RequestPullRequestCreate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "pull_request_id, pull_request_name and author_id are required",
			},
		})
	}

	pr, err := h.svc.CreatePullRequest(&req)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": echo.Map{
					"code":    "NOT_FOUND",
					"message": "author or team not found",
				},
			})
		}
		if errors.Is(err, service.ErrPRExists) {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": echo.Map{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
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
		"pr": pr,
	})
}

func (h *PullRequestHandler) Merge(c echo.Context) error {
	var req models.RequestPullRequestMerge
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
	}
	if req.PullRequestID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "pull_request_id is required",
			},
		})
	}

	pr, err := h.svc.MergePullRequest(req.PullRequestID)
	if err != nil {
		if errors.Is(err, service.ErrPRNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": echo.Map{
					"code":    "NOT_FOUND",
					"message": "pull request not found",
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
		"pr": pr,
	})
}

func (h *PullRequestHandler) Reassign(c echo.Context) error {
	var req models.RequestPullRequestReassign
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
	}
	if req.PullRequestID == "" || req.OldUserID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "pull_request_id and old_user_id are required",
			},
		})
	}

	pr, replacedBy, err := h.svc.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		if errors.Is(err, service.ErrPRNotFound) || errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": echo.Map{
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
		}
		if errors.Is(err, service.ErrPRMerged) {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": echo.Map{
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
		}
		if errors.Is(err, service.ErrReviewerNotAssigned) {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": echo.Map{
					"code":    "NOT_ASSIGNED",
					"message": "reviewer is not assigned to this PR",
				},
			})
		}
		if errors.Is(err, service.ErrNoCandidate) {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": echo.Map{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
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
		"pr":          pr,
		"replaced_by": replacedBy,
	})
}

func (h *PullRequestHandler) GetUserReviews(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"code":    "BAD_REQUEST",
				"message": "user_id is required",
			},
		})
	}

	prs, err := h.svc.GetUserReviews(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": echo.Map{
				"code":    "INTERNAL",
				"message": err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
