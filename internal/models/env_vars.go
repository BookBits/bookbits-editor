package models

type EnvVars struct {
	Port string
	DbPort string
	DbHost string
	DbUser string
	DbPassword string
	DbName string
	JWTSecretKey []byte
	DefaultAdminUserEmail string
	DefaultAdminPassword string
	GitRepo string
	GitOwner string
	GitToken string
	RedisPort string
	RedisAddr string
	RedisPassword string
}
