package internal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseHTTPMessage_Request_GET_NoHeaders(t *testing.T) {

	m := "GET / HTTP/1.1\r\n"
	m += "\r\n"

	r := bytes.NewBuffer([]byte(m))
	message, err := ParseHTTPMessage(r)

	assert.NoError(t, err)

	assert.Equal(t, 0, message.Status.Status)
	assert.Len(t, message.Body, 0)
	assert.Equal(t, GET, message.Request.Method)
	assert.Equal(t, "/", message.Request.Path)
	assert.Len(t, message.Headers, 0)

}

func Test_ParseHTTPMessage_Request_GET_WithOneHeader(t *testing.T) {

	m := "GET / HTTP/1.1\r\n"
	m += "Host: www.example.com\r\n"
	m += "\r\n"

	r := bytes.NewBuffer([]byte(m))
	message, err := ParseHTTPMessage(r)

	assert.NoError(t, err)

	assert.Equal(t, 0, message.Status.Status)
	assert.Len(t, message.Body, 0)
	assert.Equal(t, GET, message.Request.Method)
	assert.Equal(t, "/", message.Request.Path)
	assert.Equal(t, map[string]string{"host": "www.example.com"}, message.Headers)

}

func Test_ParseHTTPMessage_Request_GET_WithMultipleHeaders(t *testing.T) {

	m := "GET / HTTP/1.1\r\n"
	m += "Host: www.example.com\r\n"
	m += "Foo: bar\r\n"
	m += "\r\n"

	r := bytes.NewBuffer([]byte(m))
	message, err := ParseHTTPMessage(r)

	assert.NoError(t, err)

	assert.Equal(t, map[string]string{"host": "www.example.com", "foo": "bar"}, message.Headers)

}

func Test_ParseHTTPMessage_Request_POST_WithMultipleHeaders_WithBody(t *testing.T) {

	m := "POST / HTTP/1.1\r\n"
	m += "Host: www.example.com\r\n"
	m += "Foo: bar\r\n"
	m += "Content-Length: 19\r\n"
	m += "\r\n"
	m += "hello=world&bar=baz"

	r := bytes.NewBuffer([]byte(m))
	message, err := ParseHTTPMessage(r)

	assert.NoError(t, err)

	assert.Equal(t, POST, message.Request.Method)
	assert.Equal(t, "/", message.Request.Path)
	assert.Equal(t, map[string]string{"host": "www.example.com", "foo": "bar", "content-length": "19"}, message.Headers)
	assert.Equal(t, "hello=world&bar=baz", string(message.Body))

}
