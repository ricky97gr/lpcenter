package response

const (
	SuccessCode = 200

	ErrStruct         = 100001
	ErrDB             = 100002
	ErrAuth           = 100003
	ErrRecordNotFound = 100004
	ErrUserExist      = 100005
	ErrPluginExist    = 100006
	ErrLicenseExist   = 100007
)

var errCodeMap = map[int32]struct {
	msgCn string
	msgEn string
}{
	SuccessCode:         {msgCn: "操作成功", msgEn: "handle successfully"},
	ErrStruct:           {msgCn: "请求参数错误", msgEn: "error struct"},
	ErrDB:               {msgCn: "数据库错误", msgEn: "error db"},
	ErrAuth:             {msgCn: "认证失败，请重新登录", msgEn: "err auth info, please login again"},
	ErrRecordNotFound:   {msgCn: "查询记录不存在", msgEn: "the record in db is not found"},
	ErrUserExist:        {msgCn: "用户已存在", msgEn: "user already exists"},
	ErrPluginExist:      {msgCn: "插件已存在", msgEn: "plugin already exists"},
	ErrLicenseExist:     {msgCn: "授权已存在", msgEn: "license already exists"},
}

func GetErrMsg(code int32) string {
	if info, ok := errCodeMap[code]; ok {
		return info.msgCn
	}
	return "未知错误"
}
