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

func TestPrismaHandler_Create(t *testing.T) {
	var mockPrisma domain.Prisma
	err := faker.FakeData(&mockPrisma)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}

		tempMock := mockPrisma
		tempMock.ID = 0
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockPCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(nil)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/prisma/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/")

		err = handler.Create(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockPCase.AssertExpectations(t)
	})

}

//func TestPrismaHandler_Fetch(t *testing.T) {
//
//	t.Run("success", func(t *testing.T) {
//		mockPCase := new(mocks.PrismaUseCase)
//		handler := PrismaHandler{
//			PrismaUseCase: mockPCase,
//		}
//
//		mockPCase.On("Fetch", mock.Anything).Return([]domain.Prisma{}, nil)
//
//		e := echo.New()
//		req, err := http.NewRequest(echo.GET, "/api/v1/prisma/", strings.NewReader(""))
//		assert.NoError(t, err)
//		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//
//		rec := httptest.NewRecorder()
//		c := e.NewContext(req, rec)
//		c.SetPath("/api/v1/prisma/")
//
//		err = handler.Fetch(c)
//		require.NoError(t, err)
//		assert.Equal(t, http.StatusOK, rec.Code)
//		mockPCase.AssertExpectations(t)
//	})
//
//	t.Run("find fail", func(t *testing.T) {
//		mockPCase := new(mocks.PrismaUseCase)
//		handler := PrismaHandler{
//			PrismaUseCase: mockPCase,
//		}
//
//		mockPCase.On("Fetch", mock.Anything).Return(nil, errors.New("find fail"))
//
//		e := echo.New()
//		req, err := http.NewRequest(echo.GET, "/api/v1/prisma/", strings.NewReader(""))
//		assert.NoError(t, err)
//		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//
//		rec := httptest.NewRecorder()
//		c := e.NewContext(req, rec)
//		c.SetPath("/api/v1/prisma/")
//
//		err = handler.Fetch(c)
//		assert.Equal(t, http.StatusBadRequest, rec.Code)
//		mockPCase.AssertExpectations(t)
//	})
//}

func TestPrismaHandler_Delete(t *testing.T) {
	var mockPrisma domain.Prisma
	err := faker.FakeData(&mockPrisma)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/api/v1/prisma/"+strconv.Itoa(int(mockPrisma.ID)), strings.NewReader(""))
	assert.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}
		num := uint(mockPrisma.ID)
		mockPCase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockPCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {

		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}
		num := uint(mockPrisma.ID)
		mockPCase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/:id")
		c.SetParamNames("ids")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("parameter is nil", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}
		mockPCase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(0)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("remove fail", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}
		num := mockPrisma.ID
		mockPCase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("remove fail"))

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))
		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPrismaHandler_Update(t *testing.T) {
	var mockPrisma domain.Prisma
	err := faker.FakeData(&mockPrisma)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}

		tempMock := mockPrisma
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockPCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(nil)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/prisma/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/")

		err = handler.Update(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockPCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}
		mockPCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(nil)

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/prisma/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/prisma/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockPCase := new(mocks.PrismaUseCase)
		handler := PrismaHandler{
			PrismaUseCase: mockPCase,
		}
		mockPCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(errors.New("update fail"))
		tempMock := mockPrisma
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/prisma/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/prisma/")

		err = handler.Update(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
