package iam

type TokenRefresher interface {
	TokenRefresh() error
}
