package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAssertionIterator struct {
	data  [][]string
	index int
}

func (it *mockAssertionIterator) Next() (ok bool, key string, vals []string) {
	ok = true
	if it.index >= len(it.data) {
		ok = false
		return
	}
	line := it.data[it.index]
	key = line[0]
	vals = line[1:]
	it.index++
	return
}

func (it *mockAssertionIterator) Error() (err error) {
	return nil
}

func TestModelLoad(t *testing.T) {
	policy, err := NewAssertionSchema("sub, obj, act")
	assert.NoError(t, err)

	request, err := NewAssertionSchema("sub, dom, obj, act")
	assert.NoError(t, err)

	rolesSchema := make(RolesSchema)
	rolesSchema["g"] = RoleTypeWithDomain
	rolesSchema["g1"] = RoleTypeWithoutDomain

	matchers, err := MatchersConfig{
		Roles: rolesSchema,
		Define: map[string]string{
			"m": "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act",
		},
	}.New()
	assert.NoError(t, err)

	m := NewModel(policy, request, rolesSchema, nil, matchers)

	itr := &mockAssertionIterator{
		data: [][]string{
			{"p", "alice", "order", "get"},
			{"p", "alice", "order", "set"},
			{"g", "alice", "admin", "console"},
			{"g1", "alice", "admin"},
		},
	}

	policies, _, err := m.Load(itr)
	assert.NoError(t, err)
	assert.Len(t, policies, 2)
}
