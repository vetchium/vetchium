package middleware

type userCtx int

const (
	OrgUserCtxKey userCtx = iota
	HubUserCtxKey
)
