package response

type (
	Status  string
	Message string
)

type ValidationError struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	Tag      string `json:"tag"`
	TagValue string `json:"tag_value"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Status  Status      `json:"status"`
	Message Message     `json:"message"`
	Code    int         `json:"code"`
}
