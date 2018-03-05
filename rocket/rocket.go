package rocket

import (
	"net/http"
	"encoding/json"
	"github.com/aspcartman/exceptions"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"io"
	"crypto/md5"
	"encoding/hex"
	"time"
)

type Session struct {
	Phone, Email, Password, Token string
}

func (s *Session) Register() {
	fmt.Println(string(s.req("POST", "/devices/register", nil, mm{
		"phone": s.Phone,
	})))
}

func (s *Session) Login() {
	fmt.Println(string(s.req("GET", "/login", nil, nil)))
}

func (s *Session) req(method string, path string, query, body map[string]interface{}) []byte {
	req, err := http.NewRequest(method, s.makePath(path, query), s.makeBody(body))
	if err != nil {
		e.Throw("failed creating request", err)
	}

	req.Header = s.headers()

	switch {
	case len(s.Token) > 0:
		// ?
	case len(s.Email) > 0 && len(s.Password) > 0:
		req.SetBasicAuth(s.Email, s.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.Throw("failed doing request", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.Throw("failed reading a response body", err)
	}

	return resBody
}

func (s *Session) makePath(path string, query mm) string {
	q := url.Values{}
	for k, v := range query {
		q[k] = []string{fmt.Sprint(v)}
	}

	p := url.URL{
		Scheme:   "https",
		Host:     "rocketbank.ru",
		Path:     "/api/v5" + path,
		RawQuery: q.Encode(),
	}

	return p.String()
}

func (s *Session) makeBody(obj interface{}) io.ReadCloser {
	if obj == nil {
		return nil
	}
	data, err := json.Marshal(obj)
	if err != nil {
		e.Throw("failed marshalling body", err)
	}
	return bbody{bytes.NewReader(data)}
}

func (s *Session) headers() http.Header {
	now := fmt.Sprint(time.Now().Unix())
	return http.Header{
		"accept":          []string{"*/*"},
		"accept-language": []string{"en-us"},
		"content-length":  []string{"103"},
		"content-type":    []string{"application/json"},
		"user-agent":      []string{"rocketbank/156 CFNetwork/889.9 Darwin/17.2.0"},
		"x-app-version":   []string{"4.9.26 (156)"},
		"x-device-id":     []string{"10101010-0000-0000-0000-101010101010"},
		"x-device-idfa":   []string{"10101010-0000-0000-0000-101010101010"},
		"x-device-locale": []string{"ru_RU"},
		"x-device-os":     []string{"iOS 11.1.2"},
		"x-device-type":   []string{"iPhone7,2"},
		"x-screen-height": []string{"667.000000"},
		"x-screen-scale":  []string{"2.000000"},
		"x-screen-width":  []string{"375.000000"},
		"x-sig":           []string{sig(now)},
		"x-time":          []string{now},
	}
}

func sig(t string) string {
	hash := md5.Sum([]byte("0Jk211uvxyyYAFcSSsBK3+etfkDPKMz6asDqrzr+f7c=_" + t + "_dossantos"))
	return hex.EncodeToString(hash[:])
}

type bbody struct {
	*bytes.Reader
}

func (b bbody) Close() error {
	return nil
}

type mm map[string]interface{}
