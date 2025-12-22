package ipfs

import (
	"be/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Pinata struct {
	config *config.Config
}

type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

func NewPinata(
	config *config.Config,
) *Pinata {

	return &Pinata{
		config: config,
	}
}

func (i *Pinata) Upload(fileName string, content []byte) (string, error) {
	// 1. create buffer and writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 2. create file
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = fw.Write(content); err != nil {
		return "", fmt.Errorf("failed to write content: %w", err)
	}

	// 3. Metadata
	metadata := map[string]interface{}{
		"name": fileName,
		"keyvalues": map[string]string{
			"type":    "gp-schema",
			"version": "1.0",
		},
	}
	metadataBytes, err := json.Marshal(metadata)
	if err == nil {
		_ = writer.WriteField("pinataMetadata", string(metadataBytes))
	}

	// 4. Options
	options := map[string]bool{"pin": true}
	optionsBytes, err := json.Marshal(options)
	if err == nil {
		_ = writer.WriteField("pinataOptions", string(optionsBytes))
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 5. Create Request
	url := i.config.IPFS.Endpoint + "/pinFileToIPFS"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	// Config Header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+i.config.IPFS.JWTKey)

	//
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result PinataResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.IpfsHash, nil
}

func (i *Pinata) Remove(cid string) error {
	if cid == "" {
		return fmt.Errorf("cid is required")
	}
	url := i.config.IPFS.Endpoint + "/upin/" + cid
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+i.config.IPFS.JWTKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pinata delete error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
