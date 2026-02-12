package itmotokens

import (
	"context"
	"encoding/json"
	"html"
	"io"
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	nethtml "golang.org/x/net/html"
)

// Get performs OAuth2 Authorization Code Flow with PKCE and returns tokens.
//
//nolint:funlen // OAuth PKCE flow requires many steps.
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
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
	)
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

	formActionURL, err := url.Parse(formAction)
	if err != nil {
		return nil, errors.Wrap(err, "parse form action")
	}
	if !formActionURL.IsAbs() {
		formActionURL = resp.Request.URL.ResolveReference(formActionURL)
	}

	// Step 2: Submit the login form
	form := url.Values{
		"username":   {strconv.FormatInt(isu, 10)},
		"password":   {password},
		"rememberMe": {"on"}, // Match Python implementation
	}

	// Add any session-specific parameters found in the form
	maps.Copy(form, sessionParams)

	formReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		formActionURL.String(),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "build form request")
	}

	formReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	formReq.Header.Set("Origin", c.providerURL)
	formReq.Header.Set("Referer", resp.Request.URL.String())

	for _, cookie := range resp.Cookies() {
		formReq.AddCookie(cookie)
	}

	noRedirectClient := &http.Client{
		Transport: c.httpClient.Transport,
		Timeout:   c.httpClient.Timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	formResp, err := noRedirectClient.Do(formReq)
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
//
//nolint:gocognit,nestif // HTML parsing requires multiple steps.
func (c *Client) extractLoginFormData(htmlContent string) (string, map[string][]string, error) {
	doc, err := nethtml.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", nil, errors.Wrap(err, "parse html")
	}

	var bestForm *loginForm
	var walk func(*nethtml.Node)
	walk = func(n *nethtml.Node) {
		if n.Type == nethtml.ElementNode && n.Data == "form" {
			form := extractLoginForm(n)
			if form != nil {
				if form.hasUsername && form.hasPassword {
					bestForm = form
					return
				}
				if bestForm == nil {
					bestForm = form
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
			if bestForm != nil && bestForm.hasUsername && bestForm.hasPassword {
				return
			}
		}
	}

	walk(doc)

	if bestForm == nil || bestForm.action == "" {
		c.logger.Debug("form", zap.String("html", htmlContent))
		return "", nil, errors.New("form action not found")
	}

	return bestForm.action, bestForm.hiddenFields, nil
}

type loginForm struct {
	action       string
	hiddenFields map[string][]string
	hasUsername  bool
	hasPassword  bool
}

// submitLoginForm submits the login form with the given credentials and session parameters.
//
//nolint:gocognit,nestif // Form submission requires multiple steps.
func extractLoginForm(form *nethtml.Node) *loginForm {
	action := getAttr(form, "action")
	if action == "" {
		return nil
	}

	data := &loginForm{
		action:       html.UnescapeString(action),
		hiddenFields: make(map[string][]string),
	}

	var walk func(*nethtml.Node)
	walk = func(n *nethtml.Node) {
		if n.Type == nethtml.ElementNode && n.Data == "input" {
			name := getAttr(n, "name")
			if name != "" {
				inputType := strings.ToLower(getAttr(n, "type"))
				if name == "username" {
					data.hasUsername = true
				}
				if name == "password" {
					data.hasPassword = true
				}
				if inputType == "hidden" {
					value := html.UnescapeString(getAttr(n, "value"))
					data.hiddenFields[name] = append(data.hiddenFields[name], value)
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(form)

	return data
}

func getAttr(node *nethtml.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
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
	if errDecode := json.NewDecoder(resp.Body).Decode(&tokenData); errDecode != nil {
		return nil, errors.Wrap(errDecode, "decode refresh response")
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
