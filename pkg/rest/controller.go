package rest

type Controller interface {
	Route() *RestRoute
}
