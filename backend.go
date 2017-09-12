package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/weaveworks/flux/guid"
	"io/ioutil"
	"net/http"
)

const (
	authURL  = "http://10.45.130.193/member/auth/grant?grantType=clientCredential&appId=test1&secret=NjE3RVpfQVBQ"
	loginURL = "http://10.45.130.193/member/user/loginByMobileAndPasswd?appToken="
	apiURL   = "http://182.254.229.224:10000/terminal/api"
)

var (
	appToken string
)

func httpLogin(courierTel, passwd string) error {

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", authURL, nil)

	reqest.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)")

	response, err := client.Do(reqest)

	if err != nil {
		return errors.New("")
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)

		js, _ := simplejson.NewJson(body)
		ret := js.Get("errcode").MustInt()
		if ret == 0 {

			appToken = js.Get("data").Get("appToken").MustString()

			lgInURL := loginURL + appToken

			jsLogin := simplejson.New()

			jsLogin.Set("version", "1")
			jsLogin.Set("appId", "MocCloud")
			jsLogin.Set("accessToken", "")
			jsLogin.Set("channel", "APP") //登陆渠道--APP,EZ:快递柜,WX:微信 ,WEB :web页面

			jsParam := simplejson.New()

			jsParam.Set("timeout", 20)
			jsParam.Set("mobile", courierTel)

			jsParam.Set("passwdMd5", md5Txt("123456"))
			jsParam.Set("utype", "Courier") //Customer-C端用户,Courier-快递员,Merchant-商户（必填）

			jsLogin.Set("params", jsParam)

			jsByte, _ := jsLogin.MarshalJSON()
			jsRet, jsErr := postJSON(lgInURL, string(jsByte))
			if jsErr != nil {
				return errors.New("")
			}
			errCode := jsRet.Get("errcode").MustInt()
			if errCode == 0 {

				return nil
			}

		}
	}

	return errors.New("")

}

func md5Txt(txt string) string {

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(txt))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func postJSON(url, jsTxt string) (*simplejson.Json, error) {

	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(jsTxt)))
	if err != nil {
		return nil, errors.New("")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {

	}
	body, _ := ioutil.ReadAll(res.Body)

	js, _ := simplejson.NewJson(body)
	return js, nil

}

func httpBookBox(sn, pkgID, boxType, takenMobile, postmanMobile string) (*simplejson.Json, error) {

	js := simplejson.New()
	js.Set("uid", "1234455")
	js.Set("service", "BookBox")
	js.Set("timeout", "18000")
	js.Set("retry", "1")
	js.Set("sn", sn)
	js.Set("requestId", guid.New())

	jsData := simplejson.New()

	jsData.Set("packageId", pkgID)
	jsData.Set("boxType", boxType)
	jsData.Set("bookType", "postman")
	jsData.Set("takeMobile", takenMobile) //取件人手机
	jsData.Set("postmanMobile", postmanMobile)
	jsData.Set("callbackUrl", "http://ip:port/approot?code=123456")
	jsData.Set("bookTimeSpan", 45)
	jsData.Set("channelId", 1001)
	jsData.Set("lockTime", 10)

	js.Set("data", jsData)

	if jsTxt, err := js.MarshalJSON(); err == nil {
		jsRes, err := postJSON(apiURL, string(jsTxt))

		return jsRes, err
	}

	return nil, errors.New("")

}

func httpCancalBook(sn, parcelID string) (*simplejson.Json, error) {

	js := simplejson.New()
	js.Set("uid", "1234455")
	js.Set("service", "BookBoxCancel")
	js.Set("timeout", "18000")
	js.Set("retry", "1")
	js.Set("sn", sn)

	js.Set("requestId", guid.New())

	jsData := simplejson.New()

	jsData.Set("parcelId", parcelID) //包裹ID

	js.Set("data", jsData)
	if jsTxt, err := js.MarshalJSON(); err == nil {
		jsRes, err := postJSON(apiURL, string(jsTxt))

		return jsRes, err
	}

	return nil, errors.New("")

}
