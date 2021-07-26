package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func MakeRequest(baseUrl string, path string, method string, requestType string, data interface{}) ([]byte, error) {
	values, err := structToMap(data)
	if err != nil {
		return nil, err
	}

	log.Println(values.Encode())
	log.Println(url.ParseQuery(values.Encode()))

	request, err := http.NewRequest(strings.ToUpper(method), fmt.Sprintf("%s/%s", baseUrl, path), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = values.Encode()

	switch requestType {
	case "x-www-form-urlencoded":
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	case "form-data":
		request.Header.Add("Content-Type", "application/form-data")
		request.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	case "data-urlencode":
		request.Header.Add("Content-Type", "application/data-urlencode")
		request.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	default:
		return nil, fmt.Errorf("not support request type %s", requestType)
	}

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func structToMap(inter interface{}) (url.Values, error) {
	values := url.Values{}

	val := reflect.ValueOf(inter).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("name")

		switch typeField.Type.Kind() {
		case reflect.String:
			value := fmt.Sprintf("%s", valueField.Interface())
			values.Add(tag, value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := fmt.Sprintf("%d", valueField.Interface())
			values.Add(tag, value)
		case reflect.Float32, reflect.Float64:
			value := fmt.Sprintf("%d", valueField.Interface())
			values.Add(tag, value)
		case reflect.Bool:
			value := fmt.Sprintf("%t", valueField.Interface())
			values.Add(tag, value)
		case reflect.Slice:
			for index := 0; index < valueField.Len(); index++ {
				item := valueField.Index(index)
				switch item.Kind() {
				case reflect.String:
					value := fmt.Sprintf("%s", item.Interface())
					values.Add(tag, value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value := fmt.Sprintf("%d", valueField.Interface())
					values.Add(tag, value)
				case reflect.Float32, reflect.Float64:
					value := fmt.Sprintf("%d", valueField.Interface())
					values.Add(tag, value)
				case reflect.Bool:
					value := fmt.Sprintf("%t", valueField.Interface())
					values.Add(tag, value)
				default:
					return nil, fmt.Errorf("slice type %s not support", typeField.Name)
				}
			}
		default:
			return nil, fmt.Errorf("type %s not support", typeField.Name)
		}
	}
	return values, nil
}
