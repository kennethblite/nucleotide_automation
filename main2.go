package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
        "crypto/tls"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open("./hairpin_ct_file/"+path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err  :=  http.NewRequest("POST", uri, body)
	if err != nil{
		return nil, err
	}
	//PHPSESSID=caurbd1e9jb93d2murn0s6rns7
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", "PHPSESSID=jdsir7vbcm6uofve3sv9obp1f6")
	req.Header.Set("Referer", "https://rna-inforna.florida.scripps.edu/index.php")
	req.Header.Set("Origin", "https://rna-inforna.florida.scripps.edu")
	return req, nil
}

func RunTask( name string) {
	writename := strings.TrimSuffix(name, ".ct")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	extraParams := map[string]string{
		"ct_search_options":       "sequence_only",
	}
	///disneydb/query.php HTTP/1.1
	//request, err := newfileUploadRequest("http://localhost:8000", extraParams, "ctfile", "MI0000060.ct")
	request, err := newfileUploadRequest("https://rna-inforna.florida.scripps.edu/disneydb/excel.php", extraParams, "ctfile", name)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} 
	var bodyContent []byte
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header)
	bodyContent , err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//func WriteFile(filename string, data []byte, perm os.FileMode) error
	err = ioutil.WriteFile("./finished/"+writename+"sequence_only.xls", bodyContent, 0666)
	if err != nil{
		fmt.Println(err)
	}
	extraParams = map[string]string{
		"ct_search_options":       "sequence_plus_1",
	}
	///disneydb/query.php HTTP/1.1
	//request, err := newfileUploadRequest("http://localhost:8000", extraParams, "ctfile", "MI0000060.ct")
	request, err = newfileUploadRequest("https://rna-inforna.florida.scripps.edu/disneydb/excel.php", extraParams, "ctfile", name)
	if err != nil {
		log.Fatal(err)
	}
	client = &http.Client{}
	resp, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	} 
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header)
	bodyContent , err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//func WriteFile(filename string, data []byte, perm os.FileMode) error
	err = ioutil.WriteFile("./finished/"+writename+"sequence_plus_1.xls", bodyContent, 0666)
	if err != nil{
		fmt.Println(err)
	}
}

func main(){
    files, err := ioutil.ReadDir("./hairpin_ct_file")
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
            fmt.Println(f.Name())
	   RunTask(f.Name()) 
    }
}
