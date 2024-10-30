package vetchi

type AddOrgUserRequest struct {
	Name  string        `json:"name"  validate:"required,min=1,max=255"`
	Email string        `json:"email" validate:"required,email,min=3,max=255"`
	Roles []OrgUserRole `json:"roles" validate:"required"`
}

type DisableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}
