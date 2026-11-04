package main

import (
	"log/slog"
	"os"

	"file-temp-share/backend/docs"
	"file-temp-share/backend/internal/server"
)

//	@title			file-temp-share API
//	@version		0.1.0
//	@description	Temporary file-sharing API (WeTransfer-like).
//	@BasePath		/

//	@contact.name	file-temp-share
//	@contact.url	https://github.com/Remxin/file-temp-share

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		host = "localhost:8080"
	}
	docs.SwaggerInfo.Host = host

	engine := server.New()

	slog.Info("api listening", "addr", addr, "swagger", "http://"+host+"/swagger/index.html")
	if err := engine.Run(addr); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
