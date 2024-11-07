package db

import (
	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type CreateOpeningReq struct {
	vetchi.CreateOpeningRequest

	OrgUserID  uuid.UUID
	EmployerID uuid.UUID
}

type UpdateOpeningReq struct {
	vetchi.UpdateOpeningRequest

	OrgUserID  uuid.UUID
	EmployerID uuid.UUID
}

type GetOpeningWatchersReq struct {
	vetchi.GetOpeningWatchersRequest

	OrgUserID  uuid.UUID
	EmployerID uuid.UUID
}

type AddOpeningWatchersReq struct {
	vetchi.AddOpeningWatchersRequest

	OrgUserID  uuid.UUID
	EmployerID uuid.UUID
}

type RemoveOpeningWatcherReq struct {
	OpeningID uuid.UUID
	OrgUserID uuid.UUID
}

type ApproveOpeningStateChangeReq struct {
	OpeningID uuid.UUID
	OrgUserID uuid.UUID
}

type RejectOpeningStateChangeReq struct {
	OpeningID uuid.UUID
	OrgUserID uuid.UUID
}

type GetOpeningReq struct {
	OpeningID uuid.UUID
}

type FilterOpeningsReq struct {
	vetchi.FilterOpeningsRequest

	OrgUserID uuid.UUID
}
