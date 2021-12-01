/**
 * @Time : 2021/11/30 11:49 AM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package pan_client_go

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kplcloud/request"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

type genResult struct {
	Code    int       `json:"code"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
	TraceId string    `json:"traceId"`
	Data    GenResult `json:"data"`
}

type batchGenResult struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId"`
	Data    []GenResult `json:"data"`
}

type genRequest struct {
	Name       string `json:"name"`
	Bucket     string `json:"bucket" valid:"required"`
	TargetPath string `json:"targetPath" valid:"required"`
	ValidTime  int64  `json:"validTime"` // 分钟
	Sharer     string `json:"sharer"`
	Public     bool   `json:"public"`
}

type batchGenRequest struct {
	Name       string   `json:"name" valid:"required"`
	Bucket     string   `json:"bucket" valid:"required"`
	TargetPath []string `json:"targetPath" valid:"required"`
	ValidTime  int64    `json:"validTime"` // 分钟
	Sharer     string   `json:"sharer" valid:"required"`
	Public     bool     `json:"public"`
}

type httpClient struct {
	host, accessKey, secretKey string
	client                     *http.Client
}

func (s *httpClient) Gen(ctx context.Context, name, bucket, targetPath, sharer string, expireTime *time.Time) (res GenResult, err error) {
	timestamp := time.Now().Unix()
	nonce := rand.Int()
	if expireTime == nil {
		t := time.Now()
		expireTime = &t
	}
	validTime := (expireTime.Unix() - time.Now().Unix()) * 60 // *60秒
	req := genRequest{
		Name:       name,
		Bucket:     bucket,
		TargetPath: targetPath,
		ValidTime:  validTime,
		Sharer:     sharer,
		Public:     true,
	}
	var data map[string]interface{}
	b, _ := json.Marshal(req)
	_ = json.Unmarshal(b, &data)

	sign := encodeSign(s.accessKey, s.secretKey, fmt.Sprintf("%s:%d:%d", EncodeData(data), timestamp, nonce))

	var rs genResult
	err = request.NewRequest(fmt.Sprintf("%s%s", s.host, "/share/gen"), http.MethodPost).
		Param("accessKey", s.accessKey).
		Param("sign", sign).
		Param("timestamp", strconv.Itoa(int(timestamp))).
		Param("nonce", strconv.Itoa(nonce)).
		Header("ContentType", "application/json").
		Body(b).
		Do().Into(&rs)
	if err != nil {
		return res, err
	}

	return rs.Data, nil
}

func (s *httpClient) GenBatch(ctx context.Context, name, bucket string, targetPath []string, sharer string, expireTime *time.Time) (res []GenResult, err error) {
	timestamp := time.Now().Unix()
	nonce := rand.Int()
	if expireTime == nil {
		t := time.Now()
		expireTime = &t
	}
	validTime := (expireTime.Unix() - time.Now().Unix()) * 60 // *60秒
	req := batchGenRequest{
		Name:       name,
		Bucket:     bucket,
		TargetPath: targetPath,
		ValidTime:  validTime,
		Sharer:     sharer,
		Public:     true,
	}
	var data map[string]interface{}
	b, _ := json.Marshal(req)
	_ = json.Unmarshal(b, &data)

	sign := encodeSign(s.accessKey, s.secretKey, fmt.Sprintf("%s:%d:%d", EncodeData(data), timestamp, nonce))

	var rs batchGenResult
	err = request.NewRequest(fmt.Sprintf("%s%s", s.host, "/share/batch/gen"), http.MethodPost).
		Param("accessKey", s.accessKey).
		Param("sign", sign).
		Param("timestamp", strconv.Itoa(int(timestamp))).
		Param("nonce", strconv.Itoa(nonce)).
		Header("ContentType", "application/json").
		Body(b).
		Do().Into(&rs)
	if err != nil {
		return res, err
	}

	return rs.Data, nil

}

func (s *httpClient) ExpiresTime(ctx context.Context, code string, expireTime time.Time) (err error) {
	panic("implement me")
}

func (s *httpClient) Close(ctx context.Context) (err error) {
	return nil
}

func NewHTTPClient(host, accessKey, secretKey string, client *http.Client) (svc Service, err error) {
	if client == nil {
		dialer := &net.Dialer{
			Timeout:   time.Duration(15 * int64(time.Second)),
			KeepAlive: time.Duration(15 * int64(time.Second)),
		}

		client = &http.Client{
			Transport: &http.Transport{
				DialContext: dialer.DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		}
	}
	return &httpClient{
		accessKey: accessKey,
		secretKey: secretKey,
		client:    client,
		host:      host,
	}, nil
}
