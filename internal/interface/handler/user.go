package handler

import (
	"encoding/json"

	"example.com/user-management/internal/usecase"
	"github.com/savsgio/atreugo/v11"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserHandler interface {
	Create(ctx *atreugo.RequestCtx) error
	Get(ctx *atreugo.RequestCtx) error
	Update(ctx *atreugo.RequestCtx) error
	Delete(ctx *atreugo.RequestCtx) error
	List(ctx *atreugo.RequestCtx) error
}

type userHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) UserHandler {
	return &userHandler{
		userUseCase: userUseCase,
	}
}

func (h *userHandler) Create(ctx *atreugo.RequestCtx) error {
	var req createUserRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	user, err := h.userUseCase.Create(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	return ctx.JSONResponse(user, 201)
}

func (h *userHandler) Get(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	user, err := h.userUseCase.Get(ctx, id)
	if err != nil {
		return ctx.ErrorResponse(err, 404)
	}

	return ctx.JSONResponse(user, 200)
}

func (h *userHandler) Update(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	var req updateUserRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	user, err := h.userUseCase.Update(ctx, id, req.Name, req.Email)
	if err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	return ctx.JSONResponse(user, 200)
}

func (h *userHandler) Delete(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	if err := h.userUseCase.Delete(ctx, id); err != nil {
		return ctx.ErrorResponse(err, 404)
	}

	return ctx.JSONResponse(nil, 204)
}

func (h *userHandler) List(ctx *atreugo.RequestCtx) error {
	users, err := h.userUseCase.List(ctx)
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	return ctx.JSONResponse(users, 200)
}
