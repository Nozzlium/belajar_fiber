package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed source/chunli.txt
var chunli []byte

var app = fiber.New(fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		c.Status(fiber.StatusInternalServerError)
		return c.SendString(err.Error())
	},
})

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestGet(t *testing.T) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World")
	})

	response, err := app.Test(httptest.NewRequest("GET", "/", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World", string(bytes))
}

func TestRequestWithQuery(t *testing.T) {
	app.Get("/query", func(c *fiber.Ctx) error {
		name := c.Query("name", "World")
		return c.SendString(fmt.Sprintf("Hello, %s", name))
	})

	response1, err := app.Test(httptest.NewRequest("GET", "/query", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response1.StatusCode)

	bytes, err := io.ReadAll(response1.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World", string(bytes))

	response2, err := app.Test(httptest.NewRequest("GET", "/query?name=Shepard", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response2.StatusCode)

	bytes, err = io.ReadAll(response2.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Shepard", string(bytes))
}

func TestHeaderAndCookies(t *testing.T) {
	app.Get("/additional", func(c *fiber.Ctx) error {
		headerName := c.Get("name", "lah kocak")
		cookieName := c.Cookies("name", "gimana dah?")
		return c.SendString(fmt.Sprintf("%s %s", headerName, cookieName))
	})

	response, err := app.Test(httptest.NewRequest("GET", "/additional", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "lah kocak gimana dah?", string(bytes))

	req := httptest.NewRequest("GET", "/additional", nil)
	req.Header.Set("name", "wadidaw")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "wadidiw",
	})
	response, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "wadidaw wadidiw", string(bytes))
}

func TestRouteParam(t *testing.T) {
	app.Get("/param/:id/weleh", func(c *fiber.Ctx) error {
		id := c.Params("id", "lho?")
		return c.SendString(fmt.Sprintf("Ini dia: %s", id))
	})

	response, err := app.Test(httptest.NewRequest("GET", "/param/chunli/weleh", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Ini dia: chunli", string(bytes))
}

func TestForm(t *testing.T) {
	app.Post("/form", func(c *fiber.Ctx) error {
		name := c.FormValue("name", "lah?")
		game := c.FormValue("game", "lah?")
		return c.SendString(fmt.Sprintf("Here's: %s, she appears in %s", name, game))
	})

	reader := strings.NewReader("name=Chun-Li&game=Street Fighter 6")
	request := httptest.NewRequest("POST", "/form", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Here's: Chun-Li, she appears in Street Fighter 6", string(bytes))
}

func TestMultiForm(t *testing.T) {
	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		err = c.SaveFile(file, "./target/"+file.Filename)
		if err != nil {
			return err
		}

		return c.SendString("Success")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "chunli.txt")
	file.Write(chunli)
	writer.Close()

	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("content-type", writer.FormDataContentType())
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Success", string(bytes))
}

func TestBody(t *testing.T) {
	app.Post("/login", func(c *fiber.Ctx) error {
		body := c.Body()
		request := new(LoginRequest)
		err := json.Unmarshal(body, request)
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("Successfully logged in. Welcome, %s", request.Username))
	})

	body := strings.NewReader(`{
		"username": "Chun-Li",
		"password": "kungfu"
	}`)
	request := httptest.NewRequest("POST", "/login", body)
	request.Header.Add("content-type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, _ := io.ReadAll(response.Body)

	result := string(bytes)
	assert.Equal(t, "Successfully logged in. Welcome, Chun-Li", result)
}

func TestBodyParser(t *testing.T) {
	app.Post("/login", func(c *fiber.Ctx) error {
		request := new(LoginRequest)
		err := c.BodyParser(request)
		if err != nil {
			return err
		}

		return c.SendString(fmt.Sprintf("Successfully logged in. Welcome, %s", request.Username))
	})

	request := httptest.NewRequest("POST", "/login", strings.NewReader(`{
		"username": "Chun-Li",
		"password": "kungfu"
	}`))
	request.Header.Add("content-type", "application/json")
	response, _ := app.Test(request)
	bytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Successfully logged in. Welcome, Chun-Li", string(bytes))
}

func TestResponseBody(t *testing.T) {
	app.Get("/response_body", func(c *fiber.Ctx) error {
		return c.JSON(&LoginRequest{
			Username: "test",
			Password: "testjugajing",
		})
	})

	request := httptest.NewRequest("GET", "/response_body", nil)
	resp, _ := app.Test(request)
	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"username":"test","password":"testjugajing"}`, string(bytes))
}

func TestRoutingGroup(t *testing.T) {
	handlerFunc := func(ctx *fiber.Ctx) error {
		fmt.Println("Halo")
		return ctx.JSON(fiber.Map{
			"title": "Hello, World!",
		})
	}

	authMiddleware := func(ctx *fiber.Ctx) error {
		fmt.Println("auth maboy")
		return ctx.Next()
	}

	auth := app.Group("/auth")
	post := app.Group("/post")

	auth.Use(authMiddleware)

	post.Get("", handlerFunc)
	auth.Post("/login", handlerFunc)
	post.Post("", handlerFunc)

	app.Test(httptest.NewRequest("POST", "/auth/login", nil))
	app.Test(httptest.NewRequest("GET", "/post", nil))
	app.Test(httptest.NewRequest("POST", "/post", nil))
}

func TestErrorHandler(t *testing.T) {
	app.Get("/error", func(c *fiber.Ctx) error {
		return errors.New("chunli")
	})

	resp, _ := app.Test(httptest.NewRequest("GET", "/error", nil))
	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "chunli", string(bytes))
}
