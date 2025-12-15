package service

import (
	"errors"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
)

type CheckTss struct {
	service.Service
}

// GetPage 获取tss健康检查列表
func (e *CheckTss) GetPage(c *dto.CrossListReq, p *actions.DataPermission, list *[]models.CrossInfo, count *int64) error {
	var err error
	var data models.CrossInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			// actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).Order("created_at desc").
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Get 获取单个跨链对象
func (e *Cross) Get(d *dto.SingleCrossReq, model *models.CrossInfo) error {
	var err error

	db := e.Orm.First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err = db.Error; err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}
