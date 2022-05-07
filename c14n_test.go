package c14n_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/chenghonour/XMLC14n"
	"golang.org/x/net/html/charset"
)

func ExampleCanonicalize() {
	input := `<foo z="2" a="1"><bar /></foo>`
	decoder := xml.NewDecoder(strings.NewReader(input))
	out, err := c14n.Canonicalize(decoder)
	fmt.Println(string(out), err)
	// Output:
	// <foo a="1" z="2"><bar></bar></foo> <nil>
}

func TestCanonicalize(t *testing.T) {
	entries, err := ioutil.ReadDir("tests")
	assert.NoError(t, err)

	for _, file := range entries {
		t.Run(file.Name(), func(t *testing.T) {
			in, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/in.xml", file.Name()))
			assert.NoError(t, err)

			out, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/out.xml", file.Name()))
			assert.NoError(t, err)

			decoder := xml.NewDecoder(bytes.NewReader(in))
			decoder.CharsetReader = charset.NewReaderLabel

			actual, err := c14n.Canonicalize(decoder)
			assert.NoError(t, err)
			assert.Equal(t, out, actual)
		})
	}
}

func TestCanonicalize_NoStartElement(t *testing.T) {
	decoder := xml.NewDecoder(strings.NewReader("<!-- foo -->"))
	_, err := c14n.Canonicalize(decoder)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestCanonicalize_RawTokenError(t *testing.T) {
	_, err := c14n.Canonicalize(&errRawTokener{})
	assert.Equal(t, errDummy, err)
}

var errDummy = errors.New("dummy error")

type errRawTokener struct{}

func (e *errRawTokener) RawToken() (xml.Token, error) {
	return nil, errDummy
}
