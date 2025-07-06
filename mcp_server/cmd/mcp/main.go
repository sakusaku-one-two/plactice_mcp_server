package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type JSONRPCResponse struct {
	JSONPRC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

type MCPServer struct {
	name      string
	version   string
	tools     map[string]Tool
	resources map[string]Resource
}

func NewMCPServer(name, version string) *MCPServer {
	return &MCPServer{
		name:      name,
		version:   version,
		tools:     make(map[string]Tool),
		resources: make(map[string]Resource),
	}
}

// handleInitialize hadles the intiialize method
// hadleInitialize: initialize メソッドを処理する関数
// handles: 処理する/扱う
func (s *MCPServer) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	// servier capabilities: サーバーの機能
	// capabilities: 機能、能力
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true, // listChanged: リスト変更通知
			},
			"resource": map[string]interface{}{
				"subscribe":   true, // subscribe:購買する
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    s.name,    // name:名前
			"version": s.version, // verison: バージョン

		},
	}

	return &JSONRPCResponse{
		JSONPRC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsList hadles the tools/list method
// hadleToolsList: tools/listメソッドを処理する関数
func (s *MCPServer) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	return &JSONRPCResponse{
		JSONPRC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

// hadletoolsCall hadles the tools/call method
// hadleToolsCall: tools/call メソッドを処理する関数
func (s *MCPServer) handleToolsCall(req *JSONRPCRequest) *JSONRPCResponse {
	params, ok := req.Params.(map[string]interface{}) // type assertion:型アサーション
	if !ok {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid parameters", //parameters:パラメーター
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Tool name is required", // required : 必須の
			},
		}
	}
	// Security:ツール名の検証
	// security: セキュリティ安全
	if _, exists := s.tools[toolName]; !exists {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -326001,
				Message: "Tool not found", // found: 見つかった
			},
		}
	}

	// Execute tools:ツールを実行
	// execute: 実行する遂行する

	reuslt := s.executeTool(toolName, params["arguments"])

	return &JSONRPCResponse{
		JSONPRC: "2.0",
		ID:      req.ID,
		Result:  reuslt,
	}
}

func (s *MCPServer) handleResourcesList(req *JSONRPCRequest) *JSONRPCResponse {
	resources := make([]Resource, 0, len(s.resources))
	for _, resource := range s.resources {
		resources = append(resources, resource)
	}

	return &JSONRPCResponse{
		JSONPRC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"resources": resources,
		},
	}
}

func (s *MCPServer) hadleResourcesRead(req *JSONRPCRequest) *JSONRPCResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid parameters",
			},
		}
	}
	uri, ok := params["uri"].(string)
	if !ok {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid parameters",
			},
		}
	}

	// security uri validation

	if !strings.HasPrefix(uri, "file://") && !strings.HasPrefix(uri, "https://") {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid URI scheme",
			},
		}
	}

	// Read resource
	content := s.readResource(uri)

	return &JSONRPCResponse{
		JSONPRC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"contents": []map[string]interface{}{
				{"uri": uri,
					"mimeType": "text/plain", // mimeType
					"text":     content,
				},
			},
		},
	}
}

func (s *MCPServer) executeTool(toolName string, arguments interface{}) map[string]interface{} {
	// Example tool execution ツール実行の例
	// execution 実行遂行
	switch toolName {
	case "echo":
		args, ok := arguments.(map[string]interface{})
		if !ok {
			return map[string]interface{}{
				"error": "Invalid arguments",
			}
		}
		message, ok := args["message"].(string)
		if !ok {
			return map[string]interface{}{
				"error": "Message is required",
			}
		}
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Echo: %s", message),
				},
			},
		}
	default:
		return map[string]interface{}{
			"error": "Unknown tool",
		}

	}
}

func (s *MCPServer) readResource(uri string) string {
	// Example resource reading: リソース読み取りの例
	switch {
	case strings.HasPrefix(uri, "file://"):
		return fmt.Sprintf("Content of %s", uri)
	case strings.HasPrefix(uri, "https://"):
		return fmt.Sprintf("Web content of %s", uri)
	default:
		return "Unknown resoruce type"
	}
}

// Run starts the MCP server
// Run: MCPサーバーを開始する関数
// starts: 開始する、始める
func (s *MCPServer) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		//Skip empty lines:空行をスキップ
		if strings.TrimSpace(line) == "" {
			continue
		}
		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("JSON parsing error: %v", err)
			continue
		}

		//Process request
		resp := s.HandleRequest(&req)
		// Send response: レスポンスを送信
		// send: 送信する
		if respData, err := json.Marshal(resp); err == nil {
			fmt.Println(string(respData))
		} else {
			log.Printf("JSON marshaling error: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error:  %v", err)
	}
}

func (s *MCPServer) RegisterResource(resource Resource) {
	s.resources[resource.URI] = resource
}

func (s *MCPServer) hadlerRequest(req *JSONRPCRequest) *JSONRPCResponse {
	// Input validation:セキュリティのための入力検証
	//validateion: 検証、妥当性確認
	if req.JSONRPC != "2.0" {
		return &JSONRPCResponse{
			JSONPRC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32600,                     //INvalid Request(無効なリクエスト)
				Message: "Invalid JSON-RPC version", //version: バージョン
			},
		}
	}

	// method dispatch: メソッドの振り分け
	// dispatch: 振り分ける、発送する。

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)

	}
}

func main() {
	fmt.Println("strat mcp server")
}
