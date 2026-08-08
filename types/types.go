package types

type ErrorResp struct {
    Code    uint16 `json:"code"`
    Message string `json:"message"`
}

