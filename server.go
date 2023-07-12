package main

import (
	"fmt"
	"log"
	"net/http"
	"nicepay-service/payments"
	"os"
	"sync"

	"github.com/go-playground/validator"
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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func customErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				report.Message = fmt.Sprintf("%s is required",
					err.Field())
			case "email":
				report.Message = fmt.Sprintf("%s is not valid email",
					err.Field())
			case "gte":
				report.Message = fmt.Sprintf("%s value must be greater than %s",
					err.Field(), err.Param())
			case "lte":
				report.Message = fmt.Sprintf("%s value must be lower than %s",
					err.Field(), err.Param())
			}

			break
		}
	}

	c.Logger().Error(report)
	c.JSON(report.Code, report)
}

func main() {
	e := echo.New()

	// Middleware
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = customErrorHandler

	logFile, err := os.OpenFile("nicepay-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		/* var jsonb = []byte(reqBody)
		reqBodyCompact := new(bytes.Buffer)
		json.Compact(reqBodyCompact, jsonb)

		log.Println("Request :" + reqBodyCompact.String())
		log.Print("Response :" + string(resBody)) */
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", getAllUsers)
	ep := e.Group("/payments")
	ep.POST("/registration", payments.CreateRegistration)
	ep.POST("/", payments.CreatePayment)
	ep.POST("/show-inquery", payments.GetInquiry)

	// Start server
	e.Logger.Fatal(e.Start(":1212"))

}
