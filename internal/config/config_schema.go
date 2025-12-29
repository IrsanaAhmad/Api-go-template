package config

type Config struct {
	App struct {
		Name  string `json:"name"`
		Env   string `json:"env"`
		Key   string `json:"key"`
		Debug bool   `json:"debug"`
		URL   string `json:"url"`
		Base  string `json:"base_url"`
	} `json:"app"`

	Log struct {
		Channel             string `json:"channel"`
		DeprecationsChannel any    `json:"deprecations_channel"`
		Level               string `json:"level"`
	} `json:"log"`

	Database struct {
		Connection             string `json:"connection"`
		Host                   string `json:"host"`
		Port                   int    `json:"port"`
		DBMain                 string `json:"db_main"`
		DBSecondary            string `json:"db_secondary"`
		DBThird                string `json:"db_third"`
		Username               string `json:"username"`
		Password               string `json:"password"`
		TrustServerCertificate bool   `json:"trust_server_certificate"`
		Encrypt                string `json:"encrypt"`
		MaxIdleConn            int    `json:"max_idle_conn"`
		MaxOpenConn            int    `json:"max_open_conn"`
		MaxLifetime            int    `json:"max_lifetime_minutes"`
	} `json:"database"`

	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password any    `json:"password"`
	} `json:"redis"`

	Mail struct {
		Mailer     string `json:"mailer"`
		Host       string `json:"host"`
		Port       int    `json:"port"`
		Username   any    `json:"username"`
		Password   any    `json:"password"`
		Encryption any    `json:"encryption"`
		FromAddr   string `json:"from_address"`
		FromName   string `json:"from_name"`
	} `json:"mail"`

	Pusher struct {
		AppID   string `json:"app_id"`
		AppKey  string `json:"app_key"`
		Secret  string `json:"app_secret"`
		Host    string `json:"host"`
		Port    int    `json:"port"`
		Scheme  string `json:"scheme"`
		Cluster string `json:"cluster"`
	} `json:"pusher"`

	JWT struct {
		SecretKey       string `json:"secret_key"`
		AccessTokenExp  int    `json:"access_token_exp_minutes"`
		RefreshTokenExp int    `json:"refresh_token_exp_days"`
	} `json:"jwt"`
}
