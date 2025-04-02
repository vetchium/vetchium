package employer

type ChangeCoolOffPeriodRequest struct {
	CoolOffPeriodDays int32 `json:"cool_off_period_days" validate:"min=0,max=365"`
}
