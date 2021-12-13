package effector

import (
	"github.com/sdcxtech/casbin/core"
)

// New an deny-override	effector.
//
// The effector return allow when there is no any matched policy is deny.
// There is must be an effect column in policy.
//
// Valid effect values:
//  - allow
//  - deny
func NewDenyOverride(policyEffectKey string) core.Effector {
	return &noDenyOverrideImpl{eftKey: policyEffectKey}
}

type noDenyOverrideImpl struct {
	eftKey string
}

func (e *noDenyOverrideImpl) Execute(
	eval core.PolicyEvalFunc,
	policies core.Policies,
) (allow bool, err error) {
	allow = true
	for _, policy := range policies {
		matched, _err := eval(policy)
		if _err != nil {
			err = _err
			return
		}

		if matched {
			eft := policy[e.eftKey]
			if eft == EffectDeny {
				allow = false
				return
			}
		}
	}
	return
}
