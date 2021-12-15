package core

import "errors"

var (
	ErrNotFoundMatcher                   = errors.New("not found specified matcher")
	ErrInvalidAssertionSchema            = errors.New("invalid assertion schema")
	ErrInvalidAssertion                  = errors.New("invalid assertion")
	ErrInvalidMatchersConfig             = errors.New("invalid matchers config")
	ErrInvalidRoleTypeDefinition         = errors.New("invalid role type definition")
	ErrRoleMappingMismatchRoleType       = errors.New("role mapping mismatch role type")
	ErrUnknownAssertionType              = errors.New("unknown assertion type")
	ErrNotFoundSpecifiedRoleMappingGroup = errors.New("not found role mapping group")
	ErrUnexpectedEvalResult              = errors.New("expected eval result")
)
