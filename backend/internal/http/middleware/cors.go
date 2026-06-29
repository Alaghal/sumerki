package middleware

import "github.com/labstack/echo/v4/middleware"

func LocalCORS() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
}
