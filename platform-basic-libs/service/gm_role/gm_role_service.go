//GM角色层
package gm_role

import (
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
)

// GmRoleService
type GmRoleService struct {
}

func (this GmRoleService) Select() (list []model.GmRoleModel, err error) {
	var roleModel model.GmRoleModel
	list, err = roleModel.Select()
	if err != nil {
		return
	}
	return
}

func (this GmRoleService) Add(model2 model.GmRoleModel) (id int64, err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	id, err = roleModel.Insert()
	return
}

func (this GmRoleService) Update(model2 model.GmRoleModel) (err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	err = roleModel.Update()
	return
}

func (this GmRoleService) Delete(id int) (err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = id
	err = roleModel.Delete()
	return
}

func (this GmRoleService) GetRoles(roles []model.GmRoleModel) (list []request.GmRoleModel, err error) {
	for _, v := range roles {
		roleRes := request.GmRoleModel{
			ID:          v.ID,
			RoleName:    v.RoleName,
			Description: v.Description,
			RoleList:    v.RoleList,
		}
		apis := []string{}

		rows, err := db.Sqlx.Query("select v1 from casbin_rule where v0 = ?;", v.ID)
		if util.FilterMysqlNilErr(err) {
			logs.Logger.Sugar().Errorf("err:", err)
			continue
		}
		defer rows.Close()
		for rows.Next() {
			api := ""
			err := rows.Scan(&api)
			if err != nil {
				logs.Logger.Sugar().Errorf("err:", err)
				continue
			}
			apis = append(apis, api)
		}
		roleRes.Api = apis
		list = append(list, roleRes)
	}
	return
}
