package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	app := fiber.New()
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
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		name := c.Query("name", "World")
		return c.SendString(fmt.Sprintf("Hello, %s", name))
	})

	response1, err := app.Test(httptest.NewRequest("GET", "/", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response1.StatusCode)

	bytes, err := io.ReadAll(response1.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World", string(bytes))

	response2, err := app.Test(httptest.NewRequest("GET", "/?name=Shepard", nil))
	assert.Nil(t, err)
	assert.Equal(t, 200, response2.StatusCode)

	bytes, err = io.ReadAll(response2.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Shepard", string(bytes))
}

func TestHeaderAndCookies(t *testing.T) {
	app := fiber.New()
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
	app := fiber.New()
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
	app := fiber.New()
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
