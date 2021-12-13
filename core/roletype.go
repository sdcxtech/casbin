package core

import (
	"fmt"
	"strings"
)

type RoleType int

const (
	RoleTypeInvalid RoleType = iota
	RoleTypeWithoutDomain
	RoleTypeWithDomain
)

func RoleTypeFromLine(line string) (RoleType, error) {
	ss := strings.Split(line, ",")
	l := len(ss)
	if l == 2 {
		return RoleTypeWithoutDomain, nil
	} else if l == 3 {
		return RoleTypeWithDomain, nil
	}

	return RoleTypeInvalid, fmt.Errorf("invalid role definition: %s", line)
}
