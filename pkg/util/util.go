package util

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

func GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	path, _ = filepath.Abs(path)
	return path
}

func ParseTplStringWithDelims(tplt string, data any, left, right string) (string, error) {
	var res string
	tpl, err := template.New("").Delims(left, right).Parse(tplt)
	if err != nil {
		return res, errors.Wrap(err, "parse template")
	}
	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		return res, errors.Wrap(err, "tpl execute")
	}
	return buf.String(), nil
}

func ParseTplString(tplt string, data any) (string, error) {
	return ParseTplStringWithDelims(tplt, data, "{{", "}}")
}

func ParseTplBytesWithDelims(tplt []byte, data any, left, right string) ([]byte, error) {
	tpl, err := template.New("").Delims(left, right).Parse(string(tplt))
	// Delims
	if err != nil {
		return nil, errors.Wrap(err, "parse template")
	}
	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		return nil, errors.Wrap(err, "tpl execute")
	}
	return buf.Bytes(), nil
}

func ParseTplBytes(tplt []byte, data any) ([]byte, error) {
	return ParseTplBytesWithDelims(tplt, data, "{{", "}}")
}
