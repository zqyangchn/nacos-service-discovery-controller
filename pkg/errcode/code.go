package errcode

var (
	// Success 成功
	Success = New("00000", "Success")

	// ClientError A 组
	// 客户端错误
	ClientError        = New("A0001", "用户端错误")
	InvalidParamsError = New("A0002", "请求参数错误")

	// TokenAuthError token
	TokenAuthError = New("A0101", "鉴权失败")

	// ServerError B 组
	// 服务端错误
	ServerError = New("B0001", "系统执行出错")

	// ThirdPartyCallError C 组
	// 第三方调用错误
	ThirdPartyCallError = New("C0001", "第三方调用错误")
)
