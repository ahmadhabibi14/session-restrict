package request

type ReqSessionApprove struct {
	AccessToken string `json:"access_token" validate:"required"`
	UserId      uint64 `json:"user_id" validate:"required"`
	Role        string `json:"role" validate:"required"`
}

type ReqSessionReject struct {
	AccessToken string `json:"access_token" validate:"required"`
	UserId      uint64 `json:"user_id" validate:"required"`
	Role        string `json:"role" validate:"required"`
}
