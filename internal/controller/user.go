package controller

import (
	"errors"
	"net/http"

	"social_media/internal/router/dto"
	"social_media/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account with basic profile information
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/register [post]
// RegisterUser creates a new user account.
func (ctrl *Controller) RegisterUser(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	user, err := ctrl.svc.Register(c.Request().Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrFirstNameEmpty),
			errors.Is(err, service.ErrLastNameEmpty),
			errors.Is(err, service.ErrPasswordTooShort),
			errors.Is(err, service.ErrInvalidAge):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		Email:     user.Email,
	}

	return c.JSON(http.StatusCreated, resp)
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Router /users/login [post]
// LoginUser verifies user and returns a token.
func (ctrl *Controller) LoginUser(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	token, user, err := ctrl.svc.Login(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	resp := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Age:       user.Age,
			Email:     user.Email,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

// SearchUsers godoc
// @Summary Search users
// @Description Find users by first name, last name, or age
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param first_name query string false "First Name"
// @Param last_name query string false "Last Name"
// @Param age query int false "Age"
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} map[string]string "Failed to search users"
// @Router /users/search [get]
// SearchUsers finds users by name or age.
func (ctrl *Controller) SearchUsers(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var req dto.SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query parameters"})
	}

	users, err := ctrl.svc.SearchUsers(c.Request().Context(), req.FirstName, req.LastName, req.Age, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search users"})
	}

	var resp []dto.UserResponse
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Age:       u.Age,
			Email:     u.Email,
		})
	}

	if resp == nil {
		resp = make([]dto.UserResponse, 0)
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile information (first name, last name, age, or password) of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} map[string]string "Profile updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/profile [put]
// UpdateProfile changes the current user's profile.
func (ctrl *Controller) UpdateProfile(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := ctrl.svc.UpdateProfile(c.Request().Context(), userID, &req); err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrPasswordTooShort),
			errors.Is(err, service.ErrInvalidAge):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
