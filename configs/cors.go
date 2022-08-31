package configs

import "github.com/gin-contrib/cors"

var Cors = cors.Config{
	AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE", "PUT", "PATCH"},
	AllowHeaders:     []string{"Origin", "authorization", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	AllowCredentials: true,
	AllowOrigins:     []string{"*"},
	MaxAge:           86400,
}
