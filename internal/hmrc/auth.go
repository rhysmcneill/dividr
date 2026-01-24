package hmrc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rhysmcneill/dividr/internal/config"
	"golang.org/x/oauth2"
)

// Constants for HMRC Environments
const (
	BaseURLSandbox = "https://test-api.service.hmrc.gov.uk"
	BaseURLProd    = "https://api.service.hmrc.gov.uk"
)

var scopes = []string{
	"read:self-assessment",
	"write:self-assessment",
	"read:quarterly-obligations",
	"read:business-details",
	"hello", // Uncomment if you want to test the Hello World endpoint
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	Scope        string    `json:"scope"`
	Expiry       time.Time `json:"-"`
}

type AuthService struct {
	config *oauth2.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	// 1. Determine Environment
	env := cfg.AppEnv
	baseURL := BaseURLSandbox
	if env == "production" {
		baseURL = BaseURLProd
	}

	// 2. Get Redirect URL from Env (Crucial for Local vs Prod)
	redirectURL := cfg.HMRCRedirectURL
	if redirectURL == "" {
		// Fallback for safety, though you should always set the ENV var
		redirectURL = "http://localhost:8080/auth/hmrc/callback"
	}

	return &AuthService{
		config: &oauth2.Config{
			ClientID:     cfg.HMRCClientID,
			ClientSecret: cfg.HMRCClientSecret,
			Scopes:       scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  baseURL + "/oauth/authorize",
				TokenURL: baseURL + "/oauth/token",
			},
			RedirectURL: redirectURL,
		},
	}
}

// GenerateAuthURL creates the link to send the user to HMRC
func (s *AuthService) GenerateAuthURL() (string, string) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand.Read only fails if the system's RNG is broken
		// Fall back to timestamp-based state (less secure but functional)
		b = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}
	state := hex.EncodeToString(b)
	return s.config.AuthCodeURL(state), state
}

// ExchangeCode swaps the temporary code for permanent tokens
func (s *AuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

// RefreshAccessToken uses the stored refresh token to get a new access token
func (s *AuthService) RefreshAccessToken(refreshToken string) (*TokenResponse, error) { // 1. Determine Base URL
	endpoint := s.config.Endpoint.TokenURL

	// 2. Prepare Form Data
	data := url.Values{}
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// 3. Create Request
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 4. Execute
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			slog.Warn("failed to close response body", "error", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to refresh token: status %d, body: %s", resp.StatusCode, string(body))
	}

	// 5. Parse Response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Calculate absolute expiry time immediately
	// HMRC returns "expires_in" (seconds), usually 14400 (4 hours)
	tokenResp.Expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return &tokenResp, nil
}
