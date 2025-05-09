package employer

type CostCenterName string

type CostCenterState string

const (
	ActiveCC  CostCenterState = "ACTIVE_CC"
	DefunctCC CostCenterState = "DEFUNCT_CC"
)

type CostCenter struct {
	Name  CostCenterName  `json:"name"            validate:"required,min=3,max=64" db:"cost_center_name"`
	Notes string          `json:"notes,omitempty" validate:"max=1024"              db:"notes"`
	State CostCenterState `json:"state"                                            db:"cost_center_state"`
}

type AddCostCenterRequest struct {
	Name  CostCenterName `json:"name"            validate:"required,min=3,max=64"`
	Notes string         `json:"notes,omitempty" validate:"max=1024"`
}

type DefunctCostCenterRequest struct {
	Name CostCenterName `json:"name" validate:"required,min=3,max=64"`
}

type GetCostCenterRequest struct {
	Name CostCenterName `json:"name" validate:"required,min=3,max=64"`
}

type GetCostCentersRequest struct {
	Limit         int               `json:"limit,omitempty"          validate:"max=100"`
	PaginationKey CostCenterName    `json:"pagination_key,omitempty"`
	States        []CostCenterState `json:"states,omitempty"         validate:"validate_cc_states"`
}

func (g *GetCostCentersRequest) StatesAsStrings() []string {
	if len(g.States) == 0 {
		return []string{string(ActiveCC)}
	}

	states := []string{}
	for _, state := range g.States {
		// already validated by vator
		states = append(states, string(state))
	}
	return states
}

type RenameCostCenterRequest struct {
	OldName CostCenterName `json:"old_name" validate:"required,min=3,max=64"`
	NewName CostCenterName `json:"new_name" validate:"required,min=3,max=64"`
}

type UpdateCostCenterRequest struct {
	Name  CostCenterName `json:"name"  validate:"required,min=3,max=64"`
	Notes string         `json:"notes" validate:"required,max=1024"`
}
