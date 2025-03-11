package core

type Config struct {
	Vendor     string `json:"vendor" yaml:"vendor"`
	Endpoint   string `json:"endpoint" yaml:"endpoint"`
	AccessKey  string `json:"accessKey" yaml:"accessKey"`
	SecretKey  string `json:"secretKey" yaml:"secretKey"`
	BucketName string `json:"bucketName" yaml:"bucketName"`
}
