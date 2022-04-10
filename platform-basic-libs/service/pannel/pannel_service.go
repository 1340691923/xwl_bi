package pannel

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/model"
	"github.com/1340691923/xwl_bi/platform-basic-libs/request"
	"strconv"
)

type PannelService struct {
}

func (this *PannelService) UpdatePannelManager(managers string, id int, managerUid int32) (err error) {
	_, err = db.
		SqlBuilder.
		Update("pannel").
		SetMap(map[string]interface{}{"managers": managers}).
		Where(db.Eq{"id": id, "create_by": managerUid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *PannelService) RtListByAppid(appid int) (res map[int]model.ReportTable, err error) {

	res = map[int]model.ReportTable{}

	sql, args, err := db.
		SqlBuilder.
		Select("id,name,rt_type,data").
		From("report_table").
		Where(db.Eq{"appid": appid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Sqlx.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rtid int
		var tmp model.ReportTable
		err := rows.Scan(&rtid, &tmp.Name, &tmp.RtType, &tmp.Data)
		if err != nil {
			return nil, err
		}
		res[rtid] = tmp
	}

	return res, nil
}

func (this *PannelService) UpdatePannelRt(reportTables string, id int, managerUid int32) (err error) {
	_, err = db.SqlBuilder.
		Update("pannel").
		SetMap(map[string]interface{}{
			"report_tables": reportTables,
		}).
		Where(db.Eq{"id": id, "create_by": managerUid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *PannelService) CopyPannel(data model.Pannel, managerUid int32) (err error) {
	_, err = db.
		Sqlx.
		Exec(`insert into pannel(folder_id,pannel_name,create_by,report_tables) values (?,?,?,?)`,
		data.FolderId, data.PannelName, managerUid, data.ReportTables,
		)
	return
}

func (this *PannelService) DeleteDir(data model.PannelFolder) (err error) {
	_, err = db.Sqlx.
		Exec(`		
		DELETE p1,p2 FROM
	pannel_folder p1
	LEFT JOIN pannel p2 ON p1.id = p2.folder_id where p1.id = ? and p1.appid = ? and p1.create_by = ?`, data.Id, data.Appid, data.CreateBy)
	return
}

func (this *PannelService) DeletePannel(id int, managerUid int32) (err error) {
	_, err = db.SqlBuilder.
		Delete("pannel").
		Where(db.Eq{"id": id, "create_by": managerUid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *PannelService) MovePannel2Dir(id, folderId int, managerUid int32) (err error) {
	_, err = db.SqlBuilder.
		Update("pannel").
		SetMap(map[string]interface{}{"folder_id": folderId}).
		Where(db.Eq{"id": id, "create_by": managerUid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *PannelService) Rename(pannelName string, id int, managerUid int32) (err error) {
	_, err = db.SqlBuilder.
		Update("pannel").
		SetMap(map[string]interface{}{"pannel_name": pannelName}).
		Where(db.Eq{"id": id, "create_by": managerUid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *PannelService) AddDir(data request.NewDir) (err error) {
	_, err = db.
		SqlBuilder.
		Insert("pannel_folder").
		SetMap(map[string]interface{}{
			"appid":       data.Appid,
			"create_by":   data.CreateBy,
			"folder_name": data.FolderName,
		}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *PannelService) AddPannel(data request.NewPannel, managerUid int32) (err error) {
	_, err = db.SqlBuilder.
		Insert("pannel").
		SetMap(map[string]interface{}{
			"folder_id":   data.FolderId,
			"pannel_name": data.PannelName,
			"create_by":   managerUid,
		}).
		RunWith(db.Sqlx).
		Exec()
	return
}

type Pannel struct {
	Folder_name  string `json:"folder_name" db:"folder_name"`
	FolderType   int    `json:"folder_type" db:"folder_type"`
	FolderId     int    `json:"folder_id" db:"folder_id"`
	PannelId     int    `json:"pannel_id" db:"pannel_id"`
	PannelName   string `json:"pannel_name" db:"pannel_name"`
	Managers     string `json:"managers" db:"managers"`
	ReportTables string `json:"report_tables" db:"report_tables"`
	CreateBy     int    `json:"create_by" db:"create_by"`
}

func (this *PannelService) GetPannelList(managerUid int32, appid int) (res []Pannel, err error) {
	sql := `
			SELECT
	p1.folder_name,
	p1.id AS folder_id,
	ifnull(p2.report_tables,'') as report_tables ,
	 ifnull(p2.id,0)as pannel_id,
	ifnull(p2.pannel_name,'') as pannel_name,ifnull(p1.create_by,0) as create_by,ifnull(FIND_IN_SET(?,p2.managers),0) as folder_type,ifnull(p2.managers,'') as managers
FROM
	pannel_folder p1
	LEFT JOIN pannel p2 ON p1.id = p2.folder_id 
	
			WHERE
				(p1.create_by = ? OR FIND_IN_SET(?,p2.managers)  )   and appid= ?`

	err = db.
		Sqlx.
		Select(&res, sql, managerUid, managerUid, managerUid, appid)

	if err != nil {
		return
	}

	var userModel model.GmUserModel
	list, err := userModel.Select(strconv.Itoa(appid))
	if err != nil {
		return
	}

	userMap := map[int]string{}

	for _, v := range list {
		userMap[int(v.ID)] = v.Realname
	}

	for k, v := range res {
		if v.FolderType != 0 {
			if _, ok := userMap[v.CreateBy]; ok {
				res[k].PannelName = fmt.Sprintf("%s(%s)", res[k].PannelName, userMap[v.CreateBy])
				res[k].Folder_name = "共享文件夹"
				res[k].FolderId = -1
			}
		}
	}

	return
}

func (this *PannelService) FindRtById(id int, managerUid int32) (reportTable model.ReportTable, err error) {
	sqls, args, err := db.SqlBuilder.
		Select("*").
		From("report_table").
		Where(db.Eq{"id": id, "user_id": managerUid}).
		ToSql()

	if err != nil {
		return reportTable, err
	}

	if err := db.Sqlx.Get(&reportTable, sqls, args...); err != nil {
		if err == sql.ErrNoRows {
			return reportTable, errors.New("您无权限操作该报表或该报表已被删除")
		}
		return reportTable, err
	}

	return reportTable, err
}

func (this *PannelService) FindNameCount(data request.FindNameCount, managerUid int32) (count int, err error) {
	if err := db.SqlBuilder.
		Select("count(1)").
		From("report_table").
		Where(db.Eq{"name": data.Name, "rt_type": data.RtType, "appid": data.Appid, "user_id": managerUid}).
		RunWith(db.Sqlx).
		QueryRow().
		Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (this *PannelService) DeleteReportTableByID(data model.ReportTable, managerUid int32) (err error) {
	if _, err := db.
		SqlBuilder.
		Delete("report_table").
		Where(db.Eq{"id": data.Id, "user_id": managerUid}).
		RunWith(db.Sqlx).Exec();err != nil{
		return err
	}

	if _, err := db.
		Sqlx.
		Exec(`UPDATE pannel SET report_tables = TRIM(BOTH ',' FROM REPLACE(CONCAT(',', report_tables, ','), concat(',',?,','), ','))
						 WHERE FIND_IN_SET(?, report_tables) and create_by = ?`, data.Id, data.Id, managerUid); err != nil {
		return err
	}

	return nil
}

func (this *PannelService) ReportTableList(appid int, rtType int8, managerUid int32) (list []model.ReportTable, err error) {

	where := db.Eq{
		"appid":   appid,
		"user_id": managerUid,
	}

	if rtType != 0 {
		where["rt_type"] = rtType
	}

	sql, args, err := db.
		SqlBuilder.
		Select("*").
		From("report_table").
		Where(where).
		OrderBy("rt_type").
		ToSql()

	if err != nil {
		return nil, err
	}

	if err := db.Sqlx.Select(&list, sql, args...); err != nil {
		return nil, err
	}

	return list, nil

}
