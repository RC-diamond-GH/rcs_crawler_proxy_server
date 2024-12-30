package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"rcs_crawler_proxy_server/util"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = util.GetLogger()
var httpsPort int = util.Config.ProxySetting.HttpsPort
var httpsFakeDest string = fmt.Sprintf("127.0.0.1:%d", httpsPort)
var proxyClients []*http.Client
var localClient *http.Client

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

func handleRequestsWithCache(w http.ResponseWriter, r *http.Request, protocol string) {
	logger.Infof("%s Listener Received request method: %s, URL: %s", protocol, r.Method, r.URL.String())
	logger.Infof("Host : %s ", r.Host)

	requestInfo, err := util.BuildHTTPIdentifier(r, protocol)
	if err != nil {
		logger.Error("Error while reading HTTP request info: " + err.Error())
	}

	_, noCache := r.Header["No-Proxy-Cache"]
	r.Header.Del("No-Proxy-Cache")
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

			if err = res.SendbackResponse(w); err != nil {
				logger.Errorf("Error while sending response to client: %v", err)
			}
			return
		}
	}

	logger.Info("Processing remote request...")
	reqHeader := r.Header

	req, err := http.NewRequest(requestInfo.Method, requestInfo.URL, r.Body)
	if err != nil {
		logger.Errorf("Error while BUILDING HTTP request: %v", err)
		http.Error(w, "Proxy Server Error", http.StatusBadRequest)
		return
	}

	for k, v := range reqHeader {
		req.Header.Set(k, strings.Join(v, ","))
	}

	logger.Debug("Request headers", "headers", reqHeader)

	// Request from remote
	resp, err := doRequest(req)
	if err != nil {
		logger.Errorf("Error while DOING HTTP request %v", err.Error())
		http.Error(w, "Proxy Server Error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Construct HTTP Response Cache
	statusCode := resp.StatusCode
	logger.Info("Received response from proxy ", "status_code ", statusCode)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error while construct HTTP Response Cache: %v", err)
		http.Error(w, "Proxy Server Error", http.StatusBadGateway)
		return
	}

	response := util.HTTPResponseCache{
		StatusCode: statusCode,
		Header:     resp.Header,
		Data:       data,
	}
	response.SendbackResponse(w)

	logger.Info("Successfully forwarded response", "status_code", statusCode)
	logger.Info("Save response in cache")

	if !noCache {
		if err = util.Cache.SetHTTPResponseCache(cacheKey, response); err != nil {
			logger.Errorf("Error while saving cache: %v", err)
		}
	}
}

func handleConnect(w http.ResponseWriter) {
	targetConn, err := net.Dial("tcp", httpsFakeDest)
	if err != nil {
		http.Error(w, "Failed to connect to target server", http.StatusBadGateway)
		logger.Errorf("Failed to connect to %s: %v", httpsFakeDest, err)
		return
	}
	defer targetConn.Close()
	logger.Info("Tunnel set up")

	// Tell the client that connection has established.
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	// Get client connection.
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		logger.Printf("Failed to hijack connection: %v", err)
		return
	}

	defer clientConn.Close()

	// pivot the data
	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}

func httpsProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		logger.Info("HTTPS Listener Received CONNECT request ", "method ", r.Method, " url ", r.URL.String())
		handleConnect(w)
		return
	}
	handleRequestsWithCache(w, r, "https")
}

func httpProxyHandler(w http.ResponseWriter, r *http.Request) {
	handleRequestsWithCache(w, r, "http")
}

func httpListener() {
	port := util.Config.ProxySetting.HttpPort
	logger.Infof("Start HTTP Listener on %d", port)
	http.HandleFunc("/", httpProxyHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Errorf("Failed to start HTTPS Listener: %v", err)
	}
}

func httpsListener() {
	logger.Info("Initial HTTPS Listener...")

	cert, err := tls.LoadX509KeyPair(util.Config.ProxySetting.TLSCert, util.Config.ProxySetting.TLSKey)
	if err != nil {
		logger.Errorf("Failed to load certificate: %v", err)
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpsPort))
	if err != nil {
		logger.Errorf("Failed to start TCP Listener: %v", err)
		return
	}
	logger.Infof("TCP Listener started on :%d", httpsPort)

	logger.Info("Creating TLS Listener...")
	tlsListener := tls.NewListener(listener, tlsConfig)
	logger.Infof("TLS Listener is ready on :%d", httpsPort)

	server := &http.Server{
		Handler: http.HandlerFunc(httpsProxyHandler),
	}

	logger.Infof("Start HTTPS Listener on %d", httpsPort)

	if err = server.Serve(tlsListener); err != nil {
		logger.Errorf("Failed to start HTTPS Listener: %v", err)
	}
}

func LaunchListeners() {
	for _, proxyURL := range HTTPProxyList {
		logger.Infof("Initing proxy server: %s", proxyURL)
		proxyClients = append(proxyClients, &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		})
	}
	localClient = &http.Client{}
	go httpListener()
	go httpsListener()
}
