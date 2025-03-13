package response

import "session-restrict/src/repo/sessions"

type ResGetSessionsByUserId struct {
	ResponseCommon
	Data *[]sessions.SessionsWithKey `json:"data"`
}
