package main

import (
	"net/http"
	"nicepay-service/payments"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	users           = map[int]*User{}
	seq             = 1
	lock            = sync.Mutex{}
	defaultResponse map[string]interface{}
)

//----------
// Handlers
//----------

/*
	 func createUser(c echo.Context) error {
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

	func getUser(c echo.Context) error {
		lock.Lock()
		defer lock.Unlock()
		id, _ := strconv.Atoi(c.Param("id"))
		return c.JSON(http.StatusOK, users[id])
	}

	func updateUser(c echo.Context) error {
		lock.Lock()
		defer lock.Unlock()
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		id, _ := strconv.Atoi(c.Param("id"))
		users[id].Name = u.Name
		return c.JSON(http.StatusOK, users[id])
	}

	func deleteUser(c echo.Context) error {
		lock.Lock()
		defer lock.Unlock()
		id, _ := strconv.Atoi(c.Param("id"))
		delete(users, id)
		return c.NoContent(http.StatusNoContent)
	}

	func getAllUsers(c echo.Context) error {
		lock.Lock()
		defer lock.Unlock()
		return c.JSON(http.StatusOK, users)
	}
*/
func getAllUsers(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &User{
		ID:   seq,
		Name: "testing 1",
	}
	users[u.ID] = u
	u = &User{
		ID:   seq,
		Name: "testing 2",
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusOK, users)
}

// func createRegistration(c echo.Context) error {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	u := &User{
// 		ID: seq,
// 	}
// 	if err := c.Bind(u); err != nil {
// 		return err
// 	}
// 	users[u.ID] = u
// 	seq++
// 	return c.JSON(http.StatusCreated, u)
// }

// func createPayment(c echo.Context) error {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	u := &User{
// 		ID: seq,
// 	}
// 	if err := c.Bind(u); err != nil {
// 		return err
// 	}
// 	users[u.ID] = u
// 	seq++
// 	return c.JSON(http.StatusCreated, u)
// }

// func showInquery(c echo.Context) error {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	u := &User{
// 		ID: seq,
// 	}
// 	if err := c.Bind(u); err != nil {
// 		return err
// 	}
// 	users[u.ID] = u
// 	seq++
// 	return c.JSON(http.StatusCreated, u)
// }

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", getAllUsers)
	/*
		e.GET("/", getAllUsers)
		e.GET("/users", getAllUsers)
		e.POST("/users", createUser)
		e.GET("/users/:id", getUser)
		e.PUT("/users/:id", updateUser)
		e.DELETE("/users/:id", deleteUser) */

	ep := e.Group("/payments")
	ep.POST("/registration", payments.CreateRegistration)
	ep.POST("/", payments.CreatePayment)
	ep.POST("/show-inquery", payments.ShowInquery)

	// Start server
	e.Logger.Fatal(e.Start(":1212"))
	// e.Logger.Fatal(e.Start(":80"))
}
