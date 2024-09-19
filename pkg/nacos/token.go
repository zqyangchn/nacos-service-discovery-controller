package nacos

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Token struct {
	token      string
	expiration time.Time
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenTTL    int    `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
	Username    string `json:"username"`
}

const tokenUri = "/v1/auth/login"

func (n *Nacos) updateAccessToken() error {
	uri := n.config.scheme + schemaDelimiter +
		n.config.iPAddr + portDelimiter + strconv.FormatUint(n.config.port, 10) +
		n.config.contextPath + tokenUri

	data := url.Values{}
	data.Set("username", n.config.username)
	data.Set("password", n.config.password)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	now := time.Now()
	resp, err := n.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("Nacos token response status code: %d, body: %s",
			resp.StatusCode, string(body),
		)
	}

	tokenResponse := TokenResponse{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return err
	}

	n.token = &Token{
		token:      tokenResponse.AccessToken,
		expiration: now.Add(time.Duration(tokenResponse.TokenTTL) * time.Second),
	}

	return nil
}

func (n *Nacos) GetAccessToken() (string, error) {
	if n.token.expiration.Sub(time.Now()) < 5*time.Minute {
		if err := n.updateAccessToken(); err != nil {
			return "", err
		}
	}
	return n.token.token, nil
}
