package basicutils

import (
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}

		}
	}

	return "", errors.New("未能获取到本机 IP 地址 !")
}

func CheckURL(url string) (bool, error) {
	client := &http.Client{
		Timeout: 1 * time.Second, // 设置超时时间
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, errors.Errorf("URL responded with status %d\n", resp.StatusCode)
}
