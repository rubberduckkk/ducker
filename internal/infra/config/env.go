package config

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

func IsProd() bool {
	return config.Env == EnvProd
}
