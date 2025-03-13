package response

import "session-restrict/src/repo/sessions"

type ResGetSessionsByUserId struct {
	ResponseCommon
	Data *[]sessions.Session `json:"data"`
}

type ResSessionApprove struct {
	ResponseCommon
	Data *[]sessions.Session `json:"data"`
}

type ResSessionReject struct {
	ResponseCommon
	Data *[]sessions.Session `json:"data"`
}
