package load

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadModelSucceeds(t *testing.T) {
	{
		conf := `
[request_definition]
r = "sub, dom, obj, act"

[policy_definition]
p = "sub, obj, act"

[role_definition]
g =  "_, _, _"
g1 = "_, _"

[policy_effect]
type = "allow-override"

[matchers]
m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
    `

		v := viper.New()
		v.SetConfigType("toml")
		err := v.ReadConfig(strings.NewReader(conf))
		assert.NoError(t, err)

		_, err = LoadModel(v)
		assert.NoError(t, err)

	}
	{
		conf := `
[request_definition]
r = "sub, dom, obj, act"

[policy_definition]
p = "sub, obj, act, eft"

[role_definition]
g =  "_, _, _"
g1 = "_, _"

[policy_effect]
type = "allow-and-deny"
key = "eft"

[matchers]
m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
    `

		v := viper.New()
		v.SetConfigType("toml")
		err := v.ReadConfig(strings.NewReader(conf))
		assert.NoError(t, err)

		_, err = LoadModel(v)
		assert.NoError(t, err)
	}
	{
		conf := `
[request_definition]
r = "sub, dom, obj, act"

[policy_definition]
p = "sub, obj, act, eft"

[role_definition]
g =  "_, _, _"
g1 = "_, _"

[policy_effect]
type = "deny-override"
key = "eft"

[matchers]
m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
    `

		v := viper.New()
		v.SetConfigType("toml")
		err := v.ReadConfig(strings.NewReader(conf))
		assert.NoError(t, err)

		_, err = LoadModel(v)
		assert.NoError(t, err)
	}
}

func TestLoadModelInvalidRequest(t *testing.T) {
	{
		conf := `
[request_definition]
r = "sub, dom, , act"

[policy_definition]
p = "sub, obj, act"

[role_definition]
g =  "_, _, _"
g1 = "_, _"

[policy_effect]
type = "allow-override"

[matchers]
m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
    `

		v := viper.New()
		v.SetConfigType("toml")
		err := v.ReadConfig(strings.NewReader(conf))
		assert.NoError(t, err)

		_, err = LoadModel(v)
		assert.Error(t, err)
	}
	{
		conf := `
[request_definition]
r = "sub, dom, obj, act"

[policy_definition]
p = "sub, , act"

[role_definition]
g =  "_, _, _"
g1 = "_, _"

[policy_effect]
type = "allow-override"

[matchers]
m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
    `

		v := viper.New()
		v.SetConfigType("toml")
		err := v.ReadConfig(strings.NewReader(conf))
		assert.NoError(t, err)

		_, err = LoadModel(v)
		assert.Error(t, err)
	}
}

func TestLoadModelLeakMatchers(t *testing.T) {
	{
		conf := `
[request_definition]
r = "sub, dom, obj, act"

[policy_definition]
p = "sub, obj, act"

[role_definition]
g =  "_, _, _"
g1 = "_, _"

[policy_effect]
type = "allow-override"

[matchers]
    `

		v := viper.New()
		v.SetConfigType("toml")
		err := v.ReadConfig(strings.NewReader(conf))
		assert.NoError(t, err)

		_, err = LoadModel(v)
		assert.Error(t, err)
	}
}
