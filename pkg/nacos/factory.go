package nacos

import (
	"net/http"
)

const (
	schemaDelimiter = "://"
	portDelimiter   = ":"
)

type Config struct {
	iPAddr      string
	scheme      string
	port        uint64
	username    string
	password    string
	contextPath string
	namespaceId string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) SetIpAddr(ip string) *Config {
	c.iPAddr = ip
	return c
}

func (c *Config) SetScheme(scheme string) *Config {
	c.scheme = scheme
	return c
}

func (c *Config) SetPort(port uint64) *Config {
	c.port = port
	return c
}

func (c *Config) SetUsername(username string) *Config {
	c.username = username
	return c
}

func (c *Config) SetPassword(password string) *Config {
	c.password = password
	return c
}

func (c *Config) SetContextPath(contextPath string) *Config {
	c.contextPath = contextPath
	return c
}

func (c *Config) SetNamespaceId(namespaceId string) *Config {
	c.namespaceId = namespaceId
	return c
}

type Nacos struct {
	config *Config

	hc    *http.Client
	token *Token
}

func New(config *Config) (*Nacos, error) {
	nacos := Nacos{
		config: config,
		hc:     &http.Client{},
	}

	if err := nacos.updateAccessToken(); err != nil {
		return nil, err
	}

	return &nacos, nil
}
