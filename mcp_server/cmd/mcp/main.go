package main

import (
	"bufio"         // bufio: buffered I/O operations (バッファリングされたI/O操作)
	"encoding/json" // encoding/json: JSON encoding and decoding (JSONエンコード・デコード)
	"fmt"           // fmt: formatted I/O (フォーマット済みI/O)
	"log"           // log: simple logging package (シンプルなログ記録パッケージ)
	"os"            // os: operating system interface (オペレーティングシステムインターフェース)
	"strings"       // strings: string manipulation functions (文字列操作関数)
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
// JSONRPCRequest: JSON-RPC 2.0リクエストを表現する構造体
// represents: 表現する、示す
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"` // jsonrpc: JSON-RPC protocol version (プロトコルバージョン)
	ID      interface{} `json:"id"`      // id: request identifier (リクエスト識別子)
	Method  string      `json:"method"`  // method: RPC method name (RPCメソッド名)
	Params  interface{} `json:"params"`  // params: method parameters (メソッドパラメータ)
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
// JSONRPCResponse: JSON-RPC 2.0レスポンスを表現する構造体
// response: 応答、返答
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`          // jsonrpc: JSON-RPC protocol version
	ID      interface{}   `json:"id"`               // id: matching request identifier
	Result  interface{}   `json:"result,omitempty"` // result: method result (メソッド結果)
	Error   *JSONRPCError `json:"error,omitempty"`  // error: error object (エラーオブジェクト)
}

// JSONRPCError represents a JSON-RPC 2.0 error
// JSONRPCError: JSON-RPC 2.0エラーを表現する構造体
// error: エラー、誤り
type JSONRPCError struct {
	Code    int         `json:"code"`           // code: error code (エラーコード)
	Message string      `json:"message"`        // message: error message (エラーメッセージ)
	Data    interface{} `json:"data,omitempty"` // data: additional error data (追加エラーデータ)
}

// MCPServer represents the MCP server instance
// MCPServer: MCPサーバーインスタンスを表現する構造体
// server: サーバー、提供者
// instance: インスタンス、実例
type MCPServer struct {
	name      string              // name: server name (サーバー名)
	version   string              // version: server version (サーバーバージョン)
	tools     map[string]Tool     // tools: available tools (利用可能なツール)
	resources map[string]Resource // resources: available resources (利用可能なリソース)
}

// Tool represents an MCP tool
// Tool: MCPツールを表現する構造体
// tool: ツール、道具
type Tool struct {
	Name        string      `json:"name"`        // name: tool name (ツール名)
	Description string      `json:"description"` // description: tool description (ツール説明)
	InputSchema interface{} `json:"inputSchema"` // inputSchema: input validation schema (入力検証スキーマ)
}

// Resource represents an MCP resource
// Resource: MCPリソースを表現する構造体
// resource: リソース、資源
type Resource struct {
	URI         string `json:"uri"`         // uri: resource URI (リソースURI)
	Name        string `json:"name"`        // name: resource name (リソース名)
	Description string `json:"description"` // description: resource description (リソース説明)
	MimeType    string `json:"mimeType"`    // mimeType: MIME type (MIMEタイプ)
}

// NewMCPServer creates a new MCP server instance
// NewMCPServer: 新しいMCPサーバーインスタンスを作成する関数
// creates: 作成する、生成する
func NewMCPServer(name, version string) *MCPServer {
	return &MCPServer{
		name:      name,
		version:   version,
		tools:     make(map[string]Tool),     // make: マップを初期化
		resources: make(map[string]Resource), // initialize: 初期化する
	}
}

// RegisterTool registers a new tool with the server
// RegisterTool: サーバーに新しいツールを登録する関数
// registers: 登録する、記録する
func (s *MCPServer) RegisterTool(tool Tool) {
	s.tools[tool.Name] = tool // assign: 割り当てる
}

// RegisterResource registers a new resource with the server
// RegisterResource: サーバーに新しいリソースを登録する関数
func (s *MCPServer) RegisterResource(resource Resource) {
	s.resources[resource.URI] = resource
}

// HandleRequest processes incoming JSON-RPC requests
// HandleRequest: 受信したJSON-RPCリクエストを処理する関数
// processes: 処理する、加工する
// incoming: 入ってくる、受信する
func (s *MCPServer) HandleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	// Input validation: セキュリティのための入力検証
	// validation: 検証、妥当性確認
	if req.JSONRPC != "2.0" {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32600,                     // Invalid Request (無効なリクエスト)
				Message: "Invalid JSON-RPC version", // version: バージョン
			},
		}
	}

	// Method dispatch: メソッドの振り分け
	// dispatch: 振り分ける、発送する
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "resources/list":
		return s.handleResourcesList(req)
	case "resources/read":
		return s.handleResourcesRead(req)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32601,             // Method not found (メソッドが見つからない)
				Message: "Method not found", // found: 見つかった
			},
		}
	}
}

// handleInitialize handles the initialize method
// handleInitialize: initializeメソッドを処理する関数
// handles: 処理する、扱う
func (s *MCPServer) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	// Server capabilities: サーバーの機能
	// capabilities: 機能、能力
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05", // protocol: プロトコル
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true, // listChanged: リスト変更通知
			},
			"resources": map[string]interface{}{
				"subscribe":   true, // subscribe: 購読する
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    s.name,    // name: 名前
			"version": s.version, // version: バージョン
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsList handles the tools/list method
// handleToolsList: tools/listメソッドを処理する関数
func (s *MCPServer) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	tools := make([]Tool, 0, len(s.tools)) // make: スライスを作成
	for _, tool := range s.tools {         // range: 範囲、レンジ
		tools = append(tools, tool) // append: 追加する
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"tools": tools},
	}
}

// handleToolsCall handles the tools/call method
// handleToolsCall: tools/callメソッドを処理する関数
func (s *MCPServer) handleToolsCall(req *JSONRPCRequest) *JSONRPCResponse {
	params, ok := req.Params.(map[string]interface{}) // type assertion: 型アサーション
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,               // Invalid params (無効なパラメータ)
				Message: "Invalid parameters", // parameters: パラメータ
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Tool name is required", // required: 必須の
			},
		}
	}

	// Security: ツール名の検証
	// security: セキュリティ、安全性
	if _, exists := s.tools[toolName]; !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Tool not found", // found: 見つかった
			},
		}
	}

	// Execute tool: ツールを実行
	// execute: 実行する、遂行する
	result := s.executeTool(toolName, params["arguments"])

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleResourcesList handles the resources/list method
// handleResourcesList: resources/listメソッドを処理する関数
func (s *MCPServer) handleResourcesList(req *JSONRPCRequest) *JSONRPCResponse {
	resources := make([]Resource, 0, len(s.resources))
	for _, resource := range s.resources {
		resources = append(resources, resource)
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"resources": resources},
	}
}

// handleResourcesRead handles the resources/read method
// handleResourcesRead: resources/readメソッドを処理する関数
func (s *MCPServer) handleResourcesRead(req *JSONRPCRequest) *JSONRPCResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
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
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "URI is required", // URI: Uniform Resource Identifier
			},
		}
	}

	// Security: URI validation
	// validation: 検証、妥当性確認
	if !strings.HasPrefix(uri, "file://") && !strings.HasPrefix(uri, "https://") {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid URI scheme", // scheme: スキーム、仕組み
			},
		}
	}

	// Read resource: リソースを読み取り
	content := s.readResource(uri)

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"uri":      uri,
					"mimeType": "text/plain", // mimeType: MIMEタイプ
					"text":     content,      // text: テキスト
				},
			},
		},
	}
}

// executeTool executes a specific tool
// executeTool: 特定のツールを実行する関数
// specific: 特定の、具体的な
func (s *MCPServer) executeTool(toolName string, arguments interface{}) map[string]interface{} {
	// Example tool execution: ツール実行の例
	// execution: 実行、遂行
	switch toolName {
	case "echo":
		args, ok := arguments.(map[string]interface{})
		if !ok {
			return map[string]interface{}{
				"error": "Invalid arguments", // arguments: 引数
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
					"text": fmt.Sprintf("Echo: %s", message), // sprintf: 文字列フォーマット
				},
			},
		}
	default:
		return map[string]interface{}{
			"error": "Unknown tool", // unknown: 不明な、未知の
		}
	}
}

// readResource reads a resource by URI
// readResource: URIによってリソースを読み取る関数
func (s *MCPServer) readResource(uri string) string {
	// Example resource reading: リソース読み取りの例
	// reading: 読み取り、読書
	switch {
	case strings.HasPrefix(uri, "file://"):
		// File system access: ファイルシステムアクセス
		// access: アクセス、接近
		return fmt.Sprintf("Content of %s", uri)
	case strings.HasPrefix(uri, "https://"):
		// HTTP request: HTTPリクエスト
		// request: リクエスト、要求
		return fmt.Sprintf("Web content of %s", uri)
	default:
		return "Unknown resource type" // type: タイプ、種類
	}
}

// Run starts the MCP server
// Run: MCPサーバーを開始する関数
// starts: 開始する、始める
func (s *MCPServer) Run() {
	scanner := bufio.NewScanner(os.Stdin) // scanner: スキャナー、読み取り器

	for scanner.Scan() { // scan: スキャンする、読み取る
		line := scanner.Text() // text: テキスト、文字列

		// Skip empty lines: 空行をスキップ
		// skip: スキップする、飛ばす
		// empty: 空の、からの
		if strings.TrimSpace(line) == "" {
			continue // continue: 続ける、継続する
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			// Log error: エラーをログに記録
			log.Printf("JSON parsing error: %v", err) // parsing: 解析
			continue
		}

		// Process request: リクエストを処理
		// process: 処理する、加工する
		resp := s.HandleRequest(&req)

		// Send response: レスポンスを送信
		// send: 送信する、送る
		if respData, err := json.Marshal(resp); err == nil {
			fmt.Println(string(respData)) // println: 行を出力
		} else {
			log.Printf("JSON marshaling error: %v", err) // marshaling: マーシャリング
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err) // scanner: スキャナー
	}
}

// main function: メイン関数
// main: メイン、主要な
// function: 関数、機能
func main() {
	// Create server: サーバーを作成
	// create: 作成する、生成する
	server := NewMCPServer("CustomMCPServer", "1.0.0")

	// Register tools: ツールを登録
	// register: 登録する、記録する
	server.RegisterTool(Tool{
		Name:        "echo",
		Description: "Echo back the provided message", // provided: 提供された
		InputSchema: map[string]interface{}{
			"type": "object", // object: オブジェクト、物体
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "Message to echo back", // echo: エコー、反響
				},
			},
			"required": []string{"message"}, // required: 必須の
		},
	})

	// Register resources: リソースを登録
	server.RegisterResource(Resource{
		URI:         "file:///example.txt",
		Name:        "Example File",         // example: 例、見本
		Description: "An example text file", // file: ファイル、書類
		MimeType:    "text/plain",           // plain: プレーン、平文
	})

	// Start server: サーバーを開始
	// start: 開始する、始める
	server.Run()
}
