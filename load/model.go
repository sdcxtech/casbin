package load

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/sdcxtech/casbin/core"
	"github.com/sdcxtech/casbin/effector"
)

const (
	EffectorAllowOverride = "allow-override"
	EffectorDenyOverride  = "deny-override"
	EffectorAllowAndDeny  = "allow-and-deny"
)

// ModelFromViper load model definetion from a viper instance.
//
// You can use viper to load a concrete configure file in JSON, TOML, YAML format.
//
// TOML format example:
//
//      [request_definition]
//      r = "sub, dom, obj, act"
//
//      [policy_definition]
//      p = "sub, obj, act"
//
//      [role_definition]
//      g =  "_, _, _"
//      g1 = "_, _"
//
//      [policy_effect]
//      type = "allow-override"
//
//      [matchers]
//      m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
func ModelFromViper(v *viper.Viper, options ...ModelOption) (model *core.Model, err error) {
	optionConf, err := newLoadModelConfig(options...)
	if err != nil {
		return
	}

	request, err := core.NewAssertionSchema(v.GetString("request_definition.r"))
	if err != nil {
		err = fmt.Errorf("request_definition.r: %w", err)

		return
	}

	policy, err := core.NewAssertionSchema(v.GetString("policy_definition.p"))
	if err != nil {
		err = fmt.Errorf("policy_definition.p: %w", err)

		return
	}

	rolesSchema := make(core.RolesSchema)

	roleDefinition := v.GetStringMapString("role_definition")
	for key, v := range roleDefinition {
		roleType, _err := core.RoleTypeFromLine(v)
		if _err != nil {
			err = _err

			return
		}

		var match core.RoleDomainMatchFunc
		match, ok := optionConf.roleDomainMatchFuncs[key]

		if !ok {
			match = core.RoleDomainMatchEqual
		}

		rolesSchema[key] = core.RoleSchema{
			Type:            roleType,
			DomainMatchFunc: match,
		}
	}

	eftType := v.GetString("policy_effect.type")
	eftKey := v.GetString("policy_effect.key")

	var eft core.Effector

	switch eftType {
	case EffectorAllowOverride:
		eft = effector.NewAllowOverride()
	case EffectorDenyOverride:
		if eftKey == "" || !policy.Has(eftKey) {
			err = fmt.Errorf("%w: %s", ErrNeedEffectFieldKey, eftType)

			return
		}

		eft = effector.NewDenyOverride(eftKey)
	case EffectorAllowAndDeny:
		if eftKey == "" || !policy.Has(eftKey) {
			err = fmt.Errorf("%w: %s", ErrNeedEffectFieldKey, eftType)

			return
		}

		eft = effector.NewAllowAndDeny(eftKey)
	default:
		err = fmt.Errorf("%w: %s", ErrUnknownEffectType, eftType)
	}

	if err != nil {
		return
	}

	mattchersDefine := v.GetStringMapString("matchers")

	matchers, err := core.MatchersConfig{
		Roles:          rolesSchema,
		Define:         mattchersDefine,
		ExtensionFuncs: optionConf.extensionFuncs,
	}.New()
	if err != nil {
		err = fmt.Errorf("new matchers: %w", err)

		return
	}

	model = core.NewModel(policy, request, rolesSchema, eft, matchers)

	return model, err
}
