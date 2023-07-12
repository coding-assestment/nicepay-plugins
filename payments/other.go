package payments

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type (
	User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	users = map[int]*User{}
	seq   = 1
	lock  = sync.Mutex{}
)

func CreatePayment(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &User{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func ShowInquery(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &User{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}
