package db

import "context"

type AddOfficialEmailReq struct {
	Email   Email
	Code    string
	HubUser HubUserTO
	Context context.Context
}
