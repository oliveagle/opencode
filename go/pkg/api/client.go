package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/anomalyco/opencode/pkg/types"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *Client) ReadFile(path string) (string, error) {
	resp, err := c.get("/api/files/read", url.Values{"path": {path}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return "", err
	}

	return data["content"], nil
}

func (c *Client) WriteFile(path, content string) error {
	body := map[string]string{"path": path, "content": content}
	resp, err := c.post("/api/files/write", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) ListDir(path string) ([]types.FileInfo, error) {
	if path == "" {
		path = "."
	}
	resp, err := c.get("/api/files/list", url.Values{"path": {path}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Files []types.FileInfo `json:"files"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Files, nil
}

func (c *Client) GetFileInfo(path string) (*types.FileInfo, error) {
	resp, err := c.get("/api/files/info", url.Values{"path": {path}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data types.FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteFile(path string) error {
	resp, err := c.get("/api/files/delete", url.Values{"path": {path}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) CreateDir(path string) error {
	resp, err := c.get("/api/files/create-dir", url.Values{"path": {path}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) ListProjects() ([]*types.ProjectConfig, error) {
	resp, err := c.get("/api/projects/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Projects []*types.ProjectConfig `json:"projects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Projects, nil
}

func (c *Client) GetProject(name string) (*types.ProjectConfig, error) {
	resp, err := c.get("/api/projects/get", url.Values{"name": {name}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data types.ProjectConfig
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) CreateProject(proj *types.ProjectConfig) error {
	resp, err := c.post("/api/projects/create", proj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) DeleteProject(name string) error {
	resp, err := c.get("/api/projects/delete", url.Values{"name": {name}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) Health() error {
	resp, err := c.get("/health", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) get(path string, query url.Values) (*http.Response, error) {
	u := fmt.Sprintf("%s%s", c.baseURL, path)
	if query != nil {
		u += "?" + query.Encode()
	}
	resp, err := c.client.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) post(path string, body interface{}) (*http.Response, error) {
	u := fmt.Sprintf("%s%s", c.baseURL, path)
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(u, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
