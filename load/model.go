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

// LoadModel load model definetion from a viper instance.
//
// You can use viper to load a concrete configure file in JSON, TOML, YAML format.
//
// TOML format example:
//
// [request_definition]
// r = "sub, dom, obj, act"

// [policy_definition]
// p = "sub, obj, act"

// [role_definition]
// g =  "_, _, _"
// g1 = "_, _"

// [policy_effect]
// type = "allow-override"

// [matchers]
// m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
//
func LoadModel(v *viper.Viper, options ...LoadModelOption) (model *core.Model, err error) {
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
		rolesSchema[key] = roleType
	}

	eftType := v.GetString("policy_effect.type")
	eftKey := v.GetString("policy_effect.key")

	var eft core.Effector
	switch eftType {
	case EffectorAllowOverride:
		eft = effector.NewAllowOverride()
	case EffectorDenyOverride:
		if eftKey == "" || !policy.Has(eftKey) {
			err = fmt.Errorf("policy effect: must give an effect field key")
			return
		}

		eft = effector.NewDenyOverride(eftKey)
	case EffectorAllowAndDeny:
		if eftKey == "" || !policy.Has(eftKey) {
			err = fmt.Errorf("policy effect: must give an effect field key")
			return
		}

		eft = effector.NewAllowAndDeny(eftKey)
	default:
		err = fmt.Errorf("unknown effector: %s", eftType)
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
		return
	}

	model = core.NewModel(policy, request, rolesSchema, eft, matchers)

	return
}
