package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
	"net/http"
	"strconv"
	"time"
)

type ForumHandler struct {
	forumUC  forum.Usecase
	threadUC thread.Usecase
}

func NewForumHandler(e *echo.Echo, fUC forum.Usecase, tUC thread.Usecase) *ForumHandler {
	fh := &ForumHandler{
		forumUC:  fUC,
		threadUC: tUC,
	}

	e.POST("/forum/create", fh.CreateForum())
	e.POST("/forum/:fslug/create", fh.CreateThread())
	e.GET("/forum/:slug/details", fh.GetForumDetails())
	e.GET("/forum/:slug/threads", fh.GetForumThreads())
	e.GET("/forum/:slug/users", fh.GetForumUsers())

	return fh
}

func (fh *ForumHandler) CreateForum() echo.HandlerFunc {
	type createForumRequest struct {
		Slug  string `json:"slug" binding:"required"`
		Title string `json:"title" binding:"required"`
		User  string `json:"user" binding:"required"`
	}

	return func(c echo.Context) error {
		req := &createForumRequest{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		reqForum := &models.Forum{
			Slug:          req.Slug,
			Title:         req.Title,
			AdminNickname: req.User,
		}

		returnForum, err := fh.forumUC.AddForum(reqForum)
		if err != nil {
			if err == tools.ErrExistWithSlug {
				return c.JSON(http.StatusConflict, returnForum)
			}
			if err == tools.ErrUserDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, returnForum)
	}
}

func (fh *ForumHandler) CreateThread() echo.HandlerFunc {
	type CreateThreadRequest struct {
		Author  string    `json:"author" binding:"require"`
		Created time.Time `json:"created" binding:"omitempty"`
		Message string    `json:"message" binding:"require"`
		Title   string    `json:"title" binding:"require"`
		Slug    string    `json:"slug" binding:"omitempty"`
	}
	return func(c echo.Context) error {
		req := &CreateThreadRequest{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		if _, err := strconv.ParseInt(req.Slug, 10, 64); err == nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: tools.ErrIncorrectSlug.Error(),
			})
		}

		slug := c.Param("fslug")

		if req.Created.IsZero() {
			req.Created = time.Now()
		}

		reqThread := &models.Thread{
			Author:       req.Author,
			CreationDate: req.Created,
			About:        req.Message,
			Title:        req.Title,
			Slug:         req.Slug,
			Forum:        slug,
		}

		returnThread, err := fh.threadUC.AddThread(reqThread)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == tools.ErrUserDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == tools.ErrExistWithSlug {
				return c.JSON(http.StatusConflict, returnThread)
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, returnThread)
	}
}

func (fh *ForumHandler) GetForumDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")

		returnForum, err := fh.forumUC.GetForumBySlug(slug)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnForum)
	}
}

func (fh *ForumHandler) GetForumThreads() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")
		limit := uint64(0)
		var err error

		if l := c.QueryParam("limit"); l != "" {
			limit, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		since := c.QueryParam("since")

		desc := false
		if descVal := c.QueryParam("desc"); descVal == "true" {
			desc = true
		}

		returnThreads, err := fh.forumUC.GetForumThreads(slug, limit, since, desc)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnThreads)
	}
}

func (fh *ForumHandler) GetForumUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")

		limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
		since := c.QueryParam("since")
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		desc := false
		if descVal := c.QueryParam("desc"); descVal != "" {
			desc = true
		}

		returnUsers, err := fh.forumUC.GetForumUsers(slug, limit, since, desc)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnUsers)
	}
}
