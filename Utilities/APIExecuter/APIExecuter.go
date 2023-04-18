package Utilities

import (
	"time"

	"gopkg.in/resty.v1"
)

//GET : REST GET
func GET(apirequest APIRequest) APIResponse {
	resp := new(APIResponse)

	client := resty.New()
	//client := resty.New()
	client.SetTimeout(time.Duration(apirequest.TimeOut) * time.Millisecond)

	if apirequest.IsQueryParamsExists {
		client.SetQueryParams(apirequest.RequestQueryParams)
	}

	if apirequest.IsAuthRequired {
		client = client.SetBasicAuth(apirequest.Authentication.Username, apirequest.Authentication.Password)
	}

	if apirequest.IsURLParamsExists {
		client.SetPathParams(apirequest.RequestURLParams)
	}

	if apirequest.IsProxySet {
		client.SetProxy(apirequest.ProxyAddress)
	}

	endpoint := apirequest.BaseURL + apirequest.Action

	response, err := client.R().SetContext(apirequest.ReqCtx).SetHeaders(apirequest.Headers).Get(endpoint)

	//fmt.Printf("\nResponse Status Code: %v", response.StatusCode())
	//fmt.Printf("\nResponse Time: %v", response.Time())
	//fmt.Printf("\nResponse Status: %v", response.Status())
	//fmt.Printf("\nResponse Body: %v", response)
	//fmt.Printf("\nResponse Received At: %v", response.ReceivedAt())

	resp.Err = err
	resp.Response = response.String()
	resp.Status = response.StatusCode()
	return *resp
}

//POST : REST POST
func POST(apirequest APIRequest, requestBody interface{}) APIResponse {
	resp := new(APIResponse)

	client := resty.New()
	client.SetTimeout(time.Duration(apirequest.TimeOut) * time.Millisecond)
	if apirequest.IsQueryParamsExists {
		client.SetQueryParams(apirequest.RequestQueryParams)
	}

	// if headers != nil {
	// 	client.Header = headers
	// }

	if apirequest.IsAuthRequired {
		client.SetBasicAuth(apirequest.Authentication.Username, apirequest.Authentication.Password)
	}
	if apirequest.IsURLParamsExists {
		client.SetPathParams(apirequest.RequestURLParams)
	}
	client.Header = apirequest.HttpHeaders

	endpoint := apirequest.BaseURL + apirequest.Action

	//some, _ := jsoniter.Marshal(requestBody)
	//fmt.Println(string(some))
	if apirequest.IsFormdataExists {
		response, err := client.R().SetContext(apirequest.ReqCtx).SetHeaders(apirequest.Headers).SetFormData(apirequest.Formdata).Post(endpoint)
		resp.Err = err
		resp.Response = response.String()
		resp.Status = response.StatusCode()
	} else {
		response, err := client.R().SetContext(apirequest.ReqCtx).SetHeaders(apirequest.Headers).SetBody(requestBody).Post(endpoint)
		resp.Err = err
		resp.Response = response.String()
		resp.Status = response.StatusCode()
	}

	//fmt.Prwintln("Response Status Code: ", response.StatusCode())
	//fmt.Println("Response Time: ", response.Time())
	//fmt.Println("Response Status: ", response.Status())
	//fmt.Println("Response Body: ", response)
	//fmt.Println("Response Received At: ", response.ReceivedAt())

	return *resp
}

//PUT : REST
func PUT(apirequest APIRequest, requestBody interface{}) APIResponse {
	resp := new(APIResponse)
	client := resty.New()
	client.SetTimeout(time.Duration(apirequest.TimeOut) * time.Millisecond)
	if apirequest.IsQueryParamsExists {
		client.SetQueryParams(apirequest.RequestQueryParams)
	}
	if apirequest.IsURLParamsExists {
		client.SetPathParams(apirequest.RequestURLParams)
	}

	if apirequest.IsAuthRequired {
		client = client.SetBasicAuth(apirequest.Authentication.Username, apirequest.Authentication.Password)
	}

	endpoint := apirequest.BaseURL + apirequest.Action

	response, err := client.R().SetHeaders(apirequest.Headers).SetBody(requestBody).Put(endpoint)
	resp.Err = err
	resp.Response = response.String()
	resp.Status = response.StatusCode()

	return *resp
}

//DELETE : REST DELETE
func DELETE(apirequest APIRequest) APIResponse {
	resp := new(APIResponse)
	client := resty.New()
	client.SetTimeout(time.Duration(apirequest.TimeOut) * time.Millisecond)

	if apirequest.IsQueryParamsExists {
		client.SetQueryParams(apirequest.RequestQueryParams)
	}

	if apirequest.IsAuthRequired {
		client = client.SetBasicAuth(apirequest.Authentication.Username, apirequest.Authentication.Password)
	}

	if apirequest.IsURLParamsExists {
		client.SetPathParams(apirequest.RequestURLParams)
	}

	endpoint := apirequest.BaseURL + apirequest.Action
	response, err := client.R().SetHeaders(apirequest.Headers).Delete(endpoint)

	//fmt.Printf("\nResponse Status Code: %v", response.StatusCode())
	//fmt.Printf("\nResponse Time: %v", response.Time())
	//fmt.Printf("\nResponse Status: %v", response.Status())
	//fmt.Printf("\nResponse Body: %v", response)
	//fmt.Printf("\nResponse Received At: %v", response.ReceivedAt())

	resp.Err = err
	resp.Response = response.String()
	resp.Status = response.StatusCode()
	return *resp
}

//PATCH : REST UPDATE
func PATCH(apirequest APIRequest, requestBody interface{}) APIResponse {
	resp := new(APIResponse)
	client := resty.New()
	client.SetTimeout(time.Duration(apirequest.TimeOut) * time.Second)
	if apirequest.IsQueryParamsExists {
		client.SetQueryParams(apirequest.RequestQueryParams)
	}
	if apirequest.IsURLParamsExists {
		client.SetPathParams(apirequest.RequestURLParams)
	}

	if apirequest.IsAuthRequired {
		client = client.SetBasicAuth(apirequest.Authentication.Username, apirequest.Authentication.Password)
	}

	endpoint := apirequest.BaseURL + apirequest.Action
	response, err := client.R().SetHeaders(apirequest.Headers).SetBody(requestBody).Patch(endpoint)

	resp.Err = err
	resp.Response = response.String()
	resp.Status = response.StatusCode()

	return *resp
}

//DELETE : REST DELETE WITH BODY
func DELETE_BODY(apirequest APIRequest, requestBody interface{}) APIResponse {
	resp := new(APIResponse)
	client := resty.New()
	client.SetTimeout(time.Duration(apirequest.TimeOut) * time.Millisecond)
	if apirequest.IsQueryParamsExists {
		client.SetQueryParams(apirequest.RequestQueryParams)
	}
	if apirequest.IsURLParamsExists {
		client.SetPathParams(apirequest.RequestURLParams)
	}

	if apirequest.IsAuthRequired {
		client = client.SetBasicAuth(apirequest.Authentication.Username, apirequest.Authentication.Password)
	}

	endpoint := apirequest.BaseURL + apirequest.Action

	response, err := client.R().SetHeaders(apirequest.Headers).SetBody(requestBody).Delete(endpoint)
	resp.Err = err
	resp.Response = response.String()
	resp.Status = response.StatusCode()

	return *resp
}
