package env

import "os"

func Get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
