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
	oauthServerStorage := &security.Storage{}
	server := osin.NewServer(createServerConfig(), oauthServerStorage)
	securityHandlers := &api.SecurityHandlers{}

	if err := inject.Populate(db, oauthServerStorage, server, securityHandlers); err != nil {
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
