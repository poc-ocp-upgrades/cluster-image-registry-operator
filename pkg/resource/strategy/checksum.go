package strategy

import (
	"crypto/sha256"
	"bytes"
	"net/http"
	"runtime"
	"encoding/json"
	"fmt"
)

func Checksum(o interface{}) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%x", sha256.Sum256(data)), nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}
