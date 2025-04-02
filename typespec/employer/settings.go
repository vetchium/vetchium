package employer

type ChangeCoolOffPeriodRequest struct {
	CoolOffPeriod int32 `json:"cool_off_period" validate:"min=0,max=365"`
}
