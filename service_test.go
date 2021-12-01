/**
 * @Time : 2021/12/1 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */

package pan_client_go

import (
	"context"
	"testing"
	"time"
)

func initGRPCClient() (svc Service, err error) {
	return NewGRpcClient(
		"localhost:8082",
		"CXPUBXHF1SQOJL44DYZL",
		"HizofcKVinLfT71xXnK8HgsGG80LJhxRMurjMm2N",
	)
}

func initHTTPClient() (svc Service, err error) {
	return NewHTTPClient(
		"localhost:8080",
		"",
		"",
		nil,
	)
}
func TestGrpcClient_Gen(t *testing.T) {
	client, err := initGRPCClient()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	tm := time.Unix(time.Now().Unix()+3600, 0)
	gen, err := client.Gen(ctx, "dev", "dev-public-read", "03.pdf", "test", &tm)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(gen.URL)
	_ = client.Close(ctx)
}

func TestGrpcClient_GenBatch(t *testing.T) {
	client, err := initGRPCClient()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	tm := time.Unix(time.Now().Unix()+3600, 0)
	gen, err := client.GenBatch(ctx, "dev", "dev-public-read", []string{
		"03.pdf",
		"1591344049651.822021.jpg",
	}, "test", &tm)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range gen {
		t.Log(v.Ref, v.URL)
	}

	_ = client.Close(ctx)
}

func TestGrpcClient_ExpiresTime(t *testing.T) {
	client, err := initGRPCClient()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	tm := time.Unix(time.Now().Unix()+60, 0)
	err = client.ExpiresTime(ctx, "9fdPFCtnRyXI", tm)
	if err != nil {
		t.Error(err)
		return
	}
	_ = client.Close(ctx)
}

func TestHttpClient_Gen(t *testing.T) {
	client, err := initHTTPClient()
	if err != nil {
		t.Error(err)
		return
	}
	tm := time.Unix(time.Now().Unix()+3600, 0)
	ctx := context.Background()
	gen, err := client.Gen(ctx, "dev", "dev-public-read", "03.pdf", "test", &tm)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(gen.URL)
}

func TestHttpClient_GenBatch(t *testing.T) {
	client, err := initHTTPClient()
	if err != nil {
		t.Error(err)
		return
	}
	tm := time.Unix(time.Now().Unix()+3600, 0)
	ctx := context.Background()
	gen, err := client.GenBatch(ctx, "dev", "dev-public-read", []string{
		"03.pdf",
		"1591344049651.822021.jpg",
	}, "test", &tm)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range gen {
		t.Log(v.Ref, v.URL)
	}
}
