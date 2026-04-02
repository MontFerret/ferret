package core

type ErrorPolicy int

const (
	ErrorPolicyDefault ErrorPolicy = iota
	ErrorPolicySuppress
	ErrorPolicyFail
)
