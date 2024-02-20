package main

import (
	"BHCEupload/internal"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type BHSessionOutputConfig struct {
	BatchSize int
}

type BHResponseData struct {
	Id            int                    `json:"id"`
	Status        int                    `json:"status"`
	StatusMessage string                 `json:"status_message"`
	Nodes         map[string]interface{} `json:"nodes"`
}

type BHResponse struct {
	Data BHResponseData `json:"data"`
}

func QueryBloodhoundAPI(uri string, method string, body []byte, creds internal.Credentials) (BHResponse, error) {
	// The first HMAC digest is the token key
	digester := hmac.New(sha256.New, []byte(creds.BHTokenKey))

	// OperationKey is the first HMAC digestresource
	digester.Write([]byte(fmt.Sprintf("%s%s", method, uri)))

	// Update the digester for further chaining
	digester = hmac.New(sha256.New, digester.Sum(nil))
	datetimeFormatted := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	digester.Write([]byte(datetimeFormatted[:13]))

	// Update the digester for further chaining
	digester = hmac.New(sha256.New, digester.Sum(nil))

	// Body signing is the last HMAC digest link in the signature chain. This encodes the request body as part of
	// the signature to prevent replay attacks that seek to modify the payload of a signed request. In the case
	// where there is no body content the HMAC digest is computed anyway, simply with no values written to the
	// digester.
	if body != nil {
		digester.Write(body)
	}

	bhendpoint := fmt.Sprintf("%s%s", creds.BHUrl, uri)

	// Perform the request with the signed and expected headers
	req, err := http.NewRequest(method, bhendpoint, bytes.NewBuffer(body))
	if err != nil {
		return BHResponse{}, err
	}

	req.Header.Set("User-Agent", "simple-uploader-v0.1")
	req.Header.Set("Authorization", fmt.Sprintf("bhesignature %s", creds.BHTokenID))
	req.Header.Set("RequestDate", datetimeFormatted)
	req.Header.Set("Signature", base64.StdEncoding.EncodeToString(digester.Sum(nil)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BHResponse{}, err
	}

	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return BHResponse{}, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return BHResponse{}, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}

	var response BHResponse
	// Empty response is OK for some endpoints
	if len(respbody) == 0 {
		return BHResponse{}, nil
	}
	err = json.Unmarshal(respbody, &response)
	if err != nil {
		return BHResponse{}, err
	}
	return response, nil
}

func UploadData(data []byte, creds internal.Credentials) error {
	upload_job, err := QueryBloodhoundAPI("/api/v2/file-upload/start", "POST", nil, creds)
	if err != nil {
		return err
	}
	job_id := upload_job.Data.Id
	log.Println("Processing job ID:", job_id)
	_, err = QueryBloodhoundAPI(fmt.Sprintf("/api/v2/file-upload/%d", job_id), "POST", data, creds)
	if err != nil {
		return err
	}
	_, err = QueryBloodhoundAPI(fmt.Sprintf("/api/v2/file-upload/%d/end", job_id), "POST", nil, creds)
	if err != nil {
		return err
	}
	log.Println("Data uploaded successfully for job ID:", job_id)
	return nil
}

func processFile(path string, creds internal.Credentials) error {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = UploadData(jsonFile, creds)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	url := flag.String("url", "http://localhost:8080", "Bloodhound URL")
	tokenid := flag.String("tokenid", "", "Bloodhound Token ID")
	tokenkey := flag.String("tokenkey", "", "Bloodhound Token Key")
	dir := flag.String("dir", "./", "Directory to process")

	flag.Parse()
	internal.Banner()

	if *tokenid == "" || *tokenkey == "" {
		fmt.Println("Please provide all required flags: -tokenid, -tokenkey")
		fmt.Println("OPTIONAL: -url, -dir, -h  for help")
		return
	}

	creds := internal.Credentials{
		BHUrl:      *url,
		BHTokenID:  *tokenid,
		BHTokenKey: *tokenkey,
	}

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			sizeMB := float64(info.Size()) / 1024.0 / 1024.0
			log.Printf("Uploading file %s, size: %.2f MB", path, sizeMB)
			if sizeMB > 20000 {
				log.Printf("File %s is quite large, will most likely fail, use chophound to make it smaller, skipping.", path)
				return nil
			} else {
				err := processFile(path, creds)
				if err != nil {
					fmt.Printf("Error processing file %s: %v\n", path, err)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", dir, err)
	}
}
