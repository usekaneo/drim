package generator

type Config struct {
	Domain       string
	UseCaddy     bool
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	AuthSecret   string
	APIPort      string
	WebPort      string
}

func NewDefaultConfig() *Config {
	return &Config{
		Domain:       "",
		UseCaddy:     true,
		PostgresUser: "kaneo",
		PostgresPass: generateRandomPassword(32),
		PostgresDB:   "kaneo",
		AuthSecret:   generateRandomPassword(64),
		APIPort:      "1337",
		WebPort:      "5173",
	}
}
