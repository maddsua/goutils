package env

type AppEnv struct {
	CI    bool
	Prod  bool
	Debug bool
}

func GetAppEnv() AppEnv {
	return AppEnv{
		CI:    Env("CI").IsTrue() && !Env("NO_CI_OVERRIDE").IsTrue(),
		Prod:  !Env("PRODUCTION").IsFalse() && !Env("DEV_MODE").IsTrue(),
		Debug: Env("DEBUG").IsTrue() || Env("DEV_MODE").IsTrue(),
	}
}
