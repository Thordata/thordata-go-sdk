package thordata

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	xproxy "golang.org/x/net/proxy"
)

type ProxyProduct string

const (
	ProxyResidential ProxyProduct = "residential"
	ProxyDatacenter  ProxyProduct = "datacenter"
	ProxyMobile      ProxyProduct = "mobile"
	ProxyISP         ProxyProduct = "isp"
)

func (p ProxyProduct) defaultPort() int {
	switch p {
	case ProxyResidential:
		return 9999
	case ProxyMobile:
		return 5555
	case ProxyDatacenter:
		return 7777
	case ProxyISP:
		return 6666
	default:
		return 9999
	}
}

func (p ProxyProduct) defaultHost() string {
	// NOTE: actual working host often comes from Dashboard endpoint generator:
	// e.g. vpn9wq0d.pr.thordata.net
	// We keep historical defaults as fallback only.
	switch p {
	case ProxyResidential:
		return "t.pr.thordata.net"
	case ProxyDatacenter:
		return "dc.pr.thordata.net"
	case ProxyMobile:
		return "m.pr.thordata.net"
	case ProxyISP:
		return "isp.pr.thordata.net"
	default:
		return "pr.thordata.net"
	}
}

type ProxyConfig struct {
	Product ProxyProduct

	// Gateway credentials (Residential/Datacenter/Mobile):
	// Username is the "account username" part (without td-customer- prefix)
	Username string
	Password string

	// Endpoint overrides (prefer Dashboard shard host)
	Protocol string // http or https (IMPORTANT: many accounts require https)
	Host     string
	Port     int

	// Whitelist mode (no auth)
	NoAuth bool

	// Targeting / session options (gateway mode)
	Continent       string
	Country         string
	State           string
	City            string
	ASN             string
	SessionID       string
	SessionDuration int // minutes (1-90)
}

func (p *ProxyConfig) CountryCode(code string) *ProxyConfig {
	p.Country = strings.ToLower(strings.TrimSpace(code))
	return p
}

func (p *ProxyConfig) CityName(name string) *ProxyConfig {
	p.City = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(name, " ", "_")))
	return p
}

func (p *ProxyConfig) Session(id string) *ProxyConfig {
	p.SessionID = strings.TrimSpace(id)
	return p
}

func (p *ProxyConfig) Sticky(minutes int) *ProxyConfig {
	p.SessionDuration = minutes
	return p
}

func (p *ProxyConfig) BuildGatewayUsername() (string, error) {
	if strings.TrimSpace(p.Username) == "" {
		return "", errors.New("proxy username is required")
	}

	parts := []string{"td-customer-" + p.Username}

	// Order matters (keep consistent with spec)
	if p.Continent != "" {
		parts = append(parts, "continent-"+strings.ToLower(p.Continent))
	}
	if p.Country != "" {
		parts = append(parts, "country-"+strings.ToLower(p.Country))
	}
	if p.State != "" {
		parts = append(parts, "state-"+strings.ToLower(p.State))
	}
	if p.City != "" {
		parts = append(parts, "city-"+strings.ToLower(p.City))
	}
	if p.ASN != "" {
		asn := strings.ToUpper(p.ASN)
		if !strings.HasPrefix(asn, "AS") {
			asn = "AS" + asn
		}
		parts = append(parts, "asn-"+asn)
	}
	if p.SessionID != "" {
		parts = append(parts, "sessid-"+p.SessionID)
	}
	if p.SessionDuration != 0 {
		if p.SessionDuration < 1 || p.SessionDuration > 90 {
			return "", fmt.Errorf("session duration must be 1-90 minutes, got %d", p.SessionDuration)
		}
		if p.SessionID == "" {
			return "", errors.New("sticky session requires session id")
		}
		parts = append(parts, "sesstime-"+strconv.Itoa(p.SessionDuration))
	}

	return strings.Join(parts, "-"), nil
}

func (p *ProxyConfig) effectiveEndpoint() (scheme, host string, port int) {
	scheme = strings.ToLower(strings.TrimSpace(p.Protocol))
	if scheme == "" {
		scheme = "https" // safe default (matches real-world requirement you observed)
	}
	if scheme == "socks5" {
		scheme = "socks5h"
	}
	host = strings.TrimSpace(p.Host)
	if host == "" {
		host = p.Product.defaultHost()
	}
	port = p.Port
	if port == 0 {
		port = p.Product.defaultPort()
	}
	return
}

func (p *ProxyConfig) effectiveAuth(host string) (string, string, error) {
	if p.NoAuth {
		return "", "", nil
	}

	// ISP direct: username/password are direct
	if p.Product == ProxyISP && !strings.Contains(host, ".pr.thordata.net") {
		if p.Username == "" || p.Password == "" {
			return "", "", errors.New("isp proxy requires username/password")
		}
		return p.Username, p.Password, nil
	}

	user, err := p.BuildGatewayUsername()
	if err != nil {
		return "", "", err
	}
	if p.Password == "" && !p.NoAuth {
		return "", "", errors.New("proxy password is required (user/pass mode)")
	}
	return user, p.Password, nil
}

func (p *ProxyConfig) proxyURLAndAuth() (*url.URL, string, error) {
	scheme, host, port := p.effectiveEndpoint()

	u := &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
	}

	if p.NoAuth {
		return u, "", nil
	}

	// ISP static proxies do NOT use td-customer- username format, but your SDK currently
	// models ISP as "static direct host" with username/password provided by dashboard.
	// For simplicity here: if product is isp and host is not a gateway host, treat Username as direct login.
	if p.Product == ProxyISP && !strings.Contains(host, ".pr.thordata.net") {
		if p.Username == "" || p.Password == "" {
			return nil, "", errors.New("isp proxy requires username/password")
		}
		return u, p.Username + ":" + p.Password, nil
	}

	user, err := p.BuildGatewayUsername()
	if err != nil {
		return nil, "", err
	}
	if p.Password == "" && !p.NoAuth {
		// Some dashboards allow empty password, but generally password is required for user/pass mode.
		// Keep compatible but explicit.
		return nil, "", errors.New("proxy password is required (user/pass mode)")
	}

	return u, user + ":" + p.Password, nil
}

func basicAuthHeader(userPass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userPass))
}

func DefaultProxyFromEnv() (*ProxyConfig, error) {
	// 1) whitelist mode
	whitelist := strings.ToLower(strings.TrimSpace(os.Getenv("THORDATA_PROXY_WHITELIST")))
	if whitelist == "1" || whitelist == "true" || whitelist == "yes" || whitelist == "y" {
		p := &ProxyConfig{
			Product: ProxyResidential, // irrelevant in whitelist; keep default port 9999
			NoAuth:  true,
		}
		applyEndpointEnvOverrides(p, ProxyResidential)
		return p, nil
	}

	// 2) residential
	if u := os.Getenv("THORDATA_RESIDENTIAL_USERNAME"); u != "" {
		if pw := os.Getenv("THORDATA_RESIDENTIAL_PASSWORD"); pw != "" {
			p := &ProxyConfig{Product: ProxyResidential, Username: u, Password: pw}
			applyEndpointEnvOverrides(p, ProxyResidential)
			return p, nil
		}
	}

	// 3) datacenter
	if u := os.Getenv("THORDATA_DATACENTER_USERNAME"); u != "" {
		if pw := os.Getenv("THORDATA_DATACENTER_PASSWORD"); pw != "" {
			p := &ProxyConfig{Product: ProxyDatacenter, Username: u, Password: pw}
			applyEndpointEnvOverrides(p, ProxyDatacenter)
			return p, nil
		}
	}

	// 4) mobile
	if u := os.Getenv("THORDATA_MOBILE_USERNAME"); u != "" {
		if pw := os.Getenv("THORDATA_MOBILE_PASSWORD"); pw != "" {
			p := &ProxyConfig{Product: ProxyMobile, Username: u, Password: pw}
			applyEndpointEnvOverrides(p, ProxyMobile)
			return p, nil
		}
	}

	return nil, nil
}

func applyEndpointEnvOverrides(p *ProxyConfig, product ProxyProduct) {
	prefix := strings.ToUpper(string(product)) // RESIDENTIAL/DATACENTER/MOBILE/ISP

	// Per-product overrides first
	host := os.Getenv("THORDATA_" + prefix + "_PROXY_HOST")
	portRaw := os.Getenv("THORDATA_" + prefix + "_PROXY_PORT")
	proto := os.Getenv("THORDATA_" + prefix + "_PROXY_PROTOCOL")

	// Global overrides
	if host == "" {
		host = os.Getenv("THORDATA_PROXY_HOST")
	}
	if portRaw == "" {
		portRaw = os.Getenv("THORDATA_PROXY_PORT")
	}
	if proto == "" {
		proto = os.Getenv("THORDATA_PROXY_PROTOCOL")
	}

	if strings.TrimSpace(host) != "" {
		p.Host = strings.TrimSpace(host)
	}
	if strings.TrimSpace(proto) != "" {
		p.Protocol = strings.TrimSpace(proto)
	}
	if strings.TrimSpace(portRaw) != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(portRaw)); err == nil {
			p.Port = n
		}
	}
}

// ProxyGet sends a GET request through Thordata Proxy Network.
// If proxy is nil, it tries DefaultProxyFromEnv().
func (c *Client) ProxyGet(ctx context.Context, targetURL string, proxy *ProxyConfig) (*http.Response, error) {
	return c.ProxyRequest(ctx, http.MethodGet, targetURL, proxy, nil, nil)
}

func buildProxyTransport(proxyURL *url.URL, authUser string, authPass string) (*http.Transport, error) {
	scheme := strings.ToLower(strings.TrimSpace(proxyURL.Scheme))
	if scheme == "socks5" {
		scheme = "socks5h"
	}

	baseDialer := &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// SOCKS5/SOCKS5H
	if strings.HasPrefix(scheme, "socks5") {
		var auth *xproxy.Auth
		if authUser != "" {
			auth = &xproxy.Auth{User: authUser, Password: authPass}
		}

		d, err := xproxy.SOCKS5("tcp", proxyURL.Host, auth, baseDialer)
		if err != nil {
			return nil, err
		}
		ctxDialer, _ := d.(xproxy.ContextDialer)

		tr := &http.Transport{
			TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS12},
			ForceAttemptHTTP2: true,
		}
		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if ctxDialer != nil {
				return ctxDialer.DialContext(ctx, network, addr)
			}
			return d.Dial(network, addr)
		}
		return tr, nil
	}

	// HTTP/HTTPS proxy (CONNECT)
	tr := &http.Transport{
		Proxy:             http.ProxyURL(proxyURL),
		TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS12},
		DialContext:       baseDialer.DialContext,
		ForceAttemptHTTP2: true,
	}

	if authUser != "" {
		pa := basicAuthHeader(authUser + ":" + authPass)
		tr.ProxyConnectHeader = make(http.Header)
		tr.ProxyConnectHeader.Set("Proxy-Authorization", pa)
	}
	return tr, nil
}

func (c *Client) ProxyRequest(
	ctx context.Context,
	method string,
	targetURL string,
	proxy *ProxyConfig,
	body io.Reader,
	headers map[string]string,
) (*http.Response, error) {
	if strings.TrimSpace(targetURL) == "" {
		return nil, errors.New("targetURL is required")
	}

	if proxy == nil {
		var err error
		proxy, err = DefaultProxyFromEnv()
		if err != nil {
			return nil, err
		}
	}

	if proxy == nil {
		return nil, errors.New(
			"proxy credentials are missing" +
				"Set THORDATA_RESIDENTIAL_USERNAME/THORDATA_RESIDENTIAL_PASSWORD (or DATACENTER/MOBILE) " +
				"and THORDATA_PROXY_HOST/PORT/PROTOCOL from Dashboard endpoint generator",
		)
	}

	// We keep proxyURLAndAuth for endpoint building, but do NOT parse userPass for socks.
	proxyURL, _, err := proxy.proxyURLAndAuth()
	if err != nil {
		return nil, err
	}

	// Resolve auth safely
	host := proxyURL.Hostname()
	authUser, authPass, err := proxy.effectiveAuth(host)
	if err != nil {
		return nil, err
	}

	// Build transport (http/https/socks5h)
	tr, err := buildProxyTransport(proxyURL, authUser, authPass)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.cfg.UserAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Only set Proxy-Authorization header for HTTP/HTTPS proxies (harmless otherwise, but keep clean)
	scheme := strings.ToLower(strings.TrimSpace(proxyURL.Scheme))
	if scheme == "http" || scheme == "https" {
		if authUser != "" {
			req.Header.Set("Proxy-Authorization", basicAuthHeader(authUser+":"+authPass))
		}
	}

	client := &http.Client{
		Timeout:   c.cfg.Timeout,
		Transport: tr,
	}
	return client.Do(req)
}
