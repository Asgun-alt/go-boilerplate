package constant

type (
	Error   error
	Message string
	LogTag  string
)

const (
	ErrorToManyRequest    Message = "to many request"
	ErrorForbidden        Message = "forbidden"
	ErrorDupCheck         Message = "duplicate check"
	ErrorUnauthorized     Message = "unauthorized"
	ErrorRequest          Message = "bad request"
	ErrorGeneral          Message = "general error"
	ErrorDatabase         Message = "database error"
	ErrorInteral          Message = "internal server error"
	ErrorNotFound         Message = "not found"
	ErrorPanic            Message = "panic"
	ErrorMethodNotAllowed Message = "method not allowed"
)

const (
	SessionKey = "sess_key"
)
