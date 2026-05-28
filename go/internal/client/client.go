package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/vngcloud/greennode-cli/internal/auth"
)

const (
	maxRetries     = 3
	retryBaseDelay = 1 * time.Second
	defaultTimeout = 30 * time.Second
)

var statusMessages = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	409: "Conflict",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

var retryableStatusCodes = map[int]bool{
	500: true, 502: true, 503: true, 504: true,
}

// GreenNodeClient is an HTTP client for greennode APIs with retry and auto token refresh.
type GreenNodeClient struct {
	baseURL      string
	tokenManager *auth.TokenManager
	httpClient   *http.Client
	debug        bool
}

// NewGreenNodeClient creates a new API client.
func NewGreenNodeClient(baseURL string, tokenManager *auth.TokenManager, timeout time.Duration, verifySSL bool, debug bool) *GreenNodeClient {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	transport := &http.Transport{}
	if !verifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}

	return &GreenNodeClient{
		baseURL:      baseURL,
		tokenManager: tokenManager,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		debug: debug,
	}
}

// Get performs a GET request.
func (c *GreenNodeClient) Get(path string, params map[string]string) (interface{}, error) {
	return c.request("GET", path, params, nil)
}

// Post performs a POST request with a JSON body.
func (c *GreenNodeClient) Post(path string, body interface{}) (interface{}, error) {
	return c.request("POST", path, nil, body)
}

// Put performs a PUT request with a JSON body.
func (c *GreenNodeClient) Put(path string, body interface{}) (interface{}, error) {
	return c.request("PUT", path, nil, body)
}

// Delete performs a DELETE request with optional query params.
func (c *GreenNodeClient) Delete(path string, params map[string]string) (interface{}, error) {
	return c.request("DELETE", path, params, nil)
}

// DeleteWithBody performs a DELETE request with a JSON body.
func (c *GreenNodeClient) DeleteWithBody(path string, body interface{}) (interface{}, error) {
	return c.request("DELETE", path, nil, body)
}

// GetRaw performs a GET request and returns the raw response body.
func (c *GreenNodeClient) GetRaw(path string, params map[string]string) (string, error) {
	return c.requestRaw("GET", path, params, nil)
}

// GetAllPages fetches all pages and merges items into a single result.
func (c *GreenNodeClient) GetAllPages(path string, pageSize int) (map[string]interface{}, error) {
	if pageSize == 0 {
		pageSize = 50
	}

	var allItems []interface{}
	page := 0

	for {
		params := map[string]string{
			"page":     fmt.Sprintf("%d", page),
			"pageSize": fmt.Sprintf("%d", pageSize),
		}
		result, err := c.Get(path, params)
		if err != nil {
			return nil, err
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			break
		}

		items, _ := resultMap["items"].([]interface{})
		allItems = append(allItems, items...)

		total, _ := resultMap["total"].(float64)
		if len(allItems) >= int(total) || len(items) == 0 {
			break
		}
		page++
	}

	return map[string]interface{}{
		"items": allItems,
		"total": float64(len(allItems)),
	}, nil
}

func (c *GreenNodeClient) request(method, path string, params map[string]string, body interface{}) (interface{}, error) {
	rawBody, err := c.requestRaw(method, path, params, body)
	if err != nil {
		return nil, err
	}

	if rawBody == "" {
		return map[string]interface{}{}, nil
	}

	var result interface{}
	if err := json.Unmarshal([]byte(rawBody), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return result, nil
}

func (c *GreenNodeClient) requestRaw(method, path string, params map[string]string, body interface{}) (string, error) {
	fullURL := c.baseURL + path

	if len(params) > 0 {
		u, _ := url.Parse(fullURL)
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	token, err := c.tokenManager.GetToken()
	if err != nil {
		return "", err
	}

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		var reqBody io.Reader
		var jsonBody []byte
		if body != nil {
			jsonBody, err = json.Marshal(body)
			if err != nil {
				return "", fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewReader(jsonBody)
		}

		if c.debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, fullURL)
			if len(jsonBody) > 0 {
				fmt.Fprintf(os.Stderr, "[DEBUG] Request body: %s\n", jsonBody)
			}
		}

		req, err := http.NewRequest(method, fullURL, reqBody)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				delay := retryBaseDelay * time.Duration(1<<uint(attempt))
				time.Sleep(delay)
				continue
			}
			return "", fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, lastErr)
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if c.debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Response status: %d\n", resp.StatusCode)
			fmt.Fprintf(os.Stderr, "[DEBUG] Response body: %s\n", respBody)
		}

		// 401 — refresh token and retry once
		if resp.StatusCode == http.StatusUnauthorized {
			token, err = c.tokenManager.RefreshToken()
			if err != nil {
				return "", err
			}
			var reqBody2 io.Reader
			if len(jsonBody) > 0 {
				reqBody2 = bytes.NewReader(jsonBody)
			}
			req2, _ := http.NewRequest(method, fullURL, reqBody2)
			req2.Header.Set("Authorization", "Bearer "+token)
			req2.Header.Set("Content-Type", "application/json")
			resp2, err := c.httpClient.Do(req2)
			if err != nil {
				return "", err
			}
			respBody, _ = io.ReadAll(resp2.Body)
			resp2.Body.Close()
			resp = resp2
		}

		// Retryable server errors (5xx)
		if retryableStatusCodes[resp.StatusCode] {
			if attempt < maxRetries {
				delay := retryBaseDelay * time.Duration(1<<uint(attempt))
				time.Sleep(delay)
				continue
			}
		}

		// Non-retryable errors
		if resp.StatusCode >= 400 {
			return "", fmt.Errorf("%s", formatError(resp.StatusCode, respBody))
		}

		return string(respBody), nil
	}

	return "", fmt.Errorf("request failed after %d attempts", maxRetries+1)
}

func formatError(statusCode int, body []byte) string {
	statusText := statusMessages[statusCode]
	if statusText == "" {
		statusText = "Error"
	}

	detail := ""
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		if msg, ok := data["message"].(string); ok && msg != "" {
			detail = msg
		} else if errMsg, ok := data["error"].(string); ok && errMsg != "" {
			detail = errMsg
		} else if detailMsg, ok := data["detail"].(string); ok && detailMsg != "" {
			detail = detailMsg
		} else if errors, ok := data["errors"].([]interface{}); ok && len(errors) > 0 {
			if errObj, ok := errors[0].(map[string]interface{}); ok {
				if msg, ok := errObj["message"].(string); ok {
					detail = msg
				}
			}
		}
		if detail == "" {
			detail = string(body)
		}
	} else {
		detail = string(body)
	}

	if detail != "" {
		return fmt.Sprintf("API error (HTTP %d %s): %s", statusCode, statusText, detail)
	}
	return fmt.Sprintf("API error (HTTP %d %s)", statusCode, statusText)
}
