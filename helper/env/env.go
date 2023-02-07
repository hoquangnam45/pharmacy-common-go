package env

import "os"

func GetEnvOrDefault(key, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return env
}

func GetEnvOrDefaultT[T any](key string, convFn func(string) (T, error), defaultValue T) T {
	env, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	if val, err := convFn(env); err != nil {
		return defaultValue
	} else {
		return val
	}
}