package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"WebService/data"
	"WebService/log"
	"WebService/model"
)

var ds *data.DataStore

const VERSION = "1.0.0"

func Router(d *data.DataStore) *gin.Engine {
	ds = d
	var mode string
	if log.MaxLogLevel >= log.TRACE {
		mode = gin.DebugMode
	} else {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	router := gin.Default()
	router.Use(decodeSession)


	router.Static("/www", "www")
	api(router)

	return router
}
func api(router *gin.Engine) {
	api := router.Group("/api")
	//apiCache := api.Group("")
	api.Use(noCache)
	apiAuth := api.Group("")
	apiAuth.Use(sessionValidator)

	api.GET("/exampleeeeeeeeee", GetExample)
	api.POST("/", PostExample)

}

func decodeSession(context *gin.Context) {
	session := context.Request.Header.Get("Session")
	var err error
	if len(session) == 0 {
		session, err = context.Cookie("session")
		if err != nil {
			return
		}
		if len(session) == 0 {
			return
		}
	}
	claims, err := model.DecodeSession(session)
	if err != nil {
		DelSessionCookie(context)
		return
	}
	duration := time.Unix(claims.ExpiresAt, 0).Sub(time.Now()).Minutes()
	if duration < 0 {
		DelSessionCookie(context)
		return
	} else if duration < ds.HttpConfig.SessionTime.Minutes()/4 {
		expires := time.Now().Add(ds.HttpConfig.SessionTime).Unix()
		SetSessionCookie(context, model.NewSessionTime(claims.Uid, claims.Username, claims.Role, claims.Token, expires))
	}
	context.Set("session", claims)
}
func sessionValidator(context *gin.Context) {
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
	}
	claim := getSession(context)
	if claim == nil {
		response.Status = "error"
		response.StatusMessage = "Unauthorized"
		context.JSON(http.StatusUnauthorized, response)
		//context.Redirect(http.StatusMovedPermanently, "/login")
		context.Abort()
		return
	}

}

func getSession(context *gin.Context) *model.Claims {
	if s, ok := context.Get("session"); ok {
		if session, ok := s.(*model.Claims); ok {
			context.Set("user", session.Username)
			return session
		}
	}
	return nil
}
func noCache(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Next()
}
func getPaginationValue(context *gin.Context) (start int, limit int) {
	getStart, _ := context.GetQuery("start")
	start, _ = strconv.Atoi(getStart)
	getLimit, _ := context.GetQuery("limit")
	limit, _ = strconv.Atoi(getLimit)
	return
}
