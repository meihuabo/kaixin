package controllers

/**
 *微信测试
 */
import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"log"
	"time"
)

const TOKEN string = "kaixinmeishimeike666godlike"

type MainController struct {
	beego.Controller
}

type MsgType struct {
	MsgType string
}

type TextRequestBody struct {
	//XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
	//语音消息的参数
	MediaId     string
	Format      string
	Recognition string
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	//语音消息的参数
	MediaId string
}
type VoiceResponseBody struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	MediaId      string
}

func (c *MainController) Get() {
	//	c.Data["Website"] = "beego.me"
	//	c.Data["Email"] = "astaxie@gmail.com"
	//c.TplNames = "index.tpl"

	//获取参数
	signature := c.GetString("signature")
	echostr := c.GetString("echostr")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	checkSignature(signature, timestamp, nonce)
	c.Ctx.WriteString(echostr)
}

func (c *MainController) Post() {
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if nil != err {
		return
	}
	c.Ctx.WriteString(fmt.Sprintln(body))
	requestBody := TextRequestBody{}
	xml.Unmarshal(body, &requestBody)
	fmt.Printf(fmt.Sprint(requestBody))
	c.Ctx.WriteString(fmt.Sprintln(requestBody))

	//回复用户消息
	if &requestBody != nil {
		fmt.Printf("Wechat Service: Recv text msg [%s] from user [%s]!",
			requestBody.Content,
			requestBody.FromUserName)

		if requestBody.MsgType == "text" && "你叫什么" == requestBody.Content {
			responseTextBody, err := makeTextResponseBody(requestBody.ToUserName,
				requestBody.FromUserName,
				"你好，我是开心每时每刻，欢迎你关注我。")
			if err != nil {
				log.Println("Wechat Service: makeTextResponseBody error: ", err)
				return
			}
			fmt.Fprintf(c.Ctx.ResponseWriter, string(responseTextBody))
		} else if requestBody.MsgType == "text" && "你多大了" == requestBody.Content {
			responseTextBody, err := makeTextResponseBody(requestBody.ToUserName,
				requestBody.FromUserName,
				"我才刚出生")
			if err != nil {
				log.Println("Wechat Service: makeTextResponseBody error: ", err)
				return
			}
			fmt.Fprintf(c.Ctx.ResponseWriter, string(responseTextBody))

		} else if requestBody.MsgType == "text" {
			responseTextBody, err := makeTextResponseBody(requestBody.ToUserName,
				requestBody.FromUserName,
				"欢迎来到英雄联盟！")
			if err != nil {
				log.Println("Wechat Service: makeTextResponseBody error: ", err)
				return
			}
			fmt.Fprintf(c.Ctx.ResponseWriter, string(responseTextBody))
		} else if requestBody.MsgType == "voice" {
			responseTextBody, err := makeTextResponseBody(requestBody.ToUserName,
				requestBody.FromUserName,
				requestBody.Recognition)
			if err != nil {
				log.Println("Wechat Service: makeTextResponseBody error: ", err)
				return
			}
			fmt.Fprintf(c.Ctx.ResponseWriter, string(responseTextBody))
		}

	}
}

func checkSignature(signature string, timestamp string, nonce string) bool {
	//第一步：将token、timestamp、nonce三个参数进行字典序排序
	tmpArr := []string{TOKEN, timestamp, nonce}
	sort(tmpArr)
	//第二步：将三个参数字符串拼接成一个字符串进行sha1加密
	tmpStr := implode(tmpArr)
	tmpStr = doSha1(tmpStr)
	if signature == tmpStr {
		return true
	}
	return false
}

//字典降序排列
func sort(arr []string) {
	length := len(arr)
	if length < 2 {
		return
	}
	for i := 0; i < length-1; i++ {
		tmp := ""
		for j := i + 1; j < length; j++ {
			if arr[i] > arr[j] {
				tmp = arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
}

//将数组组合为字符串
func implode(arr []string) string {
	ret := ""
	for i := 0; i < len(arr); i++ {
		ret += arr[i]
	}
	return ret
}

//对字符串进行sha1操作
func doSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func makeVoiceResponseBody(fromUserName, toUserName, mediaId string) ([]byte, error) {
	voiceResponseBody := &VoiceResponseBody{}
	voiceResponseBody.FromUserName = fromUserName
	voiceResponseBody.ToUserName = toUserName
	voiceResponseBody.MsgType = "voice"
	voiceResponseBody.MediaId = mediaId
	voiceResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(voiceResponseBody, " ", "  ")
}

func getMsgType(c MainController) string {
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if nil != err {
		return ""
	}
	c.Ctx.WriteString("测试测试")
	//	requestBody := MsgType{}
	//	xml.Unmarshal(body, &requestBody)
	return fmt.Sprintln(body)
}
