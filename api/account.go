package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=EUR USD"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.ListAccount(ctx, db.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type DeleteAccountRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type UpdateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" bomdomg:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req UpdateAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.UpdateAccount(ctx, db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
