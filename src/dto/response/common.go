package response

type ResponseCommon struct {
	StatusCode int    `json:"status"`            // HTTP status code
	Message    string `json:"message,omitempty"` // Response message
	Error      string `json:"error,omitempty"`   // Error message (if any)
} // @name ResponseCommon

func (resp *ResponseCommon) SetStatus(code int) {
	resp.StatusCode = code
}

func (resp *ResponseCommon) SetMessage(msg string) {
	resp.Message = msg
}

func (resp *ResponseCommon) SetError(errStr string) {
	resp.Error = errStr
}
