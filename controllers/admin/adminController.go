package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"kaixin/util/weixinUtil"
	"net/http"
	"strings"
)

type AdminKefuController struct {
	beego.Controller
}

var (
	AccessToken string = "{\"access_token\":\"HSLc1l8iI4f4U9L3UzCqxtyl1JOEOZ3B_THk6k0Xi_cTI2l4CfANoVZ2HqMVBE2GG8d0qqBLpPMi3ScHle15VFibC8K8Oeri7AKmBl8A3Ak\",\"expires_in\":7200}"
)

type AccessTokenObj struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
	ErrCode     float64 `json:"errcode"`
	Errmsg      string  `json:"errmsg`
}

type Kefu struct {
	Account  string `json:"kf_account"` //完整客服账号，格式为：账号前缀@公众号微信号，账号前缀最多10个字符，必须是英文或者数字字符。如果没有公众号微信号，请前往微信公众平台设置。
	NickName string `json:"nickname"`
	Password string `json:"password"`
}

func (c *AdminKefuController) Get() {
	c.Data["AccessToken"] = stringToAccessTokenObj(AccessToken).AccessToken
	c.TplNames = "admin/kefu.tpl"
}

func (c *AdminKefuController) GetAccessToken() {

	refreshToken, err := c.GetBool("refreshToken")
	if err != nil {
		refreshToken = false
	}

	if refreshToken {
		accessToken, err := weixinUtil.FetchAccessToken()
		if err != nil {
			return
		}
		AccessToken = accessToken
	}

	//	josnStr := AccessToken
	//	tokenObj := AccessTokenObj{}
	//	json.Unmarshal([]byte(josnStr), &tokenObj)

	tokenObj := stringToAccessTokenObj(AccessToken)

	c.Data["json"] = &tokenObj
	c.ServeJson()
}

func stringToAccessTokenObj(jsonStr string) AccessTokenObj {
	josnStr := AccessToken
	tokenObj := AccessTokenObj{}
	json.Unmarshal([]byte(josnStr), &tokenObj)
	return tokenObj
}

func (c *AdminKefuController) AddKefu() {
	tokenObj := stringToAccessTokenObj(AccessToken)
	account := c.GetString("account")
	nickName := c.GetString("nickName")
	password := c.GetString("password")
	kefu := Kefu{account, nickName, password}
	//kefu := Kefu{"jiansheng@gh_4193e9200053", "剑圣", "jiansheng"}

	body, err := json.MarshalIndent(kefu, "", "")
	if err != nil {
		return
	}
	postReq, err := http.NewRequest("POST",
		strings.Join([]string{"https://api.weixin.qq.com/customservice/kfaccount/add?access_token=", tokenObj.AccessToken}, ""),
		bytes.NewReader(body))
	postReq.Header.Set("Content-Type", "application/json;encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return
	}

	b1, err := ioutil.ReadAll(resp.Body)
	k1 := &AccessTokenObj{}
	json.Unmarshal(b1, k1)
	fmt.Println(k1)

	resp.Body.Close()

	c.Data["json"] = kefu
	c.ServeJson()
}

//微信ip服务器获取接口
type WxIps struct {
	Errcode float64  `json:"errcode"`
	Errmsg  string   `json:"errmsg"`
	IpList  []string `json:"ip_list"`
}

func (c *AdminKefuController) GetCallbackIp() {
	tokenObj := stringToAccessTokenObj(AccessToken)
	getRequest, err := http.NewRequest("GET", "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token="+tokenObj.AccessToken, nil)
	//getRequest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	if err != nil {
		return
	}
	client := &http.Client{}
	resp, err := client.Do(getRequest)

	if err != nil {
		return
	}
	ips := &WxIps{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, ips)
	c.Data["json"] = ips
	c.ServeJson()
}

//微信ip服务器获取接口

//发送客服消息接口message/custom/send

type CustomTextMessageContent struct {
	Content string `json:"content"`
}
type CustomTextMessage struct {
	ToUser  string                   `json:"touser"`
	MsgType string                   `json:"msgtype"`
	Text    CustomTextMessageContent `json:"text"`
}
type MsgResult struct {
	Errcode float64 `json:"errcode"`
	Errmsg  string  `json:"errmsg"`
}

//	@router:/admin/sendCustomMessage
func (c *AdminKefuController) SendCustomMessage() {
	tokenObj := stringToAccessTokenObj(AccessToken)
	toUser := c.GetString("toUser")
	msgType := c.GetString("msgType")
	msgContent := c.GetString("msgContent")
	msg := CustomTextMessage{
		ToUser:  toUser,
		MsgType: msgType,
		Text: CustomTextMessageContent{
			Content: msgContent,
		},
	}
	requestBody, err := json.Marshal(msg)
	if nil != err {
		return
	}
	request, err := http.NewRequest("POST",
		"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+tokenObj.AccessToken,
		bytes.NewReader(requestBody))
	if nil != err {
		return
	}
	request.Header.Set("Content-Type", "application/json;encoding=utf-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if nil != err {
		fmt.Println("-----------------")
		fmt.Println(err)
		fmt.Println("-----------------")
		c.Data["json"] = "789"
		c.ServeJson()
		return
	}
	sendMsgResult := &MsgResult{}
	responseBody, err := ioutil.ReadAll(response.Body)
	if nil != err {
		c.Data["json"] = "123"
		c.ServeJson()
		return
	}
	json.Unmarshal(responseBody, sendMsgResult)
	c.Data["json"] = sendMsgResult
	c.ServeJson()
}

//发送客服消息接口
