package controller

import (
	"fmt"
	"net/http"

	"social_media/internal/router/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func getUserID(c echo.Context) (uuid.UUID, error) {
	idStr, ok := c.Get("user_id").(string)
	if !ok || idStr == "" {
		return uuid.Nil, fmt.Errorf("missing user_id in context")
	}
	return uuid.Parse(idStr)
}

// SendFriendRequest godoc
// @Summary Send a friend request
// @Description Send a friend request to another user by their ID
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.FriendRequestDTO true "Friend request data"
// @Success 201 {object} map[string]string "Friend request sent successfully"
// @Failure 400 {object} map[string]string "Invalid payload"
// @Failure 401 {object} map[string]string "Invalid user context"
// @Router /friends/request [post]
// SendFriendRequest sends a friend invitation.
func (ctrl *Controller) SendFriendRequest(c echo.Context) error {
	senderID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	var req dto.FriendRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	if err := ctrl.svc.SendFriendRequest(c.Request().Context(), senderID, req.ReceiverID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Friend request sent successfully"})
}

// ListFriendRequests godoc
// @Summary List pending friend requests
// @Description Get all pending friend requests for the current user
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.FriendRequestResponse
// @Router /friends/requests [get]
// ListFriendRequests shows all pending requests.
func (ctrl *Controller) ListFriendRequests(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	requests, err := ctrl.svc.ListFriendRequests(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list requests"})
	}

	var resp []dto.FriendRequestResponse
	for _, req := range requests {
		resp = append(resp, dto.FriendRequestResponse{
			ID:         req.ID,
			SenderID:   req.SenderID,
			ReceiverID: req.ReceiverID,
			Status:     string(req.Status),
		})
	}
	if resp == nil {
		resp = make([]dto.FriendRequestResponse, 0)
	}

	return c.JSON(http.StatusOK, resp)
}

// AcceptFriendRequest godoc
// @Summary Accept a friend request
// @Description Accept a pending friend request by its ID
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]string "Friend request accepted"
// @Failure 400 {object} map[string]string "Invalid request ID or business error"
// @Router /friends/requests/{id}/accept [put]
// AcceptFriendRequest accepts a pending invitation.
func (ctrl *Controller) AcceptFriendRequest(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request ID"})
	}

	if err := ctrl.svc.AcceptFriendRequest(c.Request().Context(), requestID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Friend request accepted"})
}

// DeclineFriendRequest godoc
// @Summary Decline a friend request
// @Description Decline a pending friend request by its ID
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]string "Friend request declined"
// @Failure 400 {object} map[string]string "Invalid request ID or business error"
// @Router /friends/requests/{id}/decline [put]
// DeclineFriendRequest rejects a pending invitation.
func (ctrl *Controller) DeclineFriendRequest(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request ID"})
	}

	if err := ctrl.svc.DeclineFriendRequest(c.Request().Context(), requestID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Friend request declined"})
}

// ListFriends godoc
// @Summary List all friends
// @Description Get the list of all accepted friends for the current user
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse
// @Router /friends [get]
// ListFriends returns all accepted friends.
func (ctrl *Controller) ListFriends(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	friends, err := ctrl.svc.ListFriends(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list friends"})
	}

	var resp []dto.UserResponse
	for _, u := range friends {
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

// RemoveFriend godoc
// @Summary Remove a friend
// @Description Remove a friend from the user's friend list (unfriend)
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "Friend ID"
// @Success 200 {object} map[string]string "Friend removed successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /friends/{id} [delete]
// RemoveFriend deletes a friend from the list.
func (ctrl *Controller) RemoveFriend(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	friendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid friend ID"})
	}

	if err := ctrl.svc.RemoveFriend(c.Request().Context(), userID, friendID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Friend removed successfully"})
}
