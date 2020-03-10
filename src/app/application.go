package app

import (
	"github.com/StefanPahlplatz/bookstore_oauth-api/src/clients/cassandra"
	"github.com/StefanPahlplatz/bookstore_oauth-api/src/domain/access_token"
	"github.com/StefanPahlplatz/bookstore_oauth-api/src/http"
	"github.com/StefanPahlplatz/bookstore_oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func StartApplication() {
	session, dbErr := cassandra.GetSession()
	if dbErr != nil {
		panic(dbErr)
	}
	session.Close()

	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))

	router.GET("oauth/access_token/:access_token_id", atHandler.GetById)

	router.Run(":8080")
}
