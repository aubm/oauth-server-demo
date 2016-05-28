package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RangelReale/osin"
	"github.com/aubm/oauth-server-demo/api"
	"github.com/aubm/oauth-server-demo/config"
	"github.com/aubm/oauth-server-demo/security"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

var appConfig config.App

func init() {
	flag.StringVar(&appConfig.Port, "port", "8080", "the tcp port for the application")
	flag.StringVar(&appConfig.DB.Addr, "db-addr", "localhost:3306", "the MySQL address")
	flag.StringVar(&appConfig.DB.User, "db-user", "root", "the mMySQL user")
	flag.StringVar(&appConfig.DB.Password, "db-password", "root", "the mMySQL password")
	flag.StringVar(&appConfig.DB.Name, "db-name", "oauthserverdemo", "the name of the MySQL database")
	flag.StringVar(&appConfig.Redis.Addr, "redis-addr", "localhost:6379", "the addr for the redis instance")
	flag.StringVar(&appConfig.Redis.Password, "redis-password", "", "the password for the redis instance")
	flag.Int64Var(&appConfig.Redis.DB, "redis-db", 0, "the Redis database to use")
	flag.IntVar(&appConfig.Security.AccessExpiration, "access-expiration", 3600, "the access token expiration time")
	flag.StringVar(&appConfig.Security.Secret, "secret", "this-is-not-really-a-secret", "the application secret")
	flag.Parse()

	fmt.Println("Parsed application parameters")
	fmt.Printf("TCP port:          %v\n", appConfig.Port)
	fmt.Printf("MySQL address:     %v\n", appConfig.DB.Addr)
	fmt.Printf("MySQL user:        %v\n", appConfig.DB.User)
	fmt.Printf("MySQL password:    %v\n", appConfig.DB.Password)
	fmt.Printf("MySQL database:    %v\n", appConfig.DB.Name)
	fmt.Printf("Redis address:     %v\n", appConfig.Redis.Addr)
	fmt.Printf("Redis password:    %v\n", appConfig.Redis.Password)
	fmt.Printf("Redis db:          %v\n", appConfig.Redis.DB)
	fmt.Printf("Access expiration: %v\n", appConfig.Security.AccessExpiration)
	fmt.Printf("Secret:            %v\n", appConfig.Security.Secret)
}

func main() {
	db, err := initDB()
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
	clientsManager := &security.ClientsManager{}
	accessDataManager := &security.AccessDataManager{}
	securityHandlers := &api.SecurityHandlers{}
	usersHandlers := &api.UsersHandlers{}
	usersManager := &security.UsersManager{}
	identityAdapter := &api.IdentityAdapter{}
	clearContextAdapter := &api.ClearContextAdapter{}
	logAdapter := &api.LogAdapter{}
	loggerInfo := log.New(os.Stdout, "info : ", log.Ldate|log.Ltime)
	loggerError := &RedPrinter{Logger: log.New(os.Stdout, "error: ", log.Ldate|log.Ltime|log.Lshortfile)}

	if err := populate(db, redisClient, oauthServerStorage, server, clientsManager, accessDataManager,
		securityHandlers, usersHandlers, usersManager, identityAdapter, clearContextAdapter, logAdapter,
		&inject.Object{Value: loggerInfo, Name: "logger_info"},
		&inject.Object{Value: loggerError, Name: "logger_error"}); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth/v1/token", securityHandlers.Token).Methods("POST")
	router.HandleFunc("/api/v1/users", usersHandlers.Create).Methods("POST")
	router.Handle("/api/v1/me", api.Adapt(http.HandlerFunc(usersHandlers.Me), identityAdapter)).Methods("GET")

	http.Handle("/", api.Adapt(router, clearContextAdapter, logAdapter))

	brand()
	fmt.Printf("Application started on port %v\n", appConfig.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", appConfig.Port), nil))
}

func populate(values ...interface{}) error {
	graph := inject.Graph{}
	for _, v := range values {
		if obj, ok := v.(*inject.Object); ok {
			graph.Provide(obj)
		} else {
			graph.Provide(&inject.Object{Value: v})
		}
	}
	return graph.Populate()
}

func createServerConfig() *osin.ServerConfig {
	serverConfig := osin.NewServerConfig()
	serverConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{}
	serverConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.REFRESH_TOKEN}
	serverConfig.AccessExpiration = int32(appConfig.Security.AccessExpiration)
	serverConfig.ErrorStatusCode = 400
	return serverConfig
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v", appConfig.DB.User, appConfig.DB.Password, appConfig.DB.Addr, appConfig.DB.Name))
	if err != nil {
		return db, err
	}
	for !pingDB(db) {
		log.Printf("database is not ready, will retry in 5 seconds")
		time.Sleep(time.Second * 5)
	}
	return db, err
}

func pingDB(db *sql.DB) bool {
	rows, err := db.Query("SELECT 1 FROM clients")
	if err == nil {
		defer rows.Close()
		return true
	}
	return false
}

func brand() {
	fmt.Print(`
.oPYo.      .oo          o  8        .oPYo.                                    ooo.
8    8     .P 8          8  8        8                                         8  '8.
8    8    .P  8 o    o  o8P 8oPYo.   'Yooo. .oPYo. oPYo. o    o .oPYo. oPYo.   8   '8 .oPYo. ooYoYo. .oPYo.
8    8   oPooo8 8    8   8  8    8       '8 8oooo8 8  '' Y.  .P 8oooo8 8  ''   8    8 8oooo8 8' 8  8 8    8
8    8  .P    8 8    8   8  8    8        8 8.     8     'b..d' 8.     8       8   .P 8.     8  8  8 8    8
'YooP' .P     8 'YooP'   8  8    8   'YooP' 'Yooo' 8      'YP'  'Yooo' 8       8ooo'  'Yooo' 8  8  8 'YooP'
:.....:..:::::..:.....:::..:..:::..:::.....::.....:..::::::...:::.....:..::::::.....:::.....:..:..:..:.....:
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

`)
}

type RedPrinter struct {
	Logger *log.Logger
}

func (p *RedPrinter) Printf(format string, v ...interface{}) {
	color.Set(color.FgRed)
	p.Logger.Printf(format, v...)
	color.Unset()
}
