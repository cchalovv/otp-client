package model

type GenerateRequest struct {
	// Data is some payload or identifier to uniquely generate OTP code
	Data string `bson:"data"`
}

type VerifyRequest struct {
	// Data is the same data used to generate code
	Data string `json:"data"`
	// Code is OTP code sent to phone/email
	Code string `bson:"code"`
}
