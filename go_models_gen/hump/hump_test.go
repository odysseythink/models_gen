package hump

import "testing"

func TestHump(t *testing.T) {
	ref := map[string]string{
		"hello_world":   "HelloWorld",
		"hello_ip":      "HelloIP",
		"guid_hello_ip": "GUIDHelloIP",
	}
	for k, v := range ref {
		v1 := BigHumpName(k)
		if v1 != v {
			t.Errorf("BigHumpName result(%s) is not match ref(%s)", v1, v)
		}
	}
}
