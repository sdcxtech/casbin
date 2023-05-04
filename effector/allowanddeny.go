package effector

import (
	"github.com/sdcxtech/casbin/core"
)

// New an allow-and-deny effector.
//
// The effector return allow when any matched policy is allow,
// and there is no any matched policy is deny.
//
// There is must be an effect column in policy.
//
// Valid effect values:
//   - allow
//   - deny
func NewAllowAndDeny(effectKey string) core.Effector {
	return &allowAndDeny{eftKey: effectKey}
}

type allowAndDeny struct {
	eftKey string
}

func (a *allowAndDeny) Execute(
	eval core.PolicyEvalFunc,
	policies core.Policies,
) (allow bool, err error) {
	for _, policy := range policies {
		matched, _err := eval(policy)
		if _err != nil {
			err = _err

			return
		}

		if matched {
			eft := policy[a.eftKey]
			if eft == EffectDeny {
				allow = false
			} else if eft == EffectAllow {
				allow = true
			}
		}
	}

	return
}
