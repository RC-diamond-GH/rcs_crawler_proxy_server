package proxy

// Uncompleted
import (
	"encoding/base64"
	"fmt"
	"strings"
)

type SSNode struct {
	Method   string
	Password string
	Host     string
	Port     string
}

func ParseSSURL(ssURL string) (*SSNode, error) {
	trimmed := strings.TrimPrefix(ssURL, "ss://")
	if idx := strings.Index(trimmed, "#"); idx != -1 {
		trimmed = trimmed[:idx]
	}

	decoded, err := base64.RawStdEncoding.DecodeString(trimmed)
	if err != nil {
		decoded, err = base64.StdEncoding.DecodeString(trimmed)
		if err != nil {
			return nil, fmt.Errorf("base64 decode failed: %v", err)
		}
	}

	parts := strings.SplitN(string(decoded), "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ss decoded format")
	}

	authPart := parts[0]
	addrPart := parts[1]

	mp := strings.SplitN(authPart, ":", 2)
	if len(mp) != 2 {
		return nil, fmt.Errorf("invalid method:password format")
	}
	method := mp[0]
	password := mp[1]

	hp := strings.SplitN(addrPart, ":", 2)
	if len(hp) != 2 {
		return nil, fmt.Errorf("invalid host:port format")
	}
	host := hp[0]
	port := hp[1]

	return &SSNode{
		Method:   method,
		Password: password,
		Host:     host,
		Port:     port,
	}, nil
}
