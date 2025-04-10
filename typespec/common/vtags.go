package common

type FilterVTagsRequest struct {
	Prefix *string `json:"prefix,omitempty"`
}

type VTagID string

type VTagName string

type VTag struct {
	ID   VTagID   `json:"id"   validate:"required,uuid"`
	Name VTagName `json:"name" validate:"required,max=32"`
}

type VTags struct {
	Tags []VTag `json:"tags"`
}
