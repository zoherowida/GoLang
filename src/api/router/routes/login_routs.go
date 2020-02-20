package routes

import (
	"api/controllers"
	"net/http"
)

var loginRoutes = []Route{
	Route{
		Uri:          "/api/login",
		Method:       http.MethodPost,
		Handler:      controllers.Login,
		AuthRequired: false,
	},
	Route{
		Uri:          "/api/register",
		Method:       http.MethodPost,
		Handler:      controllers.Register,
		AuthRequired: false,
	},
}
