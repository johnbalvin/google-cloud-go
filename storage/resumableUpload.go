package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2/jwt"
)

const bucketName = "myBucketName"            //your bucket name
const myWebPageURL = "http://localhost:8080" //as I´m working locally, this is my "webPageURL" for now, if you deploy your app, change to your page url

//this information you get it from console.cloud.google.com/apis/credentials
//select the option ->  "service account key", select your account
//it will download a json file with the information you need for this
var config = &jwt.Config{
	Email:        "someEmail@something.com",
	PrivateKey:   []byte("-----BEGIN PRIVATE KEY-----\nbablablablablablabalabal\nbablablablablablabalabalas\nbablablablablablabalabal\nbablablablablablabalabalas\nbablablablablablabalabalas\nbablablablablablabalabalas\nZ6ibLY8imBGf04PT4Z5vPFA=\n-----END PRIVATE KEY-----\n"),
	PrivateKeyID: "somethingVeryPrivate",
	Scopes:       []string{"https://www.googleapis.com/auth/devstorage.read_write"}, //you can change/add scopes if you want https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

//--------------------------

func main() {
	nameObject := "randomName" //be sure this is an unique name because if there's another one with the same name, it will overwrite it
	urlToUpload, err := resumeUploadURL(nameObject)
	if err != nil {
		log.Println("Err: ", err)
		return
	}
	fmt.Printf("\nURL:%s", urlToUpload)
	// use method: PUT , to upload the file with the url
}

func resumeUploadURL(objectName string) (string, error) {
	token, err := config.TokenSource(context.Background()).Token()
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	requestURL := "https://www.googleapis.com/upload/storage/v1/b/" + bucketName + "/o?uploadType=resumable&name=" + objectName
	pedir, _ := http.NewRequest("POST", requestURL, nil)
	pedir.Header.Add("Content-Length", "0")
	pedir.Header.Add("origin", myWebPageURL)
	pedir.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(pedir)
	if err != nil {
		return "", err
	}
	url, ok := resp.Header["Location"]
	if !ok && len(url) != 1 {
		return "", errors.New("NOT a valid url or there is not header")
	}
	fmt.Printf("\n\nBODY: %s\n\n", resp.Body) // I don´t know yet what the body is it for , I thought it should be empty
	return url[0], nil
}
