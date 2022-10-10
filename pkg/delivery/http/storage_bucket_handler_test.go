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

func TestStorageBucketHandler_Delete(t *testing.T) {
	var mockStorageBucket domain.FbStorageBucket
	err := faker.FakeData(&mockStorageBucket)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/api/v1/storageBucket/"+strconv.Itoa(int(mockStorageBucket.ID)), strings.NewReader(""))
	assert.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		num := uint(mockStorageBucket.ID)
		mockSCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {

		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		num := uint(mockStorageBucket.ID)
		mockSCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/:id")
		c.SetParamNames("ids")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("parameter is nil", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		mockSCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(0)))

		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("remove fail", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		num := mockStorageBucket.ID
		mockSCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(0), errors.New("remove fail"))

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))
		err = handler.Delete(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStorageBucketHandler_Update(t *testing.T) {
	var mockStorageBucket domain.FbStorageBucket
	err := faker.FakeData(&mockStorageBucket)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}

		tempMock := mockStorageBucket
		tempMock.Switch = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockSCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/storageBucket/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/")

		err = handler.Update(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		mockSCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(1), nil)

		t.Run("StatusUnprocessableEntity", func(t *testing.T) {
			tempMock := mockStorageBucket
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/storageBucket/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/storageBucket/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.PUT, "/api/v1/storageBucket/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/storageBucket/")

			err = handler.Update(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		mockSCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(0), errors.New("update fail"))
		tempMock := mockStorageBucket
		tempMock.Switch = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/storageBucket/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/")

		err = handler.Update(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStorageBucketHandler_Store(t *testing.T) {
	var mockStorageBucket domain.FbStorageBucket
	err := faker.FakeData(&mockStorageBucket)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}

		tempMock := mockStorageBucket
		tempMock.ID = 0
		tempMock.Switch = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockSCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/storageBucket/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/")

		err = handler.Store(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		mockSCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(1), nil)

		t.Run("StatusUnprocessableEntity", func(t *testing.T) {
			tempMock := mockStorageBucket
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/storageBucket/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/storageBucket/")

			err = handler.Store(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/storageBucket/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/storageBucket/")

			err = handler.Store(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}
		mockSCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(0), errors.New("update fail"))
		tempMock := mockStorageBucket
		tempMock.Switch = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/storageBucket/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/")

		err = handler.Store(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStorageBucketHandler_FindStorageBuckets(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}

		mockSCase.On("FindStorageBucket", mock.Anything).Return([]domain.FbStorageBucket{}, nil)

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/storageBucket/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/")

		err = handler.FindStorageBuckets(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSCase.AssertExpectations(t)
	})

	t.Run("find fail", func(t *testing.T) {
		mockSCase := new(mocks.StorageBucketUseCase)
		handler := StorageBucketHandler{
			StorageBucketUseCase: mockSCase,
		}

		mockSCase.On("FindStorageBucket", mock.Anything).Return(nil, errors.New("find fail"))

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/storageBucket/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/storageBucket/")

		err = handler.FindStorageBuckets(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockSCase.AssertExpectations(t)
	})
}
