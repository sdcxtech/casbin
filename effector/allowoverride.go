package effector

import (
	"github.com/sdcxtech/casbin/core"
)

// New an allow-override effector.
//
// The effector return allow when any matched policy is allow.
// There is no effect column, policy is default to be allow effect.
func NewAllowOverride() core.Effector {
	return &allowOverride{}
}

type allowOverride struct{}

func (a *allowOverride) Execute(
	eval core.PolicyEvalFunc, policies core.Policies,
) (allow bool, err error) {
	for _, policy := range policies {
		allow, err = eval(policy)
		if err != nil {
			return
		}

		if allow {
			return
		}
	}

	return
}
