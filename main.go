package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/aubm/oauth-server-demo/api"
	"github.com/aubm/oauth-server-demo/security"
	"github.com/facebookgo/inject"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "the tcp port for the application")
	flag.Parse()
}

func main() {
	db, err := sql.Open("mysql", "root:root@/oauthserverdemo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
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

	fmt.Printf("Server started on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func createServerConfig() *osin.ServerConfig {
	serverConfig := osin.NewServerConfig()
	serverConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{}
	serverConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.REFRESH_TOKEN}
	return serverConfig
}
