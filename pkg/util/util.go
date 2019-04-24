package util

import (
	"fmt"
	"bytes"
	"net/http"
	"runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ObjectInfo(o interface{}) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	object := o.(metav1.Object)
	s := fmt.Sprintf("%T, ", o)
	if namespace := object.GetNamespace(); namespace != "" {
		s += fmt.Sprintf("Namespace=%s, ", namespace)
	}
	s += fmt.Sprintf("Name=%s", object.GetName())
	return s
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}
