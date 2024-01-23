package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
