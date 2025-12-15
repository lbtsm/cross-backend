package apis

import (
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
)

type CheckTssApi struct {
	api.Api
}

// GetCheckList 获取tss健康检查列表
// @Summary 获取tss健康检查列表
// @Description 获取tss健康检查列表
// @Tags tss健康检查管理
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CrossData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/check/tss [get]
// @Security Bearer
func (e *CheckTssApi) GetCheckList(c *gin.Context) {
	s := service.CheckTss{}
	req := dto.CrossListReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//数据权限检查
	p := actions.GetPermissionFromContext(c)
	list := make([]models.CrossInfo, 0)
	var count int64
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}
