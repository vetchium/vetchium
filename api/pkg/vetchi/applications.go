package vetchi

type UpdateApplicationStateRequest struct {
	ID        string           `json:"id"         validate:"required"`
	FromState ApplicationState `json:"from_state" validate:"required"`
	ToState   ApplicationState `json:"to_state"   validate:"required"`
}
