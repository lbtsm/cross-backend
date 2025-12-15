package service

import (
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type Cross struct {
	service.Service
}

// GetPage 获取跨链列表
func (e *Cross) GetPage(c *dto.CrossListReq, p *actions.DataPermission, list *[]models.CrossInfo, count *int64) error {
	var err error
	var data models.CrossInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			// actions.Permission(data.TableName(), p),
		).
		Order("created_at desc").Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// // Get 获取单个跨链对象
// func (e *Cross) Get(d *dto.SingleCrossReq, model *models.CrossInfo) error {
// 	var err error

// 	db := e.Orm.First(model, d.GetId())
// 	err = db.Error
// 	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
// 		err = errors.New("查看对象不存在或无权查看")
// 		e.Log.Errorf("db error: %s", err)
// 		return err
// 	}
// 	if err = db.Error; err != nil {
// 		e.Log.Errorf("db error: %s", err)
// 		return err
// 	}
// 	return nil
// }
