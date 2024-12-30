package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type HTTPRequestIdentifier struct {
	Method string
	URL    string
	data   []byte
}

func (obj *HTTPRequestIdentifier) CacheKey() (string, error) {
	dataToHash := []byte(obj.Method + obj.URL)
	dataToHash = append(dataToHash, obj.data...)

	hash := sha256.New()
	_, err := hash.Write(dataToHash)
	if err != nil {
		return "", fmt.Errorf("error computing hash: %s", err.Error())
	}

	hashBytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashBytes), nil
}

func BuildHTTPIdentifier(reader *http.Request, protocol string) (HTTPRequestIdentifier, error) {
	method := reader.Method

	reqURL := reader.URL
	url := protocol + "://"
	if reader.Host == "" {
		url += reqURL.Host
	} else {
		url += reader.Host
	}
	query := reqURL.Query().Encode()
	if query != "" {
		url += "?" + query
	}
	//url := protocol + "://" + reqURL.Host + reqURL.Path + "?" + reqURL.Query().Encode()
	data, err := io.ReadAll(reader.Body)
	return HTTPRequestIdentifier{
		Method: method,
		URL:    url,
		data:   data,
	}, err
}

type HTTPResponseCache struct {
	StatusCode int
	Header     map[string][]string
	Data       []byte
}

func (c *HTTPResponseCache) SendbackResponse(writer http.ResponseWriter) error {
	writer.WriteHeader(c.StatusCode)
	for k, v := range c.Header {
		for _, value := range v {
			writer.Header().Add(k, value)
		}
	}
	logger.Debug("Response headers", "headers", c.Header)
	_, err := io.Copy(writer, bytes.NewReader(c.Data))

	return err
}
