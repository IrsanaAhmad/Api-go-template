package config

type Config struct {
	DatabaseURL string         `env:"DATABASE_URL"`
	App         AppConfig      `envPrefix:"APP_"`
	Database    DatabaseConfig `envPrefix:"DB_"`
	Redis       RedisConfig    `envPrefix:"REDIS_"`
	Mail        MailConfig     `envPrefix:"MAIL_"`
	Pusher      PusherConfig   `envPrefix:"PUSHER_"`
	JWT         JWTConfig      `envPrefix:"JWT_"`
}

type AppConfig struct {
	Name  string `env:"NAME"            envDefault:"MyApp"`
	Env   string `env:"ENV"             envDefault:"development"`
	Key   string `env:"KEY"`
	Debug bool   `env:"DEBUG"           envDefault:"false"`
	URL   string `env:"URL"`
	Base  string `env:"BASE_URL"`
	Port  string `env:"PORT"            envDefault:"8080"`
}

type DatabaseConfig struct {
	Connection             string `env:"CONNECTION"                envDefault:"postgres"`
	Host                   string `env:"HOST"                      envDefault:"localhost"`
	Port                   int    `env:"PORT"                      envDefault:"1433"`
	Name                   string `env:"NAME"`
	Username               string `env:"USERNAME"`
	Password               string `env:"PASSWORD"`
	TrustServerCertificate bool   `env:"TRUST_SERVER_CERTIFICATE"  envDefault:"true"`
	Encrypt                string `env:"ENCRYPT"                   envDefault:"disable"`
	MaxIdleConn            int    `env:"MAX_IDLE_CONN"             envDefault:"10"`
	MaxOpenConn            int    `env:"MAX_OPEN_CONN"             envDefault:"100"`
	MaxLifetime            int    `env:"MAX_LIFETIME_MINUTES"      envDefault:"60"`
}

type RedisConfig struct {
	Host     string `env:"HOST"     envDefault:"localhost"`
	Port     int    `env:"PORT"     envDefault:"6379"`
	Password string `env:"PASSWORD"`
}

type MailConfig struct {
	Mailer     string `env:"MAILER"`
	Host       string `env:"HOST"`
	Port       int    `env:"PORT"         envDefault:"587"`
	Username   string `env:"USERNAME"`
	Password   string `env:"PASSWORD"`
	Encryption string `env:"ENCRYPTION"`
	FromAddr   string `env:"FROM_ADDRESS"`
	FromName   string `env:"FROM_NAME"`
}

type PusherConfig struct {
	AppID   string `env:"APP_ID"`
	AppKey  string `env:"APP_KEY"`
	Secret  string `env:"APP_SECRET"`
	Host    string `env:"HOST"`
	Port    int    `env:"PORT"`
	Scheme  string `env:"SCHEME"`
	Cluster string `env:"CLUSTER"`
}

type JWTConfig struct {
	SecretKey       string `env:"SECRET_KEY"`
	AccessTokenExp  int    `env:"ACCESS_TOKEN_EXP_MINUTES"  envDefault:"30"`
	RefreshTokenExp int    `env:"REFRESH_TOKEN_EXP_DAYS"    envDefault:"7"`
}
