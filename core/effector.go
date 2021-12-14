package core

// Effector interface.
type Effector interface {
	Execute(eval PolicyEvalFunc, policies Policies) (allow bool, err error)
}

type PolicyEvalFunc func(policy Assertion) (allow bool, err error)
