package http_client_pool

import (
	"errors"
	logs "github.com/cihub/seelog"
	"go_common_lib/cache"
	"io"
	"io/ioutil"
)

func Get_n_url(url string, n int) ([]byte, error) {
	client, isIn := ClientMap[n]
	if !isIn {
		ClientMap[n] = initClient(n)
	}
	if client == nil {
		client = initClient(n)
		ClientMap[n] = client
	}
	resp, err := client.Get(url)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return nil, errors.New("res code not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func Get_20_url(url string) ([]byte, error) {
	resp, err := Http_20_client.Get(url)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return nil, errors.New("res code not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func Get_50_url(url string) ([]byte, error) {
	resp, err := Http_50_client.Get(url)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return nil, errors.New("res code not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func GetUrlPostN(url string, reqbody io.Reader, n int) ([]byte, error) {
	client, isIn := ClientMap[n]
	if !isIn {
		ClientMap[n] = initClient(n)
	}
	if client == nil {
		client = initClient(n)
		ClientMap[n] = client
	}
	resp, err := client.Post(url, "application/x-www-form-urlencoded", reqbody)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode/200 > 1 {
		logs.Error("res code not 2xx:", resp.StatusCode)
		return nil, errors.New("res code not 2xx")
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func Get_url_post(url string, reqbody io.Reader) ([]byte, error) {
	resp, err := Http_client.Post(url, "application/x-www-form-urlencoded", reqbody)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return nil, errors.New("res code not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func Get_url(url string) ([]byte, error) {
	resp, err := Http_client.Get(url)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return nil, errors.New("res code not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func Get_url_cache(url string, key string, expire int64, is_hit *bool) ([]byte, error) {
	value, err := cache.Cache_controlor.Get(key)
	if err == nil && len(value) != 0 {
		*is_hit = true
		return ([]byte(value)), nil
	}
	resp, err := Http_20_client.Get(url)
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return nil, errors.New("res code not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	cache.Cache_controlor.Update(key, string(body), expire)
	return body, err
}

func GetUrlResultFromCache(key string) ([]byte, bool) {
	value, err := cache.Cache_controlor.Get(key)
	if err == nil && len(value) != 0 {
		return ([]byte(value)), true
	}

	return nil, false
}

func UpdateUrlResultCache(key string, result string, expire int64) {
	cache.Cache_controlor.Update(key, result, expire)
}
