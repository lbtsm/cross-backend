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

type CrossApi struct {
	api.Api
}

// GetPage 获取跨链数据列表
// @Summary 获取跨链数据列表
// @Description 获取跨链数据列表
// @Tags 跨链数据管理
// @Param srcChain query string false "来源链"
// @Param dstChain query string false "目标链"
// @Param srcTxHash query string false "源链交易hash"
// @Param orderId query string false "跨链订单ID"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CrossData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/cross [get]
// @Security Bearer
func (e *CrossApi) GetPage(c *gin.Context) {
	s := service.Cross{}
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
