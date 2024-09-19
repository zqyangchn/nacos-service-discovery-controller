package nacos

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

type NamespacesResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []Namespace `json:"data"`
}

type Namespace struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	NamespaceDesc     string `json:"namespaceDesc"`
	Quota             int    `json:"quota"`
	ConfigCount       int    `json:"configCount"`
	Type              int    `json:"type"`
}

type GetNamespacesParam struct {
	AccessToken string `url:"accessToken"`
}

const namespacesUri = "/v1/console/namespaces"

func (n *Nacos) GetNamespaces(getNamespacesParam GetNamespacesParam) ([]Namespace, error) {
	uri := n.config.scheme + schemaDelimiter +
		n.config.iPAddr + portDelimiter + strconv.FormatUint(n.config.port, 10) +
		n.config.contextPath + namespacesUri

	token, err := n.GetAccessToken()
	if err != nil {
		return nil, err
	}
	getNamespacesParam.AccessToken = token
	params, err := query.Values(getNamespacesParam)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"Get namespace failed, status code: %d, body: %s", resp.StatusCode, body,
		)
	}

	namespacesResponse := &NamespacesResponse{}
	if err := json.Unmarshal(body, namespacesResponse); err != nil {
		return nil, err
	}

	return namespacesResponse.Data, nil
}
