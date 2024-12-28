package proxy

import (
	"net/url"
	"rcs_crawler_proxy_server/util"
)

var HTTPProxyList []*url.URL

func ParseHTTPProxyList() {
	for _, addr := range util.Config.Proxies {
		newURL, err := url.Parse(addr)
		if err != nil {
			util.GetLogger().Errorf("Error while parsing %s as proxy URL: %v", addr, err)
			continue
		}
		HTTPProxyList = append(HTTPProxyList, newURL)
	}
}

// 测试代理服务器列表的连通性
func PingProxyServer() {
	logger.Info("Pinging the proxy server...")
	//todo
}
