package config

type AppConfig struct {
	Port        int    `env:"PORT" envDefault:"8080"`
	PostgresUrl string `env:"POSTGRES_URL"`
	Telegram    TelegramConfig
}

type TelegramConfig struct {
	Token string `env:"TG_TOKEN,required"`
}
