package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type legacyUser struct {
	ID       int64  `json:"legacy_user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Verification is server-side and fails closed; no local account exists until it succeeds.
func verifyLegacyUser(email, password string) (*legacyUser, error) {
	endpoint := strings.TrimSpace(os.Getenv("LEGACY_VERIFY_URL"))
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.User != nil {
		return nil, &Problem{503, "Account verification is not configured. Contact your administrator."}
	}
	payload, _ := json.Marshal(map[string]string{"email": email, "password": password})
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, &Problem{503, "Account verification is unavailable. Try again later."}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Do(req)
	if err != nil {
		return nil, &Problem{503, "Account verification is unavailable. Try again later."}
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, &Problem{401, "Your IT Support credentials could not be verified."}
	case http.StatusTooManyRequests:
		return nil, &Problem{429, "Too many verification attempts. Please wait a minute and try again."}
	case http.StatusOK:
	default:
		return nil, &Problem{503, "Account verification is unavailable. Try again later."}
	}
	var result struct {
		Success          bool       `json:"success"`
		PasswordVerified bool       `json:"password_verified"`
		User             legacyUser `json:"user"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 64<<10)).Decode(&result); err != nil {
		return nil, &Problem{503, "Account verification returned an invalid response. Try again later."}
	}
	if !result.Success || !result.PasswordVerified || result.User.ID < 1 || strings.TrimSpace(result.User.Name) == "" || !strings.EqualFold(strings.TrimSpace(result.User.Email), email) {
		return nil, &Problem{401, "Your IT Support credentials could not be verified."}
	}
	result.User.Name = strings.TrimSpace(result.User.Name)
	return &result.User, nil
}
