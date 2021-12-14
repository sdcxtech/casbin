package core

import (
	"fmt"
	"strings"
)

// Role mapping type.
type RoleType int

const (
	// Invalid role mapping type.
	RoleTypeInvalid RoleType = iota
	// Role mapping without domain.
	RoleTypeWithoutDomain
	// Role mapping with domain.
	RoleTypeWithDomain
)

// Get the role type from the definition in official casbin format.
//
// * "_, _" is without domain.
// * "_, _, _" is with domain.
func RoleTypeFromLine(line string) (RoleType, error) {
	ss := strings.Split(line, ",")

	if l := len(ss); l == 2 {
		return RoleTypeWithoutDomain, nil
	} else if l == 3 {
		return RoleTypeWithDomain, nil
	}

	return RoleTypeInvalid, fmt.Errorf("invalid role definition: %s", line)
}
