package ipfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Multipart struct {
	url      string
	byteData *bytes.Buffer
	writer   *multipart.Writer
}

func NewMultipart(url string) *Multipart {
	var byteData bytes.Buffer
	writer := multipart.NewWriter(&byteData)
	return &Multipart{url, &byteData, writer}
}

func (m *Multipart) AddFile(path, name string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	fw, err := m.writer.CreateFormFile("file", name)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return err
	}
	return nil
}

func (m *Multipart) AddSubDir(dirPath, basePath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			m.AddFile(dirPath+"/"+f.Name(), basePath+"/"+f.Name())
		} else {
			m.AddSubDir(dirPath+"/"+f.Name(), basePath+"/"+f.Name())
		}
	}
	return nil
}

func (m *Multipart) Send() ([]byte, error) {
	m.writer.Close()
	req, err := http.NewRequest("POST", m.url, m.byteData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", m.writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
