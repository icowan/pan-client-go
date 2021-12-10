/**
 * @Time : 2021/11/30 11:49 AM
 * @Author : solacowa@gmail.com
 * @File : grpc
 * @Software: GoLand
 */

package pan_client_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/icowan/pan-client-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"strconv"
	"time"
)

type grpcClient struct {
	conn                 *grpc.ClientConn
	shareSvc             pb.ShareClient
	accessKey, secretKey string
}

func (s *grpcClient) Gen(ctx context.Context, name, bucket, targetPath, sharer string, expireTime *time.Time) (res GenResult, err error) {
	timestamp := time.Now().Unix()
	nonce := rand.Int()
	if expireTime == nil {
		t := time.Now()
		expireTime = &t
	}
	validTime := (expireTime.Unix() - time.Now().Unix()) / 60 // *60秒

	req := pb.GenRequest{
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

	md := metadata.Pairs(
		"accessKey", s.accessKey,
		"sign", sign,
		"timestamp", strconv.Itoa(int(timestamp)),
		"nonce", strconv.Itoa(nonce),
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	rs, err := s.shareSvc.Gen(ctx, &req)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(rs.GetData()), &res)
	return
}

func (s *grpcClient) GenBatch(ctx context.Context, name, bucket string, targetPath []string, sharer string, expireTime *time.Time) (res []GenResult, err error) {
	timestamp := time.Now().Unix()
	nonce := rand.Int()
	if expireTime == nil {
		t := time.Now()
		expireTime = &t
	}
	validTime := (expireTime.Unix() - time.Now().Unix()) * 60 // *60秒

	req := pb.BatchGenRequest{
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

	md := metadata.Pairs(
		"accessKey", s.accessKey,
		"sign", sign,
		"timestamp", strconv.Itoa(int(timestamp)),
		"nonce", strconv.Itoa(nonce),
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	rs, err := s.shareSvc.BatchGen(ctx, &req)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(rs.GetData()), &res)
	return
}

func (s *grpcClient) ExpiresTime(ctx context.Context, code string, expireTime time.Time) (err error) {
	timestamp := time.Now().Unix()
	nonce := rand.Int()
	validTime := (expireTime.Unix() - time.Now().Unix()) * 60 // *60秒

	req := pb.ExpiresRequest{
		Code:       code,
		ExtendTime: validTime,
	}

	var data map[string]interface{}
	b, _ := json.Marshal(req)
	_ = json.Unmarshal(b, &data)

	sign := encodeSign(s.accessKey, s.secretKey, fmt.Sprintf("%s:%d:%d", EncodeData(data), timestamp, nonce))

	md := metadata.Pairs(
		"accessKey", s.accessKey,
		"sign", sign,
		"timestamp", strconv.Itoa(int(timestamp)),
		"nonce", strconv.Itoa(nonce),
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err = s.shareSvc.Expires(ctx, &req)
	if err != nil {
		return
	}
	return nil
}
func (s *grpcClient) Close(ctx context.Context) (err error) {
	return s.conn.Close()
}

func NewGRpcClient(host, accessKey, secretKey string) (svc Service, err error) {
	ctx, cel := context.WithTimeout(context.Background(), time.Second*10)
	defer cel()
	conn, err := grpc.DialContext(ctx, host,
		grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		accessKey: accessKey,
		secretKey: secretKey,
		conn:      conn,
		shareSvc:  pb.NewShareClient(conn),
	}, nil
}
