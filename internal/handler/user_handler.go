package handler

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/dto"
	"github.com/ld/welfaremall/internal/middleware"
	"github.com/ld/welfaremall/internal/service"
	"github.com/ld/welfaremall/internal/util"
)

// UserHandler 用户接口处理器。
type UserHandler struct {
	userSvc service.UserService
}

// NewUserHandler 构造用户处理器。
func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// Register 注册。
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	user, err := h.userSvc.Register(req.Username, req.Password, req.RealName, req.Email, req.Phone, req.Department, req.Role)
	if err != nil {
		c.Error(fmt.Errorf("handler register: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgRegisterSuccess, user)
}

// Login 登录。
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	user, token, err := h.userSvc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, util.ErrUnauthorized) {
			c.Error(util.Unauthorized(constants.MsgPasswordIncorrect, err))
			return
		}
		c.Error(fmt.Errorf("handler login: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgLoginSuccess, dto.LoginResponse{Token: token, User: user})
}

// Me 当前用户信息。
func (h *UserHandler) Me(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	_, err = h.userSvc.GetByID(claims.UserID)
	if err != nil {
		c.Error(fmt.Errorf("handler me: %w", err))
		return
	}
	util.OK(c, nil)
}

// UpdateProfile 修改资料。
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	user, err := h.userSvc.UpdateProfile(claims.UserID, req.RealName, req.Email, req.Phone, req.Department)
	if err != nil {
		c.Error(fmt.Errorf("handler update profile: %w", err))
		return
	}
	util.OK(c, user)
}
