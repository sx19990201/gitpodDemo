package http

import (
	"database/sql"
	"fmt"
	"github.com/fire_boom/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	RoleUseCase domain.RoleUseCase
}

func InitRoleUseCaseRouter(e *echo.Echo, db *sql.DB) {
	roleRepo := sqllite.NewRoleRepository(db)
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	au := usecase.NewRoleUseCase(roleRepo, timeoutContext)
	NewRoleHandler(e, au)
}

func NewRoleHandler(e *echo.Echo, rUseCase domain.RoleUseCase) {
	handler := &RoleHandler{
		RoleUseCase: rUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/role")
		{
			authentication.GET("", handler.FindRoles)
			authentication.POST("", handler.Store)
			authentication.PUT("", handler.Update)
			authentication.DELETE("/:id", handler.Delete)
		}
	}
}

// Store @Title Store
// @Description 添加角色
// @Accept  json
// @Tags  role
// @Param id formData integer true "id"
// @Param code formData string true "编码"
// @Param remark formData integer true "说明"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/role  [POST]
func (d *RoleHandler) Store(c echo.Context) (err error) {
	var role domain.FbRole
	err = c.Bind(&role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, domain.ParamErr))
	}

	if err = roleStoreCheckParams(role); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	ctx := c.Request().Context()
	lastID, err := d.RoleUseCase.Store(ctx, &role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(lastID))
}

func roleStoreCheckParams(auth domain.FbRole) (err error) {
	if ok := utils.Empty(auth.Code); ok {
		return domain.ParamCodeEmptyErr
	}
	return
}

// FindRoles @Title FindRoles
// @Description 查询角色
// @Accept  json
// @Tags  role
// @Success 200 {object} []domain.FbRole "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/role  [GET]
func (d *RoleHandler) FindRoles(c echo.Context) (err error) {
	ctx := c.Request().Context()
	roles, err := d.RoleUseCase.FindRoles(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(roles))
}

// Delete @Title Delete
// @Description 删除角色
// @Accept  json
// @Tags  role
// @Param id formData integer true "id"
// @Success 200 "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/role/:id  [DELETE]
func (d *RoleHandler) Delete(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, domain.ParamErr))
	}

	if err = roleDeleteCheckParams(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	ctx := c.Request().Context()
	affect, err := d.RoleUseCase.Delete(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", affect)))
}

func roleDeleteCheckParams(id uint) (err error) {
	if ok := utils.Empty(id); ok {
		return domain.ParamIdEmptyErr
	}
	return
}

// Update @Title Update
// @Description 修改角色
// @Accept  json
// @Tags  role
// @Param id formData integer true "id"
// @Param code formData string true "编码"
// @Param remark formData integer true "说明"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/role  [PUT]
func (d *RoleHandler) Update(c echo.Context) (err error) {
	var role domain.FbRole
	err = c.Bind(&role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, domain.ParamErr))
	}

	if err = roleUpdateCheckParams(role); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	ctx := c.Request().Context()
	affect, err := d.RoleUseCase.Update(ctx, &role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.RoleCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", affect)))
}

func roleUpdateCheckParams(auth domain.FbRole) (result error) {
	if ok := utils.Empty(auth.ID); ok {
		return domain.ParamIdEmptyErr
	}
	if ok := utils.Empty(auth.Code); ok {
		return domain.ParamCodeEmptyErr
	}
	return
}
