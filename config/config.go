package config

type App struct {
	Port string
	DB   struct {
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int64
	}
	Security struct {
		AccessExpiration int
	}
}
