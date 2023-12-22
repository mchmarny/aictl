package url

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/k3a/html2text"
	"github.com/pkg/errors"
)

const (
	maxIdleConns     = 10
	timeoutInSeconds = 60
	clientAgent      = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"
)

var (
	reqTransport = &http.Transport{
		MaxIdleConns:          maxIdleConns,
		IdleConnTimeout:       timeoutInSeconds * time.Second,
		DisableCompression:    true,
		DisableKeepAlives:     false,
		ResponseHeaderTimeout: time.Duration(timeoutInSeconds) * time.Second,
	}
)

func getResp(url string) (resp *http.Response, err error) {
	c := http.Client{
		Timeout:   time.Duration(timeoutInSeconds) * time.Second,
		Transport: reqTransport,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating HTTP Get request")
	}

	req.Header.Set("User-Agent", clientAgent)

	return c.Do(req)
}

func GetContent(desc, url string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		return "", errors.Errorf("invalid url %s", url)
	}

	var content strings.Builder
	content.WriteString(desc)
	content.WriteString("\n")

	// get html content
	resp, err := getResp(url)
	if err != nil {
		return "", errors.Errorf("error requesting %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.Errorf("url not found %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("invalid response %s: %d - %s", url, resp.StatusCode, resp.Status)
	}

	// read html content
	html := new(strings.Builder)
	_, err = io.Copy(html, resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "error reading downloaded content from %s", url)
	}

	// convert html to text
	content.WriteString(html2text.HTML2Text(html.String()))

	return content.String(), nil
}
