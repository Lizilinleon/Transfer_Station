package util

type BillingInput struct {
	InputTokens      int64
	OutputTokens     int64
	RequestCount     int64
	InputTokenPrice  float64
	OutputTokenPrice float64
	RequestPrice     float64
}

func CalculateQuotaCharge(input BillingInput) float64 {
	inputCost := float64(input.InputTokens) * input.InputTokenPrice
	outputCost := float64(input.OutputTokens) * input.OutputTokenPrice
	requestCost := float64(input.RequestCount) * input.RequestPrice
	return inputCost + outputCost + requestCost
}
