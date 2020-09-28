package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
)

type file struct {
	filename string
	fileData []byte
}

type files []file

// @return compress file
func (f *files) zip(outputName string) error {
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	for _, fl := range *f {
		fHandler, err := w.Create(fl.filename)
		if err != nil {
			return err
		}
		_, err = fHandler.Write(fl.fileData)
		if err != nil {
			return err
		}
	}
	err := w.Close()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputName, buf.Bytes(), 0644)
}
