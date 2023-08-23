package controllers

import (
	"net/http"

	"github.com/Marcel-MD/clean-api/models"
	"github.com/Marcel-MD/clean-api/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserController interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetCurrent(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Delete(ctx *gin.Context)
	AssignRole(ctx *gin.Context)
	RemoveRole(ctx *gin.Context)
}

func NewUserController(service services.UserService) UserController {
	log.Info().Msg("Creating new user controller")

	return &userController{
		service: service,
	}
}

type userController struct {
	service services.UserService
}

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Param pagination query models.PaginationQuery false "Pagination"
// @Success 200 {array} models.User
// @Router /users [get]
func (c *userController) GetAll(ctx *gin.Context) {
	query := models.PaginationQuery{}
	err := ctx.BindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := c.service.FindAll(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (c *userController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Get current user
// @Description Get current user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Router /users/current [get]
func (c *userController) GetCurrent(ctx *gin.Context) {
	id := ctx.GetString("user_id")

	user, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Register user
// @Description Register user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.RegisterUser true "User"
// @Success 200 {object} models.Token
// @Router /users/register [post]
func (c *userController) Register(ctx *gin.Context) {
	var user models.RegisterUser
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Register(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, token)
}

// @Summary Login user
// @Description Login user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.LoginUser true "User"
// @Success 200 {object} models.Token
// @Router /users/login [post]
func (c *userController) Login(ctx *gin.Context) {
	var user models.LoginUser
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, token)
}

// @Summary Delete user
// @Description Delete user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200
// @Router /users/{id} [delete]
func (c *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Assign role to user
// @Description Assign role to user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param role path string true "Role"
// @Success 200
// @Router /users/{id}/roles/{role} [patch]
func (c *userController) AssignRole(ctx *gin.Context) {
	id := ctx.Param("id")
	role := ctx.Param("role")

	err := c.service.AssignRole(id, role)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Remove role from user
// @Description Remove role from user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param role path string true "Role"
// @Success 200
// @Router /users/{id}/roles/{role} [delete]
func (c *userController) RemoveRole(ctx *gin.Context) {
	id := ctx.Param("id")
	role := ctx.Param("role")

	err := c.service.RemoveRole(id, role)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
