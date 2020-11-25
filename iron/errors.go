package iron

import "errors"

var (
	ErrBaseIRONURLCannotBeEmpty = errors.New("base IRON URL cannot be empty")
	ErrNotImplemented           = errors.New("not implemented")
	ErrNotFound                 = errors.New("not found")
	ErrInvalidDockerCredentials = errors.New("invalid docker credentials. all fields required")
	ErrNoPublicKey              = errors.New("no public key present")
)
