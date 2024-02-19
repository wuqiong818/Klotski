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

type GameData struct {
	Score     string `json:"score"`
	StepCount string `json:"stepCount"`
	TakeTime  string `json:"takeTime"`
}

var (
	CertificationType   = "CertificationType"
	CreateRoomType      = "CreateRoomType"
	JoinRoomType        = "JoinRoomType"
	RefreshScoreType    = "RefreshScoreType"
	DiscontinueQuitType = "DiscontinueQuitType"
	GameOverType        = "GameOverType"
	HeartCheckType      = "HeartCheckType"
)

// Success 状态码定义集合
var (
	Success                = 2000
	UserCertificationError = 5000
	CreateRoomError        = 5001
	JoinRoomError          = 5002
	RefreshScoreError      = 5003
)
