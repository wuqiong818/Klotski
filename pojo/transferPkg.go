package pojo

type RequestPkg struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ResponsePkg struct {
	Type string `json:"type"`
	Code int    `json:"code"`
	Data string `json:"data"`
}

type CertificationInfo struct {
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	UserProfile string `json:"userProfile"`
}

var (
	CertificationType      = "CertificationType"
	CreateRoomType         = "CreateRoomType"
	JoinRoomType           = "JoinRoomType"
	ExchangeMutualInfoType = "ExchangeMutualInfoType"
	RefreshScoreType       = "RefreshScoreType"
	DiscontinueQuitType    = "DiscontinueQuitType"
	GameOverType           = "GameOverType"
	PingCheckType          = "PingCheckType"
)

// Success 状态码定义集合
var (
	Success                = 2000
	UserCertificationError = 5000
	CreateRoomError        = 5001
	JoinRoomError          = 5002
	RefreshScoreError      = 5003
)
