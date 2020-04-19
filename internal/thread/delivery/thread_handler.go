package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
	"net/http"
	"strconv"
	"time"
)

type ThreadHandler struct {
	threadUsecase thread.Usecase
}

func NewThreadHandler(e *echo.Echo, tUC thread.Usecase) *ThreadHandler {
	th := &ThreadHandler{
		threadUsecase: tUC,
	}

	e.POST("thread/:slug_or_id/create", th.CreatePost())
	e.POST("thread/:slug_or_id/details", th.UpdateThread())
	e.POST("thread/:slug_od_id/details", th.ThreadVote())
	e.GET("thread/:slug_or_id/details", th.GetThreadDetails())
	e.GET("thread/:slug_or_id/posts", th.GetThreadPosts())

	return th
}

func (th *ThreadHandler) CreatePost() echo.HandlerFunc {
	type createPostReq struct {
		Author  string `json:"author" binding:"require"`
		Message string `json:"message" binding:"require"`
		Parent  uint64 `json:"parent" binding:"require"`
	}
	return func(c echo.Context) error {
		req := []*createPostReq{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		slugOrID := c.Param("slug_or_id")

		posts := make([]*models.Post, 0, len(req))
		createDate := time.Now()
		for _, r := range req {
			post := &models.Post{
				Author:       r.Author,
				Message:      r.Message,
				ParentID:     r.Parent,
				CreationDate: createDate,
				IsEdited:     false,
			}

			posts = append(posts, post)
		}

		returnPosts, err := th.threadUsecase.CreatePosts(slugOrID, posts)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == tools.ErrParentPostDoesntExists {
				return c.JSON(http.StatusConflict, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, returnPosts)
	}
}

func (th *ThreadHandler) GetThreadDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		returnThread, err := th.threadUsecase.GetBySlugOrID(slugOrID)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnThread)
	}
}

func (th *ThreadHandler) UpdateThread() echo.HandlerFunc {
	type updateThreadRequest struct {
		Message string `json:"message" binding:"require"`
		Title   string `json:"title" binding:"require"`
	}
	return func(c echo.Context) error {
		req := &updateThreadRequest{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		slugOrID := c.Param("slug_or_id")
		reqThread := &models.Thread{
			About: req.Message,
			Title: req.Title,
		}

		returnThread, err := th.threadUsecase.Update(slugOrID, reqThread)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnThread)
	}
}

func (th *ThreadHandler) ThreadVote() echo.HandlerFunc {
	type voteReq struct {
		Nickname string `json:"nickname" binding:"require"`
		Voice    int64  `json:"voice" binding:"require"`
	}
	return func(c echo.Context) error {
		req := &voteReq{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		slugOrID := c.Param("slug_or_id")

		vReq := &models.Vote{
			Nickname: req.Nickname,
			Voice:    req.Voice,
		}

		returnThread, err := th.threadUsecase.Vote(slugOrID, vReq)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnThread)
	}
}

func (th *ThreadHandler) GetThreadPosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
		since, err := strconv.ParseUint(c.QueryParam("since"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		sort := c.QueryParam("sort")
		desc := false
		if descVal := c.QueryParam("desc"); descVal != "" {
			desc = true
		}

		returnPosts, err := th.threadUsecase.GetThreadPosts(slugOrID, limit, since, sort, desc)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnPosts)
	}
}
