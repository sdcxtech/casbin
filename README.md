# casbin

Another [casbin](https://casbin.org/) implementation in golang.

## Installation

```bash
go get github.com/sdcxtech/casbin
```

## Quick start

```go
package main

import (
	"fmt"
	"strings"

	"github.com/sdcxtech/casbin"
	"github.com/sdcxtech/casbin/core"
	"github.com/sdcxtech/casbin/effector"
	"github.com/sdcxtech/casbin/load"
)

func main() {
	policy, _ := core.NewAssertionSchema("sub, obj, act")
	request, _ := core.NewAssertionSchema("sub, obj, act")

	roles := core.RolesSchema{
		"g": {Type: core.RoleTypeWithoutDomain},
	}

	matchers, _ := core.MatchersConfig{
		Roles: roles,
		Define: map[string]string{
			"m": "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act",
		},
	}.New()

	model := core.NewModel(policy, request, roles, effector.NewAllowOverride(), matchers)

	policies := strings.NewReader(`
p, staff, order, get
g, admin, staff
`)

	enforcer, _ := casbin.NewEnforcer(model, load.NewCSVIterator(policies))
	allow, _ := enforcer.Enforce(casbin.Request("admin", "order", "get"))

	fmt.Println(allow) // true
}
```

## Model config

Models can also be loaded from a `viper.Viper` instance, so the concrete file
format can be TOML, YAML, JSON, or anything else Viper supports.

```toml
[request_definition]
r = "sub, dom, obj, act"

[policy_definition]
p = "sub, obj, act, eft"

[role_definition]
g = "_, _, _"
g1 = "_, _"

[policy_effect]
type = "allow-and-deny"
key = "eft"

[matchers]
m = "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act"
```

```go
model, err := load.ModelFromViper(v)
```

Supported effect types:

* `allow-override`: allow when any policy matches. No effect field is required.
* `deny-override`: allow unless a matched policy has effect `deny`.
* `allow-and-deny`: allow when a policy with effect `allow` matches and no
  matched policy has effect `deny`.

`deny-override` and `allow-and-deny` require `policy_effect.key` to name a
policy field whose value is `allow` or `deny`.

## Policies

Use `load.NewCSVIterator` to load policy and role rows from an `io.Reader`.
The first CSV column is the section key, such as `p` for policies or `g` for a
role relation. Lines starting with `#` are comments.

```csv
p, staff, order, get
g, admin, staff
```

For a role definition with domain (`g = "_, _, _"`), role rows use
`from, to, domain`:

```csv
g, alice, admin, tenant1
```

## Matchers

Matchers are CEL expressions. Request values are available as `r.<field>`,
policy values are available as `p.<field>`, and role functions use the names
from `[role_definition]`.

You can register string match helpers as CEL extension functions:

```go
model, err := load.ModelFromViper(
	v,
	load.ExtensionFuncs(
		keymatch.ToExtensionFunc("globMatch", keymatch.GlobMatch),
		keymatch.ToExtensionFunc("regexMatch", keymatch.NewRegexMatch()),
		keymatch.ToExtensionFunc("ipMatch", keymatch.IPMatch),
	),
)
```

Then use them from a matcher:

```toml
[matchers]
m = "g(r.sub, p.sub) && globMatch(r.obj, p.obj) && r.act == p.act"
```

## Per-request options

`Enforce` accepts options for one call without changing the enforcer:

* `casbin.UseMatcher(name)` selects a matcher other than the default `m`.
* `casbin.WithRawPolicies(rows)` injects extra raw policies for this call.
* `casbin.WithPolicies(policies)` injects parsed policies for this call.
* `casbin.OnlyInjectedPolicies()` ignores the enforcer's loaded policies.
* `casbin.WithRoleGraphs(key, graphs...)` injects role graphs for this call.

## Difference with the official casbin implementation

* Use google [Common Expression Language](https://github.com/google/cel-go) as the matcher expression language.
* Assertion field in policy and request only can be `string` type. So there is no support for `ABAC` model.
* Only implement the core feature checking permissions. Not include policies and roles management
  which should be implemented in a different [Bounded Context](https://martinfowler.com/bliki/BoundedContext.html).

## License

Released under the [MIT License](LICENSE).
