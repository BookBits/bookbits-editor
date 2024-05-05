package models

type EnvVars struct {
	Port string
	DbPort string
	DbHost string
	DbUser string
	DbPassword string
	DbName string
	JWTSecretKey []byte
}
