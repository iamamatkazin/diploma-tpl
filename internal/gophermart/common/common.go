package common

func PtFloatToInt(val *float64) *int {
	if val != nil {
		i := int(*val * 100)
		return &i
	}

	return nil
}
