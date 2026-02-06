package config

var Config = config{
	DB: DBConfig{
		DSN: "root:123456@tcp(127.0.0.1:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "127.0.0.1:16379",
	},
}
