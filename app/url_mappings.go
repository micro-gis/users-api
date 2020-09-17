package app

import "github.com/micro-gis/users-api/controllers"

func mapUrls() {
	router.GET("/ping", controllers.Ping)
}
