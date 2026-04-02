// Package figma handles the integration with Figma's API.
// Henry's QC workflow likely involves pushing generated decks to a design tool
// where analysts and designers can do pixel-perfect edits — layout tweaks,
// photo placement, typography adjustments, etc. — before final delivery.
//
// The flow:
// 1. Deck is generated (HTML/structured sections)
// 2. Sections are pushed to Figma as native frames in a design file
// 3. QC team edits in Figma (visual layout, branding, photos)
// 4. Edits are synced back or the final PDF is exported from Figma
//
// This uses Figma's REST API for file/node operations.
// For the MCP-based write-to-canvas flow, the frontend would use
// the Figma MCP server directly.
package figma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.figma.com/v1"

// Client wraps the Figma REST API.
type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetFile retrieves a Figma file's metadata and node tree.
func (c *Client) GetFile(ctx context.Context, fileKey string) (*FileResponse, error) {
	url := fmt.Sprintf("%s/files/%s", baseURL, fileKey)
	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}

// GetFileNodes retrieves specific nodes from a Figma file.
func (c *Client) GetFileNodes(ctx context.Context, fileKey string, nodeIDs []string) (*FileNodesResponse, error) {
	url := fmt.Sprintf("%s/files/%s/nodes", baseURL, fileKey)
	// Add node IDs as query params
	if len(nodeIDs) > 0 {
		url += "?ids="
		for i, id := range nodeIDs {
			if i > 0 {
				url += ","
			}
			url += id
		}
	}

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FileNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}

// ExportNodes exports nodes as images (PNG, SVG, PDF).
func (c *Client) ExportNodes(ctx context.Context, fileKey string, nodeIDs []string, format string) (*ExportResponse, error) {
	url := fmt.Sprintf("%s/images/%s?ids=%s&format=%s", baseURL, fileKey, joinIDs(nodeIDs), format)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ExportResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}

// PostComment adds a comment to a Figma file (useful for QC notes).
func (c *Client) PostComment(ctx context.Context, fileKey string, message string) error {
	url := fmt.Sprintf("%s/files/%s/comments", baseURL, fileKey)
	body := map[string]string{"message": message}
	bodyBytes, _ := json.Marshal(body)

	resp, err := c.doRequest(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Figma-Token", c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("figma API error %d: %s", resp.StatusCode, string(respBody))
	}
	return resp, nil
}

func joinIDs(ids []string) string {
	result := ""
	for i, id := range ids {
		if i > 0 {
			result += ","
		}
		result += id
	}
	return result
}
