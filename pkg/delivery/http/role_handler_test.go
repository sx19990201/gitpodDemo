package http

import (
	"encoding/json"
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/domain/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRoleHandler_Store(t *testing.T) {
	var mockRole domain.FbRole
	err := faker.FakeData(&mockRole)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}

		tempMock := mockRole
		tempMock.ID = 0
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/role/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/")

		err = handler.Store(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(1), nil)

		t.Run("StatusUnprocessableEntity", func(t *testing.T) {
			tempMock := mockRole
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/role/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/role/")

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
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(0), errors.New("update fail"))
		tempMock := mockRole
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

func TestRoleHandler_FindRoles(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}

		mockUCase.On("FindRoles", mock.Anything).Return([]domain.FbRole{}, nil)

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/role/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/")

		err = handler.FindRoles(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("find fail", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}

		mockUCase.On("FindRoles", mock.Anything).Return(nil, errors.New("find fail"))

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/role/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/")

		err = handler.FindRoles(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUCase.AssertExpectations(t)
	})
}

func TestRoleHandler_Delete(t *testing.T) {
	var mockRole domain.FbRole
	err := faker.FakeData(&mockRole)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/api/v1/role/"+strconv.Itoa(int(mockRole.ID)), strings.NewReader(""))
	assert.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		num := uint(mockRole.ID)
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {

		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		num := uint(mockRole.ID)
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/:id")
		c.SetParamNames("ids")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("parameter is nil", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(0)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("remove fail", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		num := mockRole.ID
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(0), errors.New("remove fail"))

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))
		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestRoleHandler_Update(t *testing.T) {
	var mockRole domain.FbRole
	err := faker.FakeData(&mockRole)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}

		tempMock := mockRole
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/role/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/")

		err = handler.Update(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(1), nil)

		t.Run("StatusUnprocessableEntity", func(t *testing.T) {
			tempMock := mockRole
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/role/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/role/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/role/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/role/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockUCase := new(mocks.RoleUseCase)
		handler := RoleHandler{
			RoleUseCase: mockUCase,
		}
		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(0), errors.New("update fail"))
		tempMock := mockRole
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/role/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/role/")

		err = handler.Update(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
