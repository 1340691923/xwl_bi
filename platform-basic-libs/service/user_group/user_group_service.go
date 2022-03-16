package user_group

import (
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"strings"
)

type UserGroupService struct {
	ManagerID int32
	Appid     int
}

func (this *UserGroupService) AddUserGroup(userCount int, uids []string, groupRemark, groupName string) (err error) {
	b, err := util.GzipCompress(strings.Join(uids, ","))
	if err != nil {
		return
	}
	userGroup := model.UserGroup{}
	userGroup.GroupRemark = groupRemark
	userGroup.GroupName = groupName
	return userGroup.Insert(this.ManagerID, this.Appid, userCount, b)
}

func (this *UserGroupService) ModifyUserGroup(id int, groupRemark, groupName string) (err error) {
	userGroup := model.UserGroup{}
	userGroup.Id = id
	userGroup.GroupName = groupName
	userGroup.GroupRemark = groupRemark
	return userGroup.ModifyUserGroup(this.ManagerID, this.Appid)
}

func (this *UserGroupService) DeleteUserGroup(id int) (err error) {
	userGroup := model.UserGroup{}
	userGroup.Id = id
	return userGroup.DeleteUserGroupById(this.ManagerID, this.Appid)
}

func (this *UserGroupService) UserGroupList() (list []model.UserGroup, err error) {
	userGroup := model.UserGroup{}
	list, err = userGroup.List(this.ManagerID, this.Appid)
	if err != nil {
		return nil, err
	}

	for index := range list {
		userListData, err := util.GzipUnCompress(list[index].UserList)
		if err != nil {
			return nil, err
		}
		list[index].UserListData = strings.Split(userListData, ",")
	}

	return list, err

}

func (this *UserGroupService) Options() (list []model.UserGroup, err error) {
	userGroup := model.UserGroup{}
	list, err = userGroup.GetSelectOptions(this.ManagerID, this.Appid)
	return
}
