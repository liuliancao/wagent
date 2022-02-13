package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

func GetToken() (string, error) {
	// TODO: token has expire timestamp, so donnot need to update every times
	url := "http://op.liuliancao.com/auth"
	data := map[string]string{
		"username": "wagent",
		"password": "wagent",
	}
	result, err := HttpPOST(url, data, "application/json")
	if err != nil {
		log.Println(url, data, err)
		return "", err
	}
	log.Println("post token 返回" + result)

	token, err := simplejson.NewJson([]byte(result))
	t, err := token.Get("data").Get("token").String()
	if err != nil {
		return "", err
	}
	return t, nil
}

func HttpPOST(url string, data interface{}, contentType string) (res string, err error) {
	client := &http.Client{Timeout: 5 * time.Second}

	byte, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(string(byte)))
	if err != nil {
		log.Println(err)
		return
	}
	resp, _ := client.Do(request)
	if err != nil {
		log.Fatalf("post %s with %v get err %v", url, data, err)
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), nil
}

func HttpGET(url string) (res string, err error) {

	client := &http.Client{Timeout: 5 * time.Second}

	request, err := http.NewRequest("GET", url, nil)

	resp, _ := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}
func CMDBHttpPOST(url string, data interface{}, contentType string) (res string, err error) {
	token, err := GetToken()
	if err != nil {
		log.Fatalf("get token failed %v", err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	byte, err := json.Marshal(data)
	request, err := http.NewRequest("POST", url, strings.NewReader(string(byte)))

	request.Header.Add("Authorization", "Bearer "+token)
	resp, _ := client.Do(request)
	if err != nil {
		log.Fatalf("post %s with %v get err %v", url, data, err)
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), nil
}
func CMDBHttpPUT(url string, data interface{}, contentType string) (res string, err error) {
	token, err := GetToken()
	if err != nil {
		log.Fatalf("get token failed %v", err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	byte, err := json.Marshal(data)
	request, err := http.NewRequest("PUT", url, strings.NewReader(string(byte)))

	request.Header.Add("Authorization", "Bearer "+token)
	resp, _ := client.Do(request)
	if err != nil {
		log.Fatalf("put %s with %v get err %v", url, data, err)
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), nil
}

func CMDBHttpGET(url string) (res string, err error) {
	token, err := GetToken()
	if err != nil {
		log.Fatalf("get token failed %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", token)

	resp, _ := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}
