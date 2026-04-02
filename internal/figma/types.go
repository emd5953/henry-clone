package figma

// FileResponse is the response from GET /v1/files/:key
type FileResponse struct {
	Name         string   `json:"name"`
	LastModified string   `json:"lastModified"`
	Version      string   `json:"version"`
	Document     Document `json:"document"`
}

type Document struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Children []Node `json:"children"`
}

type Node struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Children []Node `json:"children,omitempty"`
}

// FileNodesResponse is the response from GET /v1/files/:key/nodes
type FileNodesResponse struct {
	Nodes map[string]NodeData `json:"nodes"`
}

type NodeData struct {
	Document Node `json:"document"`
}

// ExportResponse is the response from GET /v1/images/:key
type ExportResponse struct {
	Images map[string]string `json:"images"` // nodeID -> image URL
}
