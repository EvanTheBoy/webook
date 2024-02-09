package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(192.168.183.134:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "192.168.183.134:6379",
	},
}
