package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/aubm/oauth-server-demo/api"
	"github.com/aubm/oauth-server-demo/config"
	"github.com/aubm/oauth-server-demo/security"
	"github.com/facebookgo/inject"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

var appConfig config.App

func init() {
	flag.StringVar(&appConfig.Port, "port", "8080", "the tcp port for the application")
	flag.StringVar(&appConfig.DB.User, "db-user", "root", "the mMySQL user")
	flag.StringVar(&appConfig.DB.Password, "db-password", "root", "the mMySQL password")
	flag.StringVar(&appConfig.DB.Name, "db-name", "oauthserverdemo", "the name of the MySQL database")
	flag.StringVar(&appConfig.Redis.Addr, "redis-addr", "localhost:6379", "the addr for the redis instance")
	flag.StringVar(&appConfig.Redis.Password, "redis-password", "", "the password for the redis instance")
	flag.Int64Var(&appConfig.Redis.DB, "redis-db", 0, "the Redis database to use")
	flag.IntVar(&appConfig.Security.AccessExpiration, "access-expiration", 3600, "the access token expiration time")
	flag.Parse()
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", appConfig.DB.User, appConfig.DB.Password, appConfig.DB.Name))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: appConfig.Redis.Addr, Password: appConfig.Redis.Password, DB: appConfig.Redis.DB})
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	oauthServerStorage := &security.Storage{}
	server := osin.NewServer(createServerConfig(), oauthServerStorage)
	securityHandlers := &api.SecurityHandlers{}
	clientsManager := &security.ClientsManager{}
	accessDataManager := &security.AccessDataManager{}

	if err := inject.Populate(db, oauthServerStorage, server, securityHandlers,
		clientsManager, redisClient, accessDataManager); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/auth/v1/token", securityHandlers.Token)
	http.HandleFunc("/api/v1/me", securityHandlers.Me)

	fmt.Printf("Server started on port %v\n", appConfig.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", appConfig.Port), nil))
}

func createServerConfig() *osin.ServerConfig {
	serverConfig := osin.NewServerConfig()
	serverConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{}
	serverConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.REFRESH_TOKEN}
	serverConfig.AccessExpiration = int32(appConfig.Security.AccessExpiration)
	return serverConfig
}
