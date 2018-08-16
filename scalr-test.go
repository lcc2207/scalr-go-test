package main

import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "time"
import "crypto/hmac"
import "crypto/sha256"
import "encoding/base64"
import "strings"

const scalrSignatureVersion = "V1-HMAC-SHA256"
const APIKeyID = "xxxx"
const APIKeySecret = "xxxx"
const ScalrUrl = "http://demo.scalr.club/"

// generate signature
func computeHmac256(stringToSign string, APIKeySecret string) string {
	key := []byte(APIKeySecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(stringToSign))
  signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
  return scalrSignatureVersion+" "+signature
}

// build url for processing
func buildurl(queryurl string) string{
  return strings.Join([]string{ScalrUrl, queryurl}, "")
}

// build the signing string
func signstr(timestamp string, queryurl string) string{
    return strings.Join([]string{"GET", timestamp, queryurl, "", ""}, "\n")
}

// process the request
func processrequest(APIKeyID string, timestamp string, xsig string, queryurl string, raction string) []byte{
  // build url
  url := buildurl(queryurl)

  // create HTTP connection
  client := &http.Client{}
  req, err := http.NewRequest(raction, url, nil)
  // set Headers
  req.Header.Set("X-Scalr-Key-Id", APIKeyID )
  req.Header.Set("X-Scalr-Date", timestamp)
  req.Header.Set("X-Scalr-Signature", xsig)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-Scalr-Debug", "1")

  // execute request
  res, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
  return responseData

}

func main() {
  queryurl := "/api/v1beta0/user/6/projects/"
  timestamp := time.Now().Format(time.RFC3339)
  stringToSign := signstr(timestamp, queryurl)
  xsig := computeHmac256(stringToSign, APIKeySecret)
  dowork := processrequest(APIKeyID, timestamp, xsig, queryurl, "GET")
  fmt.Println(string(dowork))

}
