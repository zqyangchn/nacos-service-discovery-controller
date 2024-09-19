package nacos

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	servicePageNo   = 1
	servicePageSize = 200
)

type ServiceResponse struct {
	Count int      `json:"count"`
	Doms  []string `json:"doms"`
}

type GetServiceParam struct {
	AccessToken string `url:"accessToken"`

	NamespaceId string `url:"namespaceId"`

	PageNo   int64 `url:"pageNo"`
	PageSize int64 `url:"pageSize"`
}

const serviceUri = "/v1/ns/service/list"

func (n *Nacos) getService(getServiceParam GetServiceParam) (*ServiceResponse, error) {
	uri := n.config.scheme + schemaDelimiter +
		n.config.iPAddr + portDelimiter + strconv.FormatUint(n.config.port, 10) +
		n.config.contextPath + serviceUri

	token, err := n.GetAccessToken()
	if err != nil {
		return nil, err
	}
	getServiceParam.AccessToken = token
	params, err := query.Values(getServiceParam)
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
			"Get service failed, status code: %d, body: %s", resp.StatusCode, body,
		)
	}

	serviceResponse := &ServiceResponse{}
	if err := json.Unmarshal(body, serviceResponse); err != nil {
		return nil, err
	}

	return serviceResponse, nil
}

func (n *Nacos) GetService(NamespaceId string) ([]string, error) {
	token, err := n.GetAccessToken()
	if err != nil {
		return nil, err
	}

	services := make([]string, 0)

	pageNo, total := int64(servicePageNo), -1
	param := GetServiceParam{
		AccessToken: token,
		NamespaceId: NamespaceId,
		PageNo:      pageNo,
		PageSize:    servicePageSize,
	}

GETMORESERVICE:
	resp, err := n.getService(param)
	if err != nil {
		return nil, err
	}
	if total == -1 {
		total = resp.Count
	}

	services = append(services, resp.Doms...)
	if pageNo*servicePageSize < int64(total) {
		pageNo++
		param.PageNo = pageNo
		goto GETMORESERVICE
	}

	return services, nil
}
