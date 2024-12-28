package proxy

import (
	"io"
	"net/http"
	"strings"

	"rcs_crawler_proxy_server/util"

	"github.com/sirupsen/logrus"
)

var logger logrus.Logger
var proxyClients []*http.Client
var localClient *http.Client

func httpProxyHandler(writer http.ResponseWriter, reader *http.Request) {

	logger.Info("Received request ", "method ", reader.Method, "url ", reader.URL.String())

	requestInfo, err := util.BuildHTTPIdentifier(reader, "http")
	if err != nil {
		logger.Error("Error while reading HTTP request info: " + err.Error())
	}

	_, noCache := reader.Header["No-Proxy-Cache"]
	reader.Header.Del("No-Proxy-Cache")
	cacheKey, err := requestInfo.CacheKey()
	if err != nil {
		logger.Error("Error while calculating cacheKey: " + err.Error())
	}

	if !noCache {
		logger.Info("Handle request with cache.")
		hasCache, err := util.Cache.Exists(cacheKey)
		hasCache = cacheKey != "" && hasCache

		if err != nil {
			logger.Errorf("Error while checking cache: %v", err)
		}

		if hasCache {
			logger.Info("Trigger Cache")
			res, err := util.Cache.GetHTTPResponseCache(cacheKey)
			if err != nil {
				logger.Errorf("Error while getting from cache: %v", err)
			}

			if err = res.SendbackResponse(writer); err != nil {
				logger.Errorf("Error while sending response to client: %v", err)
			}
			return
		}
	}

	logger.Info("Processing remote request...")
	reqHeader := reader.Header

	req, err := http.NewRequest(requestInfo.Method, requestInfo.URL, reader.Body)
	if err != nil {
		logger.Errorf("Error while BUILDING HTTP request: %v", err)
		http.Error(writer, "Proxy Server Error", http.StatusBadRequest)
		return
	}

	for k, v := range reqHeader {
		req.Header.Set(k, strings.Join(v, ","))
	}

	logger.Debug("Request headers", "headers", reqHeader)

	// Request from remote
	resp, err := doRequest(req)
	//resp, err := proxyClient.Do(req)
	if err != nil {
		logger.Errorf("Error while DOING HTTP request %v", err.Error())
		http.Error(writer, "Proxy Server Error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Construct HTTP Response Cache
	statusCode := resp.StatusCode
	logger.Info("Received response from proxy ", "status_code ", statusCode)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error while construct HTTP Response Cache: %v", err)
		http.Error(writer, "Proxy Server Error", http.StatusBadGateway)
		return
	}

	response := util.HTTPResponseCache{
		StatusCode: statusCode,
		Header:     resp.Header,
		Data:       data,
	}
	response.SendbackResponse(writer)

	logger.Info("Successfully forwarded response", "status_code", statusCode)
	logger.Info("Save response in cache")

	if !noCache {
		if err = util.Cache.SetHTTPResponseCache(cacheKey, response); err != nil {
			logger.Errorf("Error while saving cache: %v", err)
		}
	}

}
func doRequest(req *http.Request) (*http.Response, error) {
	logger.Infof("Direct request: %s", req.URL)
	resp, err := localClient.Do(req)
	if err == nil {
		return resp, err

	}
	// Use proxy
	logger.Infof("Request %s with proxy", req.URL)
	for _, p := range proxyClients {
		resp, err := p.Do(req)
		if err == nil {
			return resp, err
		}
	}
	return resp, err
}
func HttpListener() {
	logger = *util.GetLogger()

	logger.Info("Initing proxy servers...")
	for _, proxyURL := range HTTPProxyList {
		logger.Infof("Initing proxy server: %s", proxyURL)
		proxyClients = append(proxyClients, &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		})
	}
	localClient = &http.Client{}

	logger.Info("Start HTTP Listener on 8080")
	http.HandleFunc("/", httpProxyHandler)
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
