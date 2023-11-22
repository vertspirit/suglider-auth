package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	sqltable "suglider-auth/init/sql_table"
	mariadb "suglider-auth/internal/database"
	"suglider-auth/internal/redis"
	"suglider-auth/configs"
	"suglider-auth/pkg/api-server"
	"suglider-auth/pkg/time_convert"
	"suglider-auth/pkg/logger"
	"log/slog"
)

type (
	authApiSettings = api_server.AuthApiSettings
)

const ServiceName string = "Suglider-Auth"

var (
	Version           string
	Build             string
	swag              gin.HandlerFunc
)

// @title                      Suglider-Auth API Doc
// @version                    1.0
// @description                API Doc for Suglider Auth
// @termsOfService             http://swagger.io/terms/
// @contact.name               Above&Beyond
// @contact.email              geek@openmind.np
// @license.name               MIT
// @license.url                https://github.com/tony/suglider-auth/blob/master/LICENSE
// @schemes                    http https
// @BasePath                   /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       X-API-KEY

func init() {
	logger.SlogMultiWriter(ServiceName)
	sqltable.SugliderTableInit()
	_, _, err := time_convert.ConvertTimeFormat(configs.ApplicationConfig.Session.Timeout)
	if err != nil {
		errorMessage := fmt.Sprintf("TTL string convert to duration failed: %v", err)
		slog.Error(errorMessage)
		panic(err)
	}
}

func main() {
	apiServer := &authApiSettings {
		Name:             ServiceName,
		Version:          Version,
		TemplatePath:     configs.ApplicationConfig.Server.TemplatePath,
		StaticPath:       configs.ApplicationConfig.Server.StaticPath,
		SessionsPath:     configs.ApplicationConfig.Session.Path,
		SessionsHttpOnly: configs.ApplicationConfig.Session.HttpOnly,
		CasbinConfig:     configs.ApplicationConfig.Server.CasbinConfig,
		CasbinTable:      configs.ApplicationConfig.Server.CasbinTable,
		CasbinCache:      configs.ApplicationConfig.Server.CasbinCache,
		ReadTimeout:      configs.ApplicationConfig.Server.ReadTimeout,
		WriteTimeout:     configs.ApplicationConfig.Server.WriteTimeout,
		MaxHeaderBytes:   configs.ApplicationConfig.Server.MaxHeaderBytes,
		EnableRbac:       configs.ApplicationConfig.Server.EnableRbac,
		EnablePprof:      configs.Args.Debug,
	}

	if configs.Args.Subpath != "" {
		apiServer.SubpathPrefix = configs.Args.Subpath
	}

	fmt.Printf("[%s] Version:    %s\n", ServiceName, Version)
	fmt.Printf("[%s] Build Date: %s\n", ServiceName, Build)

	addr := fmt.Sprintf("%s:%d", configs.Args.Host, configs.Args.Port)
	apiServer.StartServer(addr, swag)

	mariadb.Close()
	slog.Info("Close database connection.")

	redis.Close()
	slog.Info("Close redis connection.")

}
