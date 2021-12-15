package load

import "errors"

var (
	ErrNeedEffectFieldKey = errors.New("must give an effect field key")
	ErrUnknownEffectType  = errors.New("unknown policy effect type")
)
