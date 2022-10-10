package middleware

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	Conn *sql.DB
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		c.Response().Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		c.Response().Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		return next(c)
	}
}

// LOGS will handle the LOGS middleware
func (m *GoMiddleware) LOGS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware(db *sql.DB) *GoMiddleware {
	return &GoMiddleware{
		Conn: db,
	}
}

// RequestID will handle the LOGS middleware
func (m *GoMiddleware) RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := uuid.New().String()
		c.Request().Header.Set(echo.HeaderXRequestID, requestId)
		if c.Request().RequestURI != "/api/v1/wdg/log" {
			return next(c)
		}
		return next(c)
	}
}

// UpdateWunderGraph will handle the UpdateWunderGraph middleware
func (m *GoMiddleware) UpdateWunderGraph(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := c.Request().Context()
		next(c)
		//wdg, err := wundergraph.GetWdg(m.Conn)
		//if err != nil {
		//	// TODO err
		//	return c.JSON(http.StatusBadRequest, "")
		//}
		//ctx := c.Request().Context()
		//wdg.ReloadConfig(ctx)
		//wdg.ReloadOpertionTs(ctx)
		//wdg.ReloadServerTs(ctx)

		return nil
	}
}
