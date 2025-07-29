package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// -- decoding

func DecodeStreamedV(r io.Reader, to any) error {
	return json.NewDecoder(r).Decode(to)
}

// --- request
func ReqWithJSON(req *http.Request, o any) error {
	req.Header.Add("Content-Type", "application/json")
	bodyBytes, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("marshal body err: %v", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return nil
}

// --- response
type RespStatus string

const (
	RespOk   RespStatus = "ok"
	RespFail RespStatus = "fails"
)

type Resp struct {
	Status RespStatus `json:"status"`
	Data   any        `json:"data"`
}

func NewOkResp(data any) Resp {
	return Resp{Status: RespOk, Data: data}
}
func NewFailResp() Resp {
	return Resp{Status: RespFail, Data: nil}
}

func (r Resp) Bytes() []byte {
	w := new(bytes.Buffer)
	json.NewEncoder(w).Encode(r)
	return w.Bytes()
}

// --- error response
func WriteFailResp(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Write(NewFailResp().Bytes())
}

func WriteInternalErrResp(w http.ResponseWriter) {
	WriteFailResp(w, http.StatusInternalServerError)
}

func WriteNotfoundResp(w http.ResponseWriter) {
	WriteFailResp(w, http.StatusNotFound)
}

func WriteUnauthResp(w http.ResponseWriter) {
	WriteFailResp(w, http.StatusUnauthorized)
}

// --- success response
func WriteOKResp(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusOK)
	w.Write(NewOkResp(v).Bytes())
}
