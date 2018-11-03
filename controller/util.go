package controller

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func SetSessionCookie(context *gin.Context, session string) {
	i := int(ds.HttpConfig.SessionTime.Seconds())
	context.SetCookie("session", session, i, "/", "", false, true)
}

func DelSessionCookie(context *gin.Context) {
	context.SetCookie("session", "", -1, "/", "", false, true)
}

func convertMillToTime(t int64) time.Time {
	if t <= 0 {
		return time.Time{}
	}
	return time.Unix(0, t*int64(time.Millisecond))
}
func convertTimeToMill(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / int64(time.Millisecond)
}

func checkRoles(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 20)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
