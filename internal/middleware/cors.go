package middleware

import (
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

// CORSMiddleware sets up CORS only in development.
// In production it becomes a no-op.
func CORSMiddleware() echo.MiddlewareFunc {
	if os.Getenv("APP_ENV") != "dev" {
		// No-op outside of dev
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}
	}

	return emw.CORSWithConfig(emw.CORSConfig{
		// Allow any origin on localhost (any port)
		AllowOriginFunc: func(origin string) (bool, error) {
			u, err := url.Parse(origin)
			if err != nil {
				return false, err
			}
			host := u.Hostname()
			scheme := u.Scheme
			// accept only http or https on host "localhost"
			ok := (scheme == "http" || scheme == "https") && host == "localhost"
			return ok, nil
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
		},
	})
}
