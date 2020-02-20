package routes

import (
	"api/controllers"
	"net/http"
)

var twitterRoutes = []Route{
	Route{
		Uri:          "/api/twitter/tweets/search",
		Method:       http.MethodPost,
		Handler:      controllers.SearchTweets,
		AuthRequired: true,
	},
	Route{
		Uri:          "/api/twitter/tweets",
		Method:       http.MethodGet,
		Handler:      controllers.GetAllTweet,
		AuthRequired: true,
	},
	Route{
		Uri:          "/api/twitter/tweets",
		Method:       http.MethodPost,
		Handler:      controllers.CreateTweet,
		AuthRequired: true,
	},
	Route{
		Uri:          "/api/twitter/tweets/{id}",
		Method:       http.MethodGet,
		Handler:      controllers.GetTweet,
		AuthRequired: true,
	},
	Route{
		Uri:          "/api/twitter/tweets/{id}",
		Method:       http.MethodPut,
		Handler:      controllers.UpdateTweet,
		AuthRequired: true,
	},
	Route{
		Uri:          "/api/twitter/tweets/{id}",
		Method:       http.MethodDelete,
		Handler:      controllers.DeleteTweet,
		AuthRequired: true,
	},
}
