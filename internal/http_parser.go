package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type HTTPMethod int
type HTTPMessageType int

func (h HTTPMethod) String() string {
	m := map[HTTPMethod]string{
		GET:     "GET",
		POST:    "POST",
		PUT:     "PUT",
		PATCH:   "PATCH",
		DELETE:  "DELETE",
		CONNECT: "CONNECT",
		OPTIONS: "OPTIONS",
	}

	return m[h]
}

func HTTPMethodFromString(s string) (HTTPMethod, error) {
	m := map[string]HTTPMethod{
		"GET":     GET,
		"POST":    POST,
		"PUT":     PUT,
		"PATCH":   PATCH,
		"DELETE":  DELETE,
		"CONNECT": CONNECT,
		"OPTIONS": OPTIONS,
	}

	method, ok := m[s]
	if !ok {
		return HTTPMethod(-1), fmt.Errorf("error: method: parse: invalid method string '%s'", s)
	}

	return method, nil
}

const (
	GET HTTPMethod = iota
	POST
	PUT
	PATCH
	DELETE
	CONNECT
	OPTIONS
)

const (
	HTTPMessageRequest HTTPMessageType = iota
	HTTPMessageResponse
)

type HTTPRequest struct {
	Method  HTTPMethod
	Path    string
	Version string
}

type HTTPStatus struct {
	Status       int
	ReasonPhrase string
	Version      string
}

type HTTPMessage struct {
	Request HTTPRequest
	Status  HTTPStatus
	Headers map[string]string
	Body    []byte
}

func (m HTTPMessage) Serialize() []byte {
	s := fmt.Sprintf("%s %d %s\r\n", m.Status.Version, m.Status.Status, m.Status.ReasonPhrase)

	for key, value := range m.Headers {
		s += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	s += "\r\n"

	if len(m.Body) > 0 {
		s += string(m.Body)
	}

	return []byte(s)
}

func parseRequestLine(s string) (HTTPRequest, error) {
	r := bufio.NewReader(bytes.NewBuffer([]byte(s)))

	method, err := r.ReadString(' ')
	if err != nil {
		return HTTPRequest{}, err
	}

	m, err := HTTPMethodFromString(strings.TrimSpace(method))
	if err != nil {
		return HTTPRequest{}, err
	}

	requestPath, err := r.ReadString(' ')
	if err != nil {
		return HTTPRequest{}, err
	}

	return HTTPRequest{
		Method:  m,
		Path:    strings.TrimSpace(requestPath),
		Version: "HTTP/1.1",
	}, nil
}

func parseHeaderLine(s string) (string, string, error) {
	b := []byte(s)
	r := bufio.NewReader(bytes.NewBuffer(b))

	key, err := r.ReadString(':')
	if err != nil {
		return "", "", err
	}

	valueBytes := make([]byte, len(b)-len([]byte(key)))
	if _, err := r.Read(valueBytes); err != nil {
		return "", "", err
	}

	key = strings.ToLower(strings.TrimSpace(key[:len(key)-1]))
	value := strings.ToLower(strings.TrimSpace(string(valueBytes)))

	return key, value, nil
}

func ParseHTTPMessage(r io.Reader) (HTTPMessage, error) {
	reader := bufio.NewReader(r)

	startLine, err := reader.ReadString('\n')
	if err != nil {
		return HTTPMessage{}, err
	}

	startLine = strings.TrimSpace(startLine)
	requestLine, err := parseRequestLine(startLine)
	if err != nil {
		return HTTPMessage{}, err
	}

	m := HTTPMessage{Headers: map[string]string{}, Request: requestLine}
	contentLength := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return HTTPMessage{}, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		k, v, err := parseHeaderLine(line)
		if err != nil {
			return HTTPMessage{}, err
		}

		m.Headers[k] = v

		if k == "content-length" {
			contentLength, err = strconv.Atoi(v)
			if err != nil {
				return HTTPMessage{}, err
			}
		}
	}

	if contentLength > 0 {
		m.Body = make([]byte, contentLength)
		if _, err := io.ReadFull(reader, m.Body); err != nil {
			return HTTPMessage{}, err
		}
	}

	return m, nil
}
