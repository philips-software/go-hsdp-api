package iam

type TokenRefresher interface {
	TokenRefresh() error
}

type HTTPStatus interface {
	StatusCode() int
}
