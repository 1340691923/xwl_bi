//BI用户层
package gm_user

import (
	"errors"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/jwt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"go.uber.org/zap"
	"time"
)

// GmUserService
type GmUserService struct {
}

func (this GmUserService) CheckLogin(username, password string) (token string, err error) {
	var model2 model.GmUserModel
	model2.Password = password
	model2.Username = username
	gmUser, err := model2.GetUserByUP()

	if err != nil {
		logs.Logger.Error("登陆失败", zap.Error(err))
		err = errors.New("用户验证失败")
		return
	}

	if gmUser.IsDel == 1 {
		err = errors.New("您的账号已被封禁")
		return
	}

	db.SqlBuilder.
		Update("gm_user").
		SetMap(map[string]interface{}{"last_login_time": time.Now().Format(util.TimeFormat)}).
		Where(db.Eq{"id": gmUser.ID}).
		RunWith(db.Sqlx).
		Exec()

	token, err = jwt.GenerateToken(gmUser)
	if err != nil {
		return
	}
	return
}

func (this GmUserService) GetRoleInfo(roleId int32) (gminfo model.GmRoleModel, err error) {
	var model2 model.GmRoleModel
	gminfo, err = model2.GetById(int(roleId))
	if err != nil {
		return
	}
	return
}

func (this GmUserService) IsExitUser(claims *jwt.Claims) bool {
	var model2 model.GmUserModel
	model2.Username = claims.Username
	model2.RoleId = claims.RoleId
	return model2.Exsit()
}
