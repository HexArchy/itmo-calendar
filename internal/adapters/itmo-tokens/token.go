package itmotokens

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// loginAction extraction patterns.
var (
	// kcContextLoginActionRe extracts loginAction from Keycloakify SPA's kcContext JS object.
	// Example: "loginAction": "https://id.itmo.ru/auth/realms/itmo/login-actions/authenticate?..."
	kcContextLoginActionRe = regexp.MustCompile(`"loginAction"\s*:\s*"(https?:[^"]+)"`)

	// htmlFormActionRe extracts action from a classic HTML <form> element (legacy Keycloak themes).
	htmlFormActionRe = regexp.MustCompile(`<form[^>]+action="([^"]+)"`)
)

// Get performs OAuth2 Authorization Code Flow with PKCE and returns tokens.
func (c *Client) Get(ctx context.Context, isu int64, password string) (*entities.UserTokens, error) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, errors.Wrap(err, "generate code verifier")
	}

	// Step 1: Load the login page and extract the form action URL.
	loginAction, cookies, err := c.fetchLoginAction(ctx, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "fetch login page")
	}

	// Step 2: Submit credentials and get the authorization code.
	code, err := c.submitCredentials(ctx, loginAction, cookies, isu, password)
	if err != nil {
		return nil, err
	}

	// Step 3: Exchange the authorization code for tokens.
	tokens, err := c.exchangeCode(ctx, code, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "exchange code")
	}

	tokens.ISU = isu

	return tokens, nil
}

// fetchLoginAction loads the Keycloak login page and extracts the URL to POST credentials to.
func (c *Client) fetchLoginAction(ctx context.Context, codeVerifier string) (string, []*http.Cookie, error) {
	authURL := c.providerURL + "/protocol/openid-connect/auth"
	params := url.Values{
		"protocol":              {"oauth2"},
		"response_type":         {"code"},
		"client_id":             {c.clientID},
		"redirect_uri":          {c.redirectURI},
		"scope":                 {"openid"},
		"state":                 {"im_not_a_browser"},
		"code_challenge_method": {"S256"},
		"code_challenge":        {getCodeChallenge(codeVerifier)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", nil, errors.Wrap(err, "build request")
	}

	setBrowserHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, errors.Wrap(err, "read response")
	}

	loginAction := extractLoginAction(string(body))
	if loginAction == "" {
		c.logger.Warn("login page has no loginAction",
			zap.Int("status", resp.StatusCode),
			zap.String("url", resp.Request.URL.String()),
		)
		c.logger.Debug("login page body", zap.String("html", string(body)))

		return "", nil, errors.New("loginAction not found in login page")
	}

	// Resolve relative URLs against the response URL.
	parsed, err := url.Parse(loginAction)
	if err != nil {
		return "", nil, errors.Wrap(err, "parse loginAction URL")
	}

	if !parsed.IsAbs() {
		loginAction = resp.Request.URL.ResolveReference(parsed).String()
	}

	return loginAction, resp.Cookies(), nil
}

// extractLoginAction extracts the login form action URL from a Keycloak page.
// Supports both Keycloakify SPA (kcContext JS) and classic server-rendered HTML forms.
func extractLoginAction(body string) string {
	// Primary: Keycloakify SPA — loginAction in kcContext JS object.
	if m := kcContextLoginActionRe.FindStringSubmatch(body); m != nil {
		// Keycloakify escapes forward slashes in JS string literals.
		return strings.ReplaceAll(m[1], `\/`, `/`)
	}

	// Fallback: classic HTML form with action attribute.
	if m := htmlFormActionRe.FindStringSubmatch(body); m != nil {
		return strings.ReplaceAll(m[1], "&amp;", "&")
	}

	return ""
}

// submitCredentials POSTs username/password to the login action URL.
// On success, Keycloak responds with 302 redirect containing the authorization code.
func (c *Client) submitCredentials(
	ctx context.Context,
	loginAction string,
	cookies []*http.Cookie,
	isu int64,
	password string,
) (string, error) {
	form := url.Values{
		"username":   {strconv.FormatInt(isu, 10)},
		"password":   {password},
		"rememberMe": {"on"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginAction, strings.NewReader(form.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "build login request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Don't follow redirects — we need the Location header with the auth code.
	noRedirect := &http.Client{
		Transport: c.httpClient.Transport,
		Timeout:   c.httpClient.Timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := noRedirect.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "send login request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", errors.Errorf("authentication failed (status %d)", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")

	u, err := url.Parse(loc)
	if err != nil {
		return "", errors.Wrap(err, "parse redirect location")
	}

	code := u.Query().Get("code")
	if code == "" {
		return "", errors.Errorf("no authorization code in redirect: %s", loc)
	}

	return code, nil
}

// exchangeCode exchanges an authorization code for access and refresh tokens.
func (c *Client) exchangeCode(ctx context.Context, code, codeVerifier string) (*entities.UserTokens, error) {
	tokenURL := c.providerURL + "/protocol/openid-connect/token"

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {c.clientID},
		"code":          {code},
		"redirect_uri":  {c.redirectURI},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("token exchange failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int64  `json:"expires_in"`
		RefreshExpiresIn int64  `json:"refresh_expires_in"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, errors.Wrap(err, "decode token response")
	}

	now := time.Now()

	return &entities.UserTokens{
		AccessToken:           tokenResp.AccessToken,
		RefreshToken:          tokenResp.RefreshToken,
		AccessTokenExpiresAt:  now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		RefreshTokenExpiresAt: now.Add(time.Duration(tokenResp.RefreshExpiresIn) * time.Second),
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

// Refresh exchanges a refresh token for a new token pair.
func (c *Client) Refresh(ctx context.Context, isu int64, refreshToken string) (*entities.UserTokens, error) {
	tokenURL := c.providerURL + "/protocol/openid-connect/token"

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {c.clientID},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("refresh failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int64  `json:"expires_in"`
		RefreshExpiresIn int64  `json:"refresh_expires_in"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	now := time.Now()

	return &entities.UserTokens{
		ISU:                   isu,
		AccessToken:           tokenResp.AccessToken,
		RefreshToken:          tokenResp.RefreshToken,
		AccessTokenExpiresAt:  now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		RefreshTokenExpiresAt: now.Add(time.Duration(tokenResp.RefreshExpiresIn) * time.Second),
	}, nil
}

func setBrowserHeaders(req *http.Request) {
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	)
	req.Header.Set(
		"Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	)
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
}
