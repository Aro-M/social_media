package controller

import (
	"bytes"
	"encoding/json"
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

func TestController_RegisterUser(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	req := dto.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password",
		Age:       25,
	}
	user := &entity.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Age:       req.Age,
	}

	mockSvc.On("Register", mock.Anything, &req).Return(user, nil)

	body, _ := json.Marshal(req)
	httpreq := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	httpreq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpreq, rec)

	if assert.NoError(t, ctrl.RegisterUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var resp dto.UserResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, user.ID, resp.ID)
		assert.Equal(t, user.Email, resp.Email)
	}
}

func TestController_LoginUser(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: "password",
	}
	user := &entity.User{
		ID:    uuid.New(),
		Email: req.Email,
	}
	token := "some-token"

	mockSvc.On("Login", mock.Anything, &req).Return(token, user, nil)

	body, _ := json.Marshal(req)
	httpreq := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	httpreq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpreq, rec)

	if assert.NoError(t, ctrl.LoginUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp dto.LoginResponse
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, token, resp.Token)
		assert.Equal(t, user.ID, resp.User.ID)
	}
}

func TestController_UpdateProfile(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	userID := uuid.New()
	newName := "NewName"
	req := dto.UpdateProfileRequest{
		FirstName: &newName,
	}

	mockSvc.On("UpdateProfile", mock.Anything, userID, &req).Return(nil)

	body, _ := json.Marshal(req)
	httpreq := httptest.NewRequest(http.MethodPut, "/users/profile", bytes.NewBuffer(body))
	httpreq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(httpreq, rec)
	c.Set("user_id", userID.String())

	if assert.NoError(t, ctrl.UpdateProfile(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Profile updated successfully")
	}
}

func TestController_UpdateProfile_Unauthorized(t *testing.T) {
	e := echo.New()
	mockSvc := new(MockService)
	ctrl := New(mockSvc)

	req := dto.UpdateProfileRequest{}
	body, _ := json.Marshal(req)
	httpreq := httptest.NewRequest(http.MethodPut, "/users/profile", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	c := e.NewContext(httpreq, rec)
	// user_id is NOT set in context

	err := ctrl.UpdateProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
