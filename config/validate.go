package config

import (
	"fmt"
	"net/url"
	"strconv"

	"ticket-sync-mvc/framework/env"
)

func Validate() error {
	port, e := strconv.Atoi(env.Get("PORT", "8001"))
	if e != nil || port < 1 || port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535")
	}
	for _, k := range []string{"COOKIE_SECURE"} {
		v := env.Get(k, "0")
		if v != "0" && v != "1" {
			return fmt.Errorf("%s must be 0 or 1", k)
		}
	}
	if v := env.Get("PUBLIC_ORIGIN", ""); v != "" {
		u, e := url.Parse(v)
		if e != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
			return fmt.Errorf("PUBLIC_ORIGIN must be an http(s) origin without path or credentials")
		}
	}
	return nil
}
