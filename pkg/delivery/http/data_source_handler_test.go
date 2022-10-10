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

type testCase struct {
	Name string
	in   domain.FbDataSource
	want error
}

func TestUpdateCheckParams(t *testing.T) {
	//tests := []testCase{
	//	{
	//		Name: "success",
	//		in: domain.FbDataSource{
	//			ID:         1,
	//			Name:       "1",
	//			SourceType: 1,
	//		},
	//		want: nil,
	//	},
	//	{
	//		Name: "id empty",
	//		in: domain.FbDataSource{
	//			ID:         0,
	//			Name:       "1",
	//			SourceType: 1,
	//		},
	//		want: domain.ParamIdEmptyErr,
	//	},
	//	{
	//		Name: "name empty",
	//		in: domain.FbDataSource{
	//			ID:         1,
	//			Name:       "",
	//			SourceType: 1,
	//		},
	//		want: domain.ParamNameEmptyErr,
	//	},
	//	{
	//		Name: "sourceType empty",
	//		in: domain.FbDataSource{
	//			ID:         1,
	//			Name:       "1",
	//			SourceType: 0,
	//		},
	//		want: domain.ParamSourceTypeEmptyErr,
	//	},
	//}

	//for _, row := range tests {
	//	err := dataSourceUpdateCheckParams(row.in)
	//	assert.Equal(t, err, row.want)
	//}
}

func TestDataSourceHandler_Delete(t *testing.T) {
	var mockDataSource domain.FbDataSource
	err := faker.FakeData(&mockDataSource)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/api/v1/dataSource/"+strconv.Itoa(int(mockDataSource.ID)), strings.NewReader(""))
	assert.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		num := uint(mockDataSource.ID)
		mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(num)))

		err = handler.Delete(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {

		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		num := uint(mockDataSource.ID)
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
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
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
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		num := mockDataSource.ID
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

func TestDataSourceHandler_Update(t *testing.T) {
	var mockDataSource domain.FbDataSource
	err := faker.FakeData(&mockDataSource)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}

		tempMock := mockDataSource
		tempMock.SourceType = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbDataSource")).Return(int64(1), nil)

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
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbDataSource")).Return(int64(1), nil)

		t.Run("StatusUnprocessableEntity", func(t *testing.T) {
			tempMock := mockDataSource
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
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbDataSource")).Return(int64(0), errors.New("update fail"))
		tempMock := mockDataSource
		tempMock.SourceType = 1
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

func TestDataSourceHandler_Store(t *testing.T) {
	var mockDataSource domain.FbDataSource
	err := faker.FakeData(&mockDataSource)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}

		tempMock := mockDataSource
		tempMock.ID = 0
		tempMock.SourceType = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbDataSource")).Return(int64(1), nil)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/dataSource/", strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/")

		err = handler.Store(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("parameter is invalid", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbDataSource")).Return(int64(1), nil)

		t.Run("StatusUnprocessableEntity", func(t *testing.T) {
			tempMock := mockDataSource
			j, err := json.Marshal(tempMock)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/dataSource/", strings.NewReader(string(j)))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/dataSource/")

			err = handler.Store(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("StatusBadRequest", func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/api/v1/dataSource/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/dataSource/")

			err = handler.Store(c)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("update fail", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}
		mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbDataSource")).Return(int64(0), errors.New("update fail"))
		tempMock := mockDataSource
		tempMock.SourceType = 1
		j, err := json.Marshal(tempMock)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/dataSource/", strings.NewReader(string(j)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/")

		err = handler.Store(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDataSourceHandler_FindDataSources(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}

		mockUCase.On("FindDataSources", mock.Anything).Return([]domain.FbDataSource{}, nil)

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/dataSource/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/")

		err = handler.FindDataSources(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("find fail", func(t *testing.T) {
		mockUCase := new(mocks.DataSourceUseCase)
		handler := DataSourceHandler{
			DataSourceUseCase: mockUCase,
		}

		mockUCase.On("FindDataSources", mock.Anything).Return(nil, errors.New("find fail"))

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/dataSource/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/dataSource/")

		err = handler.FindDataSources(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUCase.AssertExpectations(t)
	})
}
