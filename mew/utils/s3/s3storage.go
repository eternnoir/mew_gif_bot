package s3

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

type S3Service struct {
	Key    string
	Secret string
	Bucket string
}

func NewS3Service(key, secret, bucket string) *S3Service {
	return &S3Service{Key: key, Secret: secret, Bucket: bucket}
}

func (s3s *S3Service) getConnection() *s3.S3 {
	AWSAuth := aws.Auth{
		AccessKey: s3s.Key, // change this to yours
		SecretKey: s3s.Secret,
	}
	region := aws.USEast
	return s3.New(AWSAuth, region)
}

func (s3s *S3Service) PutFile(buffer []byte, filepath string) error {
	return nil
}

func (s3s *S3Service) GetFilesUrl(path string) ([]string, error) {
	conn := s3s.getConnection()
	conn.Bucket(s3s.Bucket).URL("gifs")
	resps, err := conn.Bucket(s3s.Bucket).List(path, "/", "s", 1000)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, k := range resps.Contents {
		ret = append(ret, s3s.getUrl(path+"/"+k.Key, conn))
	}

	return ret, nil
}

func (s3s *S3Service) getUrl(filepath string, s3 *s3.S3) string {
	return s3.Bucket(s3s.Bucket).URL(filepath)
}
