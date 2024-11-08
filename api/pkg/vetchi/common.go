package vetchi

type CountryCode string

type TimeZone string

type Currency string

type HubUserShort struct {
	FullName string `json:"full_name"`
	Handle   string `json:"handle"`
}

type OrgUserShort struct {
	OrgUserName  string `json:"org_user_name"`
	OrgUserEmail string `json:"org_user_email"`
}
