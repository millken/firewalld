package server

type ServerConfig struct {
	CmdApiAddr  string `json:"cmd_api_addr"`
	HttpApiAddr string `json:"http_api_addr"`
}
