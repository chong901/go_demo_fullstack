package constants

const (
	SessionRole       = "role"
	SessionLoginUser  = "loginUser"
	SessionCookieLang = "cookieLang"
	SessionHeaderLang = "headerLang"

	ContextSetLang  = "contextLang"
	DefaultLang     = "en-US"
	TemplateLangStr = "lang"

	SocketRoomMosi = "mosi"

	SocketEventPoUpdate           = "poUpdate"
	SocketEventJobUpdate          = "jobUpdate"
	SocketEventInventoryUpdate    = "inventoryUpdate"
	SocketEventReportOutputUpdate = "reportOutputUpdate"

	FromBrowser = "browser"
	FromPi      = "Pi"

	UserFieldStr = "User"

	RequestQueryIdString = "id"

	BindJson  = 1
	BindForm  = 2
	BindQuery = 3
)
