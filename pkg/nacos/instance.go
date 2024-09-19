package nacos

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/logger"
)

type UpdateInstanceParam struct {
	AccessToken string `url:"accessToken"`

	Enable bool `url:"enabled"` //required,the instance can be access or not

	Ip          string `url:"ip"`          //required
	Port        uint64 `url:"port"`        //required
	ServiceName string `url:"serviceName"` //required
	NamespaceId string `url:"namespaceId"` //required

	Weight      float64 `url:"weight"`      //required,it must be lager than 0
	Healthy     bool    `url:"healthy"`     //required,the instance is health or not
	Ephemeral   bool    `url:"ephemeral"`   //optional
	ClusterName string  `url:"clusterName"` //optional
	GroupName   string  `url:"groupName"`   //optional,default:DEFAULT_GROUP
	Metadata    string  `url:"metadata"`    //optional
}

const instanceUri = "/v1/ns/instance"

func (n *Nacos) UpdateInstance(updateInstanceParam UpdateInstanceParam) error {
	uri := n.config.scheme + schemaDelimiter +
		n.config.iPAddr + portDelimiter + strconv.FormatUint(n.config.port, 10) +
		n.config.contextPath + instanceUri

	token, err := n.GetAccessToken()
	if err != nil {
		return err
	}
	updateInstanceParam.AccessToken = token
	params, err := query.Values(updateInstanceParam)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf(
			"Update instance failed, status code: %d, body: %s", resp.StatusCode, body,
		)
	}

	return nil
}

func (n *Nacos) RetryUpdateInstance(attempts int, sleep time.Duration, param UpdateInstanceParam) error {
	var err error

	for i := 0; i < attempts; i++ {
		if err = n.UpdateInstance(param); err == nil {
			return nil
		}
		logger.Debug("RetryUpdateInstance 失败",
			zap.Int("Attempt", i+1), zap.Duration("时间间隔", sleep),
		)
		time.Sleep(sleep)
	}

	return errors.Errorf("after %d attempts, last error: %s", attempts, err)
}

type GetInstanceParam struct {
	AccessToken string `url:"accessToken"`

	ServiceName string `url:"serviceName"` //required
	Ip          string `url:"ip"`          //required
	Port        uint64 `url:"port"`        //required
	NamespaceId string `url:"namespaceId"` //required
	HealthyOnly bool   `url:"healthyOnly"` //required,the instance is health or not
	ClusterName string `url:"cluster"`     //optional
	GroupName   string `url:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool   `url:"ephemeral"`   //optional
}

type Instance struct {
	Service     string            `json:"service"`
	IP          string            `json:"ip"`
	Port        uint64            `json:"port"`
	ClusterName string            `json:"clusterName"`
	Weight      float64           `json:"weight"`
	Healthy     bool              `json:"healthy"`
	InstanceID  string            `json:"instanceId"`
	Metadata    map[string]string `json:"metadata"`
}

func (n *Nacos) GetInstance(getInstanceParam GetInstanceParam) (*Instance, error) {
	uri := n.config.scheme + schemaDelimiter +
		n.config.iPAddr + portDelimiter + strconv.FormatUint(n.config.port, 10) +
		n.config.contextPath + instanceUri

	token, err := n.GetAccessToken()
	if err != nil {
		return nil, err
	}
	getInstanceParam.AccessToken = token
	params, err := query.Values(getInstanceParam)
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
			"Get instance failed, status code: %d, body: %s", resp.StatusCode, body,
		)
	}

	instance := &Instance{}
	if err := json.Unmarshal(body, instance); err != nil {
		return nil, err
	}

	return instance, nil
}

func (n *Nacos) RetryGetInstance(attempts int, sleep time.Duration, param GetInstanceParam) (*Instance, error) {
	var err error
	var instance *Instance

	for i := 0; i < attempts; i++ {
		instance, err = n.GetInstance(param)
		if err == nil {
			return instance, nil
		}
		logger.Debug("RetryGetInstance 失败",
			zap.Int("Attempt", i+1), zap.Duration("时间间隔", sleep),
		)
		time.Sleep(sleep)
	}

	return nil, errors.Errorf("after %d attempts, last error: %s", attempts, err)
}

type ListInstanceParam struct {
	AccessToken string `url:"accessToken"`

	ServiceName string `url:"serviceName"` //required
	NamespaceId string `url:"namespaceId"` //optional
	HealthyOnly bool   `url:"healthyOnly"` //optional,the instance is health or not
	ClusterName string `url:"cluster"`     //optional
	GroupName   string `url:"groupName"`   //optional,default:DEFAULT_GROUP
}

type ListInstanceResponse struct {
	Name        string `json:"name"`
	GroupName   string `json:"groupName"`
	Clusters    string `json:"clusters"`
	CacheMillis int    `json:"cacheMillis"`
	Hosts       []struct {
		IP          string  `json:"ip"`
		Port        int     `json:"port"`
		Weight      float64 `json:"weight"`
		Healthy     bool    `json:"healthy"`
		Enabled     bool    `json:"enabled"`
		Ephemeral   bool    `json:"ephemeral"`
		ClusterName string  `json:"clusterName"`
		ServiceName string  `json:"serviceName"`
		Metadata    struct {
			PreservedRegisterSource string `json:"preserved.register.source"`
		} `json:"metadata"`
		InstanceHeartBeatInterval int `json:"instanceHeartBeatInterval"`
		InstanceHeartBeatTimeOut  int `json:"instanceHeartBeatTimeOut"`
		IPDeleteTimeout           int `json:"ipDeleteTimeout"`
	} `json:"hosts"`
	LastRefTime              int64  `json:"lastRefTime"`
	Checksum                 string `json:"checksum"`
	AllIPs                   bool   `json:"allIPs"`
	ReachProtectionThreshold bool   `json:"reachProtectionThreshold"`
	Valid                    bool   `json:"valid"`
}

const instanceListUri = "/v1/ns/instance/list"

func (n *Nacos) ListInstance(listInstanceParam ListInstanceParam) (*ListInstanceResponse, error) {
	uri := n.config.scheme + schemaDelimiter + n.config.iPAddr +
		portDelimiter + strconv.FormatUint(n.config.port, 10) + n.config.contextPath + instanceListUri

	token, err := n.GetAccessToken()
	if err != nil {
		return nil, err
	}
	listInstanceParam.AccessToken = token
	params, err := query.Values(listInstanceParam)
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
			"Get instance failed, status code: %d, body: %s", resp.StatusCode, body,
		)
	}

	listInstanceResponse := &ListInstanceResponse{}
	if err := json.Unmarshal(body, listInstanceResponse); err != nil {
		return nil, err
	}

	return listInstanceResponse, nil
}
