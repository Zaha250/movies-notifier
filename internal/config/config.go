package config

type AppConfig struct {
	Port        int    `env:"PORT" envDefault:"8080"`
	PostgresURL string `env:"POSTGRES_URL"`
	Telegram    TelegramConfig
}

type TelegramConfig struct {
	Token string `env:"TG_TOKEN,required"`
}
