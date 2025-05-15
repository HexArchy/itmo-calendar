package itmotokens

import (
	"context"
	"encoding/json"
	"html"
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

// Get performs OAuth2 Authorization Code Flow with PKCE and returns tokens.
func (c *Client) Get(ctx context.Context, isu int64, password string) (*entities.UserTokens, error) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, errors.Wrap(err, "generate code verifier")
	}
	codeChallenge := getCodeChallenge(codeVerifier)

	// Step 1: Get the login page
	authURL := c.providerURL + "/protocol/openid-connect/auth"
	params := url.Values{
		"protocol":              {"oauth2"},
		"response_type":         {"code"},
		"client_id":             {c.clientID},
		"redirect_uri":          {c.redirectURI},
		"scope":                 {"openid"},
		"state":                 {"im_not_a_browser"},
		"code_challenge_method": {"S256"},
		"code_challenge":        {codeChallenge},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "build auth request")
	}

	// Set request headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "auth request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read auth response")
	}

	// Extract both the form action URL and session-related parameters
	formAction, sessionParams, err := c.extractLoginFormData(string(body))
	if err != nil {
		return nil, errors.Wrap(err, "extract form data")
	}

	// Step 2: Submit the login form
	form := url.Values{
		"username":   {strconv.FormatInt(isu, 10)},
		"password":   {password},
		"rememberMe": {"on"}, // Match Python implementation
	}

	// Add any session-specific parameters found in the form
	for k, v := range sessionParams {
		form[k] = v
	}

	formReq, err := http.NewRequestWithContext(ctx, http.MethodPost, formAction, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "build form request")
	}

	formReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	formReq.Header.Set("Origin", c.providerURL)
	formReq.Header.Set("Referer", resp.Request.URL.String())

	for _, cookie := range resp.Cookies() {
		formReq.AddCookie(cookie)
	}

	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		c.httpClient.CheckRedirect = nil
	}()

	formResp, err := c.httpClient.Do(formReq)
	if err != nil {
		return nil, errors.Wrap(err, "form submit")
	}
	defer formResp.Body.Close()

	if formResp.StatusCode != http.StatusFound {
		bodyBytes, _ := io.ReadAll(formResp.Body)
		bodyStr := string(bodyBytes)

		// Extract error message if present
		errorPattern := regexp.MustCompile(`<span[^>]*class="[^"]*invalid-feedback[^"]*"[^>]*>(.*?)</span>`)
		if matches := errorPattern.FindStringSubmatch(bodyStr); len(matches) > 1 {
			return nil, errors.Errorf("authentication failed: %s", strings.TrimSpace(matches[1]))
		}

		return nil, errors.Errorf("unexpected form response: %d", formResp.StatusCode)
	}

	// Step 3: Handle the redirect and extract the authorization code
	loc := formResp.Header.Get("Location")
	u, err := url.Parse(loc)
	if err != nil {
		return nil, errors.Wrap(err, "parse redirect location")
	}

	code := u.Query().Get("code")
	if code == "" {
		return nil, errors.New("no code in redirect")
	}

	// Step 4: Exchange the code for tokens
	tokenPair, err := c.exchangeCode(ctx, code, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "exchange code")
	}

	tokenPair.ISU = isu

	return tokenPair, nil
}

// exchangeCode exchanges authorization code for tokens.
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
		return nil, errors.Wrap(err, "create token request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send token request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("unexpected token response: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, errors.Wrap(err, "decode token response")
	}

	// Calculate expiration times
	now := time.Now()
	accessExpires := now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	refreshExpires := now.Add(time.Duration(tokenResp.RefreshExpiresIn) * time.Second)

	return &entities.UserTokens{
		AccessToken:           tokenResp.AccessToken,
		RefreshToken:          tokenResp.RefreshToken,
		AccessTokenExpiresAt:  accessExpires,
		RefreshTokenExpiresAt: refreshExpires,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

// extractLoginFormData extracts both the form action URL and hidden form fields.
func (c *Client) extractLoginFormData(htmlContent string) (string, map[string][]string, error) {
	formActionRe := regexp.MustCompile(`(?s)<form[^>]*\s+id="kc-form-login"[^>]*\s+action="([^"]+)"`)
	matches := formActionRe.FindStringSubmatch(htmlContent)
	if len(matches) < 2 {
		c.logger.Debug("form", zap.String("html", htmlContent))
		return "", nil, errors.New("form action not found")
	}

	formAction := html.UnescapeString(matches[1])

	hiddenFieldsRe := regexp.MustCompile(`<input[^>]+type="hidden"[^>]+name="([^"]+)"[^>]+value="([^"]*)"[^>]*>`)
	paramMatches := hiddenFieldsRe.FindAllStringSubmatch(htmlContent, -1)

	params := make(map[string][]string)
	for _, match := range paramMatches {
		if len(match) >= 3 {
			params[match[1]] = []string{html.UnescapeString(match[2])}
		}
	}

	return formAction, params, nil
}

// Refresh exchanges a refresh token for a new access token.
func (c *Client) Refresh(ctx context.Context, isu int64, refreshToken string) (*entities.UserTokens, error) {
	tokenURL := c.providerURL + "/protocol/openid-connect/token"
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {c.clientID},
		"refresh_token": {refreshToken},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "build refresh request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "refresh request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("unexpected refresh response: %d %s", resp.StatusCode, string(b))
	}
	var tokenData struct {
		AccessToken           string `json:"access_token"`
		RefreshToken          string `json:"refresh_token"`
		AccessTokenExpiresIn  int64  `json:"expires_in"`
		RefreshTokenExpiresIn int64  `json:"refresh_expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return nil, errors.Wrap(err, "decode refresh response")
	}

	now := time.Now()
	return &entities.UserTokens{
		ISU:                   isu,
		AccessToken:           tokenData.AccessToken,
		RefreshToken:          tokenData.RefreshToken,
		AccessTokenExpiresAt:  now.Add(time.Duration(tokenData.AccessTokenExpiresIn) * time.Second),
		RefreshTokenExpiresAt: now.Add(time.Duration(tokenData.RefreshTokenExpiresIn) * time.Second),
	}, nil
}
