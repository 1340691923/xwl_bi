package request

// GmRoleModel
type GmRoleModel struct {
	ID          int      `json:"id" db:"id"`
	RoleName    string   `json:"name" db:"role_name"`
	Description string   `json:"description" db:"description"`
	RoleList    string   `json:"routes" db:"role_list"`
	Api         []string `json:"api"`
}

type AnalysisFilter struct {
	FilterType string `json:"filterType"`
	Filts      []struct {
		FilterType string `json:"filterType"`
		Filts      []struct {
			ColumnName string      `json:"columnName"`
			Comparator string      `json:"comparator"`
			FilterType string      `json:"filterType"`
			Ftv        interface{} `json:"ftv"`
		} `json:"filts,omitempty"`
		Relation   string      `json:"relation,omitempty"`
		ColumnName string      `json:"columnName,omitempty"`
		Comparator string      `json:"comparator,omitempty"`
		Ftv        interface{} `json:"ftv,omitempty"`
	} `json:"filts"`
	Relation string `json:"relation"`
}

type Zhibiao struct {
	EventName        string         `json:"eventName"`
	EventNameDisplay string         `json:"eventNameDisplay"`
	Relation         AnalysisFilter `json:"relation"`
}

type FunnelReqData struct {
	UserGroup         []int          `json:"userGroup"`
	ZhibiaoArr        []Zhibiao      `json:"zhibiaoArr"`
	WhereFilter       AnalysisFilter `json:"whereFilter"`
	WindowTime        int            `json:"windowTime"`
	WindowTimeFormat  string         `json:"windowTimeFormat"`
	Date              []string       `json:"date"`
	Appid             int            `json:"appid"`
	WhereFilterByUser AnalysisFilter `json:"whereFilterByUser"`
	GroupBy           []string       `json:"groupBy"`
}

type TraceReqData struct {
	EventNames        []string       `json:"eventNames"`
	UserGroup         []int          `json:"userGroup"`
	ZhibiaoArr        []Zhibiao      `json:"zhibiaoArr"`
	WhereFilter       AnalysisFilter `json:"whereFilter"`
	WindowTime        int            `json:"windowTime"`
	WindowTimeFormat  string         `json:"windowTimeFormat"`
	Date              []string       `json:"date"`
	Appid             int            `json:"appid"`
	WhereFilterByUser AnalysisFilter `json:"whereFilterByUser"`
	GroupBy           []string       `json:"groupBy"`
}

type RetentionReqData struct {
	UserGroup  []int `json:"userGroup"`
	ZhibiaoArr []struct {
		EventName        string         `json:"eventName"`
		EventNameDisplay string         `json:"eventNameDisplay"`
		Relation         AnalysisFilter `json:"relation"`
	} `json:"zhibiaoArr"`
	WhereFilter       AnalysisFilter `json:"whereFilter"`
	WindowTime        int            `json:"windowTime"`
	WindowTimeFormat  string         `json:"windowTimeFormat"`
	Date              []string       `json:"date"`
	Appid             int            `json:"appid"`
	WhereFilterByUser AnalysisFilter `json:"whereFilterByUser"`
	GroupBy           []string       `json:"groupBy"`
}

type FormulaDimension struct {
	SelectAttr []string       `json:"selectAttr"`
	EventName  string         `json:"eventName"`
	Relation   AnalysisFilter `json:"relation"`
}

type EventZhibiao struct {
	SelectAttr        []string         `json:"selectAttr,omitempty"`
	Typ               int              `json:"typ"`
	EventName         string           `json:"eventName,omitempty"`
	EventNameDisplay  string           `json:"eventNameDisplay"`
	Relation          AnalysisFilter   `json:"relation,omitempty"`
	ScaleType         string           `json:"scaleType,omitempty"`
	Operate           string           `json:"operate,omitempty"`
	One               FormulaDimension `json:"one,omitempty"`
	Two               FormulaDimension `json:"two,omitempty"`
	DivisorNoGrouping bool             `json:"divisor_no_grouping"`
}

type EventReqData struct {
	UserGroup         []int          `json:"userGroup"`
	ZhibiaoArr        []EventZhibiao `json:"zhibiaoArr"`
	GroupBy           []string       `json:"groupBy"`
	WhereFilter       AnalysisFilter `json:"whereFilter"`
	WhereFilterByUser AnalysisFilter `json:"whereFilterByUser"`
	Date              []string       `json:"date"`
	WindowTimeFormat  string         `json:"windowTimeFormat"`
	Appid             int            `json:"appid"`
}

type UserAttrReqData struct {
	UserGroup         []int          `json:"userGroup"`
	ZhibiaoArr        []string       `json:"zhibiaoArr"`
	GroupBy           []string       `json:"groupBy"`
	WhereFilterByUser AnalysisFilter `json:"whereFilterByUser"`
	Appid             int            `json:"appid"`
}

type UserListReqData struct {
	UI    []string `json:"ui"`
	Appid int      `json:"appid"`
}

type NewPannel struct {
	PannelName string `json:"pannel_name"`
	FolderId   int    `json:"folder_id"`
}

type NewDir struct {
	FolderName string `db:"folder_name" json:"folder_name"`
	FolderType int8   `db:"folder_type" json:"folder_type"` //0为自己创建的
	CreateBy   int    `db:"create_by" json:"create_by"`
	Appid      int    `db:"appid" json:"appid"`
}

type FindRtById struct {
	Appid int `db:"appid" json:"appid"`
	Id    int `json:"id"`
}

type FindNameCount struct {
	Appid  int    `db:"appid" json:"appid"`
	Name   string `db:"name" json:"name"`
	RtType int8   `db:"rt_type" json:"rt_type"`
}

type GetPannelList struct {
	Appid int `db:"appid" json:"appid"`
}

type AddUserGroup struct {
	Ids    []string `json:"uids"`
	Name   string   `json:"name"`
	Remark string   `json:"remark"`
	Appid  int      `json:"appid"`
}

type ModifyUserGroup struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
	Appid  int    `json:"appid"`
}

type DeleteUserGroup struct {
	Id    int `json:"id"`
	Appid int `json:"appid"`
}

type UserGroupList struct {
	Appid int `json:"appid"`
}

type UserEventDetailReq struct {
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	Appid      int      `json:"appid"`
	UserID     string   `json:"userId"`
	EventName  string   `json:"eventName"`
	OrderBy    string   `json:"orderBy"`
	Date       []string `json:"date"`
	EventNames []string `json:"eventNames"`
}

type UserEventListReq struct {
	Uid   int `json:"uid"`
	Appid int `json:"appid"`
}

type UserEventCountReq struct {
	Appid            int      `json:"appid"`
	WindowTimeFormat string   `json:"windowTimeFormat"`
	UserID           string   `json:"userId"`
	EventNames       []string `json:"eventNames"`
	Date             []string `json:"date"`
}

type LoadPropQuotasReq struct {
	EventName string `json:"event_name"`
	Appid     int    `json:"appid"`
}

type RolesDelReq struct {
	Id int `json:"id"`
}

type UserListReq struct {
	Appid int `json:"appid,omitempty" `
}

type DeleteUserReq struct {
	Id int32 `json:"id"`
}

type GetUserByIdReq struct {
	Id int32 `json:"id"`
}

type ReportCountReq struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Appid     int    `json:"appid"`
}

type EventFailDescReq struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Appid     int    `json:"appid"`
	DataName  string `json:"data_name"`
}

type AddDebugDeviceIDReq struct {
	Appid    int    `json:"appid"`
	Remark   string `json:"remark"`
	DeviceID string `json:"deviceID"`
}

type DelDebugDeviceIDReq struct {
	Appid    int    `json:"appid"`
	DeviceID string `json:"deviceID"`
}

type DebugDeviceIDListReq struct {
	Appid int `json:"appid"`
}
type UserUpdateReq struct {
	Id       int    `json:"id"`
	Realname string `json:"realname"`
	RoleId   int32  `json:"role_id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserAddReq struct {
	Realname string `json:"realname"`
	RoleId   int32  `json:"role_id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserBanReq struct {
	Typ int `json:"typ"`
	Id  int `json:"id"`
}

type AttrManagerReq struct {
	Appid int `json:"appid"`
	Typ   int `json:"typ"`
}

type UpdateAttrInvisibleReq struct {
	Appid           int    `json:"appid"`
	AttributeSource int    `json:"attribute_source"`
	AttributeName   string `json:"attribute_name"`
	Status          int    `json:"status"`
}

type AttrManagerByMetaReq struct {
	Appid     int    `json:"appid"`
	Typ       int    `json:"typ"`
	EventName string `json:"event_name"`
}

type UpdateShowNameReq struct {
	Appid     int    `json:"appid"`
	EventName string `json:"event_name"`
	ShowName  string `json:"show_name"`
}

type UpdateAttrShowNameReq struct {
	Appid         int    `json:"appid"`
	AttributeName string `json:"attribute_name"`
	Typ           int    `json:"typ"`
	ShowName      string `json:"show_name"`
}

type GetCalcuSymbolDataReq struct {
	Appid     int    `json:"appid"`
	EventName string `json:"event_name"`
}

type GetAnalyseSelectOptionsReq struct {
	Appid int `json:"appid"`
}

type GmOperaterLogList struct {
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
	UserId         int      `json:"operater_id"`
	RoleId         int      `json:"operater_role_id"`
	OperaterAction string   `json:"operater_action"`
	Date           []string `json:"date"`
}
