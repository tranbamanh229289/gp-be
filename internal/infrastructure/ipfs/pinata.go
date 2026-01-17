package ipfs

import (
	"be/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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
	if fileName == "" || len(content) == 0 {
		return "", fmt.Errorf("fileName and content cannot be empty")
	}

	// 1. Create buffer and writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 2. Create file form field
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := fw.Write(content); err != nil {
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
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := writer.WriteField("pinataMetadata", string(metadataBytes)); err != nil {
		return "", fmt.Errorf("failed to write metadata field: %w", err)
	}

	// 4. Options
	options := map[string]bool{"pin": true}
	optionsBytes, err := json.Marshal(options)
	if err != nil {
		return "", fmt.Errorf("failed to marshal options: %w", err)
	}
	if err := writer.WriteField("pinataOptions", string(optionsBytes)); err != nil {
		return "", fmt.Errorf("failed to write options field: %w", err)
	}

	// 5. Close writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 6. Create request
	url, _ := url.JoinPath(i.config.IPFS.Endpoint, "pinFileToIPFS")

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+i.config.IPFS.JWTKey)

	// 7. HTTP client
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Pinata response status:", resp.StatusCode)
	fmt.Println("Pinata response body:", string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Pinata upload failed, status %d: %s", resp.StatusCode, string(respBody))
	}

	// 8. Parse response
	var result PinataResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal Pinata response: %w", err)
	}

	fmt.Println("[DEBUG] Upload success, IpfsHash:", result.IpfsHash)
	return result.IpfsHash, nil
}

func (i *Pinata) Remove(cid string) error {
	if cid == "" {
		return fmt.Errorf("cid is required")
	}
	url, _ := url.JoinPath(i.config.IPFS.Endpoint, "unpin", cid)

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

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println("Pinata response status:", resp.StatusCode)
	fmt.Println("Pinata response body:", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pinata delete error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
