package response

const (
	ResponseStatusError        Status = "error"
	ResponseStatusUnauthorized Status = "unauthorized"
	ResponseStatusForbidden    Status = "forbidden"
	ResponseStatusNotFound     Status = "not found"
	ResponseStatusSuccess      Status = "success"
	ResponseStatusnknown       Status = "uknown"
)

var (
	ResponseErrorBadRequest    Message = "bad request"
	ResponseErrorGeneral       Message = "general error"
	ResponseErrorDatabase      Message = "database error"
	ResponseErrorInteral       Message = "internal server error"
	ResponseErrorUnknown       Message = "something went wrong"
	ResponseErrorToManyRequest Message = "Too many requests. Try again later."
)

const (
	ResponseMessageSuccess Message = "success"
)
