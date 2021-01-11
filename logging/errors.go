package logging

import "errors"

var (
	ErrMissingCredentialsOrIAMClient = errors.New("missing signing credentials or IAM client")
	ErrNothingToPost                 = errors.New("nothing to post")
	ErrMissingSharedKey              = errors.New("missing shared key")
	ErrMissingSharedSecret           = errors.New("missing shared secret")
	ErrMissingBaseURL                = errors.New("missing base URL")
	ErrMissingProductKey             = errors.New("missing ProductKey")
	ErrBatchErrors                   = errors.New("batch errors. check Invalid map for details")
	ErrResponseError                 = errors.New("unexpected HSDP response error")
)
