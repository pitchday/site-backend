package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cheviz/pitchdayBackend/config"
	"io"
	"net/http"
	"net/url"
	"os"
)

func ImportImage(imageUrl string, userId string) (s3Url string, err error) {
	requestUrl, err := url.Parse(imageUrl)
	if err != nil {
		Logger.Println("There was an error parsing the imageUrl")
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestUrl.String(), nil)

	//Do request
	resp, err := client.Do(req)
	if err != nil {
		Logger.Printf("There was an error downloading the requested image: %s\n", err)
		return
	}
	defer resp.Body.Close()

	imageName := fmt.Sprintf("%s%s", userId, "jpg")
	imagePath := fmt.Sprintf("./uploads/%s", imageName)

	//create file to temporarily store image
	file, err := os.OpenFile(imagePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer file.Close()

	io.Copy(file, resp.Body)
	defer os.Remove(imagePath)

	s3Url, err = putImageInBucket(imagePath, userId)

	return
}

func putImageInBucket(imagePath string, userId string) (url string, err error) {
	credentials := credentials.NewStaticCredentials(config.Conf.AWSAccessKey, config.Conf.AWSSecretKey, "")

	conf := aws.NewConfig()
	conf.WithCredentials(credentials)
	conf.WithRegion("us-east-1")

	svc := session.Must(session.NewSession(conf))

	uploader := s3manager.NewUploader(svc)

	f, err := os.Open(imagePath)
	if err != nil {
		err = fmt.Errorf("failed to open file %q, %v", imagePath, err)
		Logger.Println(err)
		return
	}

	fileKey := fmt.Sprintf("contributors/%s.jpg", userId)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.Conf.AWSBucketName),
		ContentType: aws.String("image/jpeg"),
		Key:    aws.String(fileKey),
		Body:   f,
	})
	if err != nil {
		err = fmt.Errorf("failed to upload file %q, %v", imagePath, err)
		Logger.Println(err)
		return
	}

	url = fmt.Sprintf("https://s3.amazonaws.com/pitchday/%s", fileKey)

	return
}
