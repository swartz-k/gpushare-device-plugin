package client

import (
	"flag"
	"fmt"
	"k8s.io/client-go/rest"
	"os"
	"testing"
	"time"
)

var (
	clientCert string
	clientKey  string
	token      string
	timeout    int
)

func TestNewKubeletClient(t *testing.T) {
	flag.StringVar(&clientCert, "client-cert", "", "")
	flag.StringVar(&clientKey, "client-key", "", "")
	flag.StringVar(&token, "token", "", "")
	flag.IntVar(&timeout, "timeout", 10, "")

	flag.Parse()

	if clientCert == "" && clientKey == "" && token == "" {
		tokenByte, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			t.Skipf("in cluster mode, find token failed, error: %v", err)
		}
		token = string(tokenByte)
	}

	c, err := NewKubeletClient(&KubeletClientConfig{
		Address: "127.0.0.1",
		Port:    10250,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure:   true,
			ServerName: "kubelet",
			CertFile:   clientCert,
			KeyFile:    clientKey,
		},
		BearerToken: token,
		HTTPTimeout: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	podsList, err := c.GetNodeRunningPods()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(podsList)
}
