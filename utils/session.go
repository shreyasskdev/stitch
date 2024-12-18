package utils

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CreateSession(c echo.Context, name string) {
	var sess, err = session.Get(name, c)
	if err != nil {
		panic(err)
	}

	if sess.IsNew {
		sess.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			MaxAge:   60 * 60 * 24 * 365,
		}
	}
}

func AddSessionValue(c echo.Context, name string, key string, value interface{}) {
	var sess, err = session.Get(name, c)
	if err != nil {
		panic(err)
	}

	sess.Values[key] = value
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		panic(err)
	}
}

func GetSessionValue(c echo.Context, name string, key string) interface{} {
	var sess, err = session.Get(name, c)
	if err != nil {
		panic(err)
	}

	var value = sess.Values[key]
	return value
}

func DeleteSession(c echo.Context, name string) {
	var sess, err = session.Get(name, c)
	if err != nil {
		panic(err)
	}

	sess.Options.MaxAge = -1
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		panic(err)
	}
}

func GetSessionAll(c echo.Context, name string) map[interface{}]interface{} {
	var sess, err = session.Get(name, c)
	if err != nil {
		fmt.Println(err.Error())
		return map[interface{}]interface{}{}
	}

	return sess.Values
}

func DeleteSessionValue(c echo.Context, name string, key string) {
	var sess, err = session.Get(name, c)
	if err != nil {
		fmt.Println(err)
		return
	}
	delete(sess.Values, key)

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		panic(err)
	}
}
