package podio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type File struct {
	Id                uint64  `json:"file_id"`
	Name              string `json:"name"`
	Link              string `json:"link"`
	Size              int    `json:"size"`
	Push              Push   `json:"push"`
	MimeType          string `json:"mimetype"`
	HostedBy          string `json:"hosted_by"`
	HumanizedHostedBy string `json:"hosted_by_humanized_name"`
}

// https://developers.podio.com/doc/files/get-files-4497983
func (client *Client) GetFiles() (files []File, err error) {
	err = client.Request("GET", "/file", nil, nil, &files)
	return
}

// https://developers.podio.com/doc/files/get-file-22451
func (client *Client) GetFile(fileId int) (file *File, err error) {
	err = client.Request("GET", fmt.Sprintf("/file/%d", fileId), nil, nil, &file)
	return
}

func (client *Client) GetFileContents(url string) ([]byte, error) {
	body, err := client.GetFileReader(url)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(body)
	defer body.Close()

	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (client *Client) GetFileReader(url string) (io.ReadCloser, error) {
	link := fmt.Sprintf("%s?oauth_token=%s", url, client.authToken.AccessToken)
	client.Emitter.FireBackground("podio.request", "GET", url, struct {
		Token string `json:"access_token"`
	}{
		client.authToken.blurToken(),
	})

	resp, err := http.Get(link)

	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// https://developers.podio.com/doc/files/upload-file-1004361
func (client *Client) CreateFile(name string, contents []byte) (file *File, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("source", name)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(contents)
	if err != nil {
		return nil, err
	}

	err = writer.WriteField("filename", name)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	err = client.Request("POST", "/file", headers, body, &file)
	return
}

// https://developers.podio.com/doc/files/replace-file-22450
func (client *Client) ReplaceFile(oldFileId, newFileId int) error {
	path := fmt.Sprintf("/file/%d/replace", newFileId)
	params := map[string]interface{}{
		"old_file_id": oldFileId,
	}

	return client.RequestWithParams("POST", path, nil, params, nil)
}

// https://developers.podio.com/doc/files/attach-file-22518
func (client *Client) AttachFile(fileId int, refType string, refId int) error {
	path := fmt.Sprintf("/file/%d/attach", fileId)
	params := map[string]interface{}{
		"ref_type": refType,
		"ref_id":   refId,
	}

	return client.RequestWithParams("POST", path, nil, params, nil)
}

// https://developers.podio.com/doc/files/delete-file-22453
func (client *Client) DeleteFile(fileId int) error {
	path := fmt.Sprintf("/file/%d", fileId)
	return client.Request("DELETE", path, nil, nil, nil)
}
