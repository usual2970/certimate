package apisix

type apiResponse interface{}

type apiResponseBase struct{}

var _ apiResponse = (*apiResponseBase)(nil)
