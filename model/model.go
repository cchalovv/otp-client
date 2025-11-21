package model

type CreateRequest struct {
	// Data is some payload or identifier to uniquely generate OTP code
	Data string
}

type VerifyRequest struct {
	// Data is the same data used to generate code
	Data string
	// Code is OTP code sent to phone/email
	Code string
}
