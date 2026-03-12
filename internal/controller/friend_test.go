package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"social_media/internal/entity"
	"social_media/internal/router/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestController_SendFriendRequest(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	userID := uuid.New()
	receiverID := uuid.New()
	reqBody := dto.FriendRequestDTO{ReceiverID: receiverID}

	mockSvc.On("SendFriendRequest", mock.Anything, userID, receiverID).Return(nil)

	body, _ := json.Marshal(reqBody)
	httpreq := httptest.NewRequest(http.MethodPost, "/friends/request", bytes.NewBuffer(body))
	httpreq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpreq, rec)
	c.Set("user_id", userID.String())

	if assert.NoError(t, ctrl.SendFriendRequest(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "successfully")
	}
}

func TestController_ListFriends(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	userID := uuid.New()
	friends := []*entity.User{
		{ID: uuid.New(), FirstName: "Friend1"},
	}

	mockSvc.On("ListFriends", mock.Anything, userID).Return(friends, nil)

	httpreq := httptest.NewRequest(http.MethodGet, "/friends", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpreq, rec)
	c.Set("user_id", userID.String())

	if assert.NoError(t, ctrl.ListFriends(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp []dto.UserResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, friends[0].ID, resp[0].ID)
	}
}

func TestController_RemoveFriend(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	userID := uuid.New()
	friendID := uuid.New()

	mockSvc.On("RemoveFriend", mock.Anything, userID, friendID).Return(nil)

	httpreq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/friends/%s", friendID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpreq, rec)
	c.Set("user_id", userID.String())
	c.SetParamNames("id")
	c.SetParamValues(friendID.String())

	if assert.NoError(t, ctrl.RemoveFriend(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "removed successfully")
	}
}
