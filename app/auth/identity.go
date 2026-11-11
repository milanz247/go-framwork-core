package auth

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"ticket-sync-mvc/app/repositories"
	"ticket-sync-mvc/framework/security"
	"ticket-sync-mvc/framework/support"
)

type Problem struct {
	Status  int
	Message string
}

func (p *Problem) Error() string { return p.Message }
func Authenticate(v support.M, register bool) (int64, error) {
	email := strings.ToLower(strings.TrimSpace(support.Str(v["email"])))
	pass := support.Str(v["password"])
	if register {
		if !strings.Contains(email, "@") || pass == "" {
			return 0, &Problem{400, "Enter your existing IT Support email and password"}
		}
		if repositories.UserByEmail(email) != nil {
			return 0, &Problem{409, "Email is already registered"}
		}
		verified, e := verifyLegacyUser(email, pass)
		if e != nil {
			return 0, e
		}
		hash, e := bcrypt.GenerateFromPassword([]byte(pass), 12)
		if e != nil {
			return 0, &Problem{400, "Password exceeds the supported length"}
		}
		return repositories.CreateUser(verified.Name, strings.ToLower(strings.TrimSpace(verified.Email)), string(hash))
	}
	u := repositories.UserByEmail(email)
	if u == nil || !security.PasswordOK(support.Str(u["password_hash"]), pass) {
		return 0, &Problem{401, "Email or password is incorrect"}
	}
	id := support.Number(u["id"])
	if id < 1 {
		return 0, errors.New("invalid user")
	}
	return id, nil
}
