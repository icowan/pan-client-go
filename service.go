/**
 * @Time : 2021/11/30 11:38 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package pan_client_go

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Middleware func(service Service) Service

type GenResult struct {
	Err     string `json:"err,omitempty"`
	Origin  string `json:"origin"`
	Ref     string `json:"ref,omitempty"`
	Success bool   `json:"success,omitempty"`
	URL     string `json:"url"`
}

type Service interface {
	// Gen 获取短链接
	Gen(ctx context.Context, name, bucket, targetPath, sharer string, expireTime *time.Time) (res GenResult, err error)
	// GenBatch 批量生成短链接
	GenBatch(ctx context.Context, name, bucket string, targetPath []string, sharer string, expireTime *time.Time) (res []GenResult, err error)
	// ExpiresTime 更新过期时间
	ExpiresTime(ctx context.Context, code string, expireTime time.Time) (err error)
	// Close 关闭连接
	Close(ctx context.Context) (err error)
}

func EncodeData(params map[string]interface{}) string {
	dataParams := url.Values{}
	for index, v := range params {
		encodeI(v, index, dataParams)
	}
	s, _ := url.QueryUnescape(dataParams.Encode())
	return s
}

func encodeI(i interface{}, parentKey string, values url.Values) {
	switch t := i.(type) {
	case string:
		values.Add(parentKey, t)
	case bool:
		values.Add(parentKey, strconv.FormatBool(t))
	case float64:
		values.Add(parentKey, strconv.FormatFloat(t, 'g', 10, 64))
	//case []string:
	//	for index, value := range t {
	//		values.Add(sliceKey(parentKey, index), value)
	//	}
	case []interface{}:
		for index, v := range t {
			encodeI(v, sliceKey(parentKey, index), values)
		}
	case map[string]interface{}:
		for index, v := range t {
			encodeI(v, nestKey(parentKey, index), values)
		}
	}
}

func sliceKey(key string, index int) string {
	return key + "." + strconv.Itoa(index)
}

func nestKey(key string, nestKey string) string {
	return key + "." + nestKey
}

func encodeSign(accessKey, secretKey, data string) string {
	//fmt.Println("encodeSign", fmt.Sprintf("%s:%s:%s", accessKey, data, secretKey))
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", accessKey, data, secretKey)))
	return hex.EncodeToString(h.Sum(nil))
}
