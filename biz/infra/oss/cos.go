package oss

import (
	"github.com/li1553770945/personal-file-service/biz/infra/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

func NewCosClient(conf *config.Config) *cos.Client {
	u, _ := url.Parse(conf.CosConfig.Endpoint)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.CosConfig.Ak, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: conf.CosConfig.Sk, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return client
}
