package entity

type Temperature struct {
	TempC float64 `json:"temp_C"`
	TempF string  `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
