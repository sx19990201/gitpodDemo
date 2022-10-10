package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/domain/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationHandler_Store(t *testing.T) {
	var mockAuth domain.FbAuthentication
	err := faker.FakeData(&mockAuth)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}

		tempMock := mockAuth
		tempMock.ID = 0
		tempMock.AuthSupplier = "1"
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/")

		err = handler.Store(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(1), nil)

		t.Run("StatusBadRequest", func(t *testing.T) {
			tempMock := mockAuth
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/auth/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/auth/")

			err = handler.Store(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/auth/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/auth/")

			err = handler.Store(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(0), errors.New("update fail"))
		tempMock := mockAuth
		tempMock.AuthSupplier = "1"
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/auth/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/")

		err = handler.Store(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthenticationHandler_FindAuthentication(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}

		mockUCase.On("FindAuthentication", mock.Anything).Return([]domain.FbAuthentication{}, nil)

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/auth/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/")

		err = handler.FindAuthentication(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("find fail", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}

		mockUCase.On("FindAuthentication", mock.Anything).Return(nil, errors.New("find fail"))

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/auth/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/")

		err = handler.FindAuthentication(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUCase.AssertExpectations(t)
	})
}

func TestAuthenticationHandler_Delete(t *testing.T) {
	var mockAuthentication domain.FbAuthentication
	err := faker.FakeData(&mockAuthentication)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/api/v1/auth/"+strconv.Itoa(int(mockAuthentication.ID)), strings.NewReader(""))
	assert.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		num := uint(mockAuthentication.ID)
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {

		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		num := uint(mockAuthentication.ID)
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/:id")
		c.SetParamNames("ids")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("parameter is nil", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(0)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("remove fail", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		num := mockAuthentication.ID
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(0), errors.New("remove fail"))

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))
		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthenticationHandler_Update(t *testing.T) {
	var mockAuthentication domain.FbAuthentication
	err := faker.FakeData(&mockAuthentication)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}

		tempMock := mockAuthentication
		tempMock.AuthSupplier = "1"
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/dataSource/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/")

		err = handler.Update(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(1), nil)

		t.Run("StatusBadRequest", func(t *testing.T) {
			tempMock := mockAuthentication
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/dataSource/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/dataSource/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/dataSource/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/dataSource/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockUCase := new(mocks.AuthenticationUseCase)
		handler := AuthenticationHandler{
			AuthenticationUseCase: mockUCase,
		}
		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(0), errors.New("update fail"))
		tempMock := mockAuthentication
		tempMock.AuthSupplier = "1"
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/dataSource/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/")

		err = handler.Update(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
