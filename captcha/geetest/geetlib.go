package geetest

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"math"
)

const (
	verName     = "3.3.0"                  // SDK版本编号
	sdkLang     = "golang"                 // SD的语言类型
	apiUrl      = "http://api.geetest.com" //极验验证API URL
	registerUrl = "/register.php"          //register url
	validateUrl = "/validate.php"          //validate url

	//极验验证二次验证表单数据
	FN_GEETEST_CHALLENGE = "geetest_challenge"
	FN_GEETEST_VALIDATE  = "geetest_validate"
	FN_GEETEST_SECCODE   = "geetest_seccode"

	//极验验证API服务状态Session Key
	GtServerStatusSessionKey = "gt_server_status"
)

type GeetLib struct {
	captchaId   string //公钥
	privateKey  string //私钥
	userId      string
	responseStr interface{}
	debugCode   bool //调试开关，是否输出调试日志
}

func NewGeetLib(captchaId string, privateKey string, debug ... bool) *GeetLib {
	flag := false
	if len(debug)>0 {
		flag = debug[0]
	}
	return &GeetLib{captchaId: captchaId, privateKey: privateKey, debugCode: flag}
}

//获取本次验证初始化返回字符串
func (g *GeetLib) GetResponseStr() interface{} {

	return g.responseStr
}

func (g *GeetLib) GetVersionInfo() string {
	return verName
}

//预处理失败后的返回格式串
func (g *GeetLib) getFailPreProcessRes() interface{} {

	rnd1 := rand.Intn(100)
	rnd2 := rand.Intn(100)

	md5Str1 := g.md5Encode(strconv.Itoa(rnd1))
	md5Str2 := g.md5Encode(strconv.Itoa(rnd2))
	challenge := md5Str1 + md5Str2[0:2]
	result := map[string]interface{}{
		"success":   0,
		"gt":        g.captchaId,
		"challenge": challenge}
	return result
}

//预处理成功后的标准串
func (g *GeetLib) getSuccessPreProcessRes(challenge string) interface{} {

	g.gtlog("challenge:" + challenge)
	result := map[string]interface{}{
		"success":   1,
		"gt":        g.captchaId,
		"challenge": challenge}
	return result
}

// 验证初始化预处理 1表示初始化成功，0表示初始化失败

func (g *GeetLib) PreProcess(userid string) int {

	g.userId = userid;
	return g.preProcess()
}

//验证初始化预处理 1表示初始化成功，0表示初始化失败
func (g *GeetLib) preProcess() int {

	if g.registerChallenge() != 1 {

		g.responseStr = g.getFailPreProcessRes()
		return 0
	}
	return 1
}

// 用captchaID进行注册，更新challenge 1表示注册成功，0表示注册失败

func (g *GeetLib) registerChallenge() int {

	GET_URL := apiUrl + registerUrl + "?gt=" + g.captchaId
	if len(g.userId) != 0 {
		GET_URL = GET_URL + "&user_id=" + g.userId
		g.userId = ""
	}
	g.gtlog("GET_URL:" + GET_URL)
	result_str, err := g.readContentFromGet(GET_URL)
	if err != nil {
		g.gtlog("exception:register api")
		return 0
	}
	g.gtlog("register_result:" + result_str)
	if 32 == len(result_str) {
		g.responseStr = g.getSuccessPreProcessRes(g.md5Encode(result_str + g.privateKey))
		return 1
	} else {
		g.gtlog("gtServer register challenge failed")
		return 0
	}

	return 0
}

//发送请求，获取服务器返回结果
func (g *GeetLib) readContentFromGet(getURL string) (string, error) {

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	response, err := client.Get(getURL)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}

	body, _eer := ioutil.ReadAll(response.Body)
	if _eer != nil {
		return "", _eer
	}
	return string(body), nil
}

//判断一个表单对象值是否为空
func (g *GeetLib) objIsEmpty(gtObj interface{}) bool {
	if gtObj == nil {
		return true
	}
	if gtstr, ok := gtObj.(string); ok {
		if len(gtstr) == 0 {
			return true
		}
	}
	return false
}

//检查客户端的请求是否合法,三个只要有一个为空，则判断不合法
func (g *GeetLib) resquestIsLegal(challenge string, validate string, seccode string) bool {

	if g.objIsEmpty(challenge) || g.objIsEmpty(validate) || g.objIsEmpty(seccode) {
		return false
	}
	return true
}

func (g *GeetLib) EnhencedValidateRequest(challenge string, validate string, seccode string, userid string) int {
	g.userId = userid
	return g.enhencedValidateRequest(challenge, validate, seccode)
}

// 服务正常的情况下使用的验证方式,向gt-server进行二次验证,获取验证结果
func (g *GeetLib) enhencedValidateRequest(challenge string, validate string, seccode string) int {

	if !g.resquestIsLegal(challenge, validate, seccode) {
		return 0
	}
	g.gtlog("request legitimate")

	url := apiUrl + validateUrl
	query := fmt.Sprintf("seccode=%s&sdk=%s", seccode, sdkLang+"_"+verName)

	if len(g.userId) != 0 {
		query = query + "&user_id=" + g.userId
		g.userId = ""
	}
	g.gtlog(query)

	if len(validate) <= 0 {
		return 0
	}

	if !g.checkResultByPrivate(challenge, validate) {
		return 0
	}
	g.gtlog("checkResultByPrivate")

	response, err := g.postValidate(url, query)
	if err != nil {
		g.gtlog("response: " + response)
		return 0
	}
	g.gtlog("response: " + response)
	g.gtlog("md5: " + g.md5Encode(seccode))
	if strings.Compare(response, g.md5Encode(seccode)) == 0 {
		return 1
	}
	return 0

}

// failback使用的验证方式 1表示验证成功0表示验证失败

func (g *GeetLib) FailbackValidateRequest(challenge string, validate string, seccode string) int {

	g.gtlog("in failback validate")

	if !g.resquestIsLegal(challenge, validate, seccode) {
		return 0
	}
	g.gtlog("request legitimate")
	validateStr := strings.Split(validate, "_")
	if len(validateStr) < 3 {
		return 0
	}
	encodeAns := validateStr[0]
	encodeFullBgImgIndex := validateStr[1]
	encodeImgGrpIndex := validateStr[2]
	g.gtlog(fmt.Sprintf(
		"encode----challenge:%s--ans:%s,bg_idx:%s,grp_idx:%s",
		challenge, encodeAns, encodeFullBgImgIndex, encodeImgGrpIndex))
	decodeAns := g.decodeResponse(challenge, encodeAns)
	decodeFullBgImgIndex := g.decodeResponse(challenge, encodeFullBgImgIndex)
	decodeImgGrpIndex := g.decodeResponse(challenge, encodeImgGrpIndex)
	g.gtlog(fmt.Sprintf("decode----ans:%s,bg_idx:%s,grp_idx:%s", decodeAns,
		decodeFullBgImgIndex, decodeImgGrpIndex))
	validateResult := g.validateFailImage(decodeAns, decodeFullBgImgIndex, decodeImgGrpIndex)
	return validateResult
}

func (g *GeetLib) validateFailImage(ans int, full_bg_index int,
	img_grp_index int) int {

	full_bg_name := g.md5Encode(strconv.Itoa(full_bg_index))[0:9]
	bg_name := g.md5Encode(strconv.Itoa(img_grp_index))[10:19]
	answer_decode := ""

	// 通过两个字符串奇数和偶数位拼接产生答案位
	for i := 0; i < 9; i++ {
		if i%2 == 0 {
			answer_decode = answer_decode + string(full_bg_name[i])
		} else if i%2 == 1 {
			answer_decode = answer_decode + string(bg_name[i])
		} else {
			g.gtlog("exception")
		}
	}

	x_decode := answer_decode[4: ]
	x_int, _ := strconv.ParseInt(x_decode, 16, 32); // 16 to 10
	result := int(x_int) % 200
	if result < 40 {
		result = 40
	}
	if math.Abs(float64(ans-result)) <= 3.0 { // 3.0为容差值
		return 1
	} else {
		return 0
	}
}

// 解码随机参数

func (g *GeetLib) decodeResponse(challenge string, str string) int {
	if len(str) > 100 {
		return 0
	}

	shuzi := []int{1, 2, 5, 10, 50 }
	chongfu := "";
	keys := make(map[string]int, 0)

	count := 0;
	for i := 0; i < len(challenge); i++ {
		item := string(challenge[i])
		if strings.Contains(chongfu, item) {
			continue
		} else {
			value := shuzi[count%5]
			chongfu = chongfu + item
			count++
			keys[item] = value
		}
	}

	res := 0
	for j := 0; j < len(str); j++ {
		res = res + keys[string(str[j])]
	}
	res = res - g.decodeRandBase(challenge)
	return res
}

// 输入的两位的随机数字,解码出偏移量

func (g *GeetLib) decodeRandBase(challenge string) int {
	base := challenge[32:34]
	tempArray := make([]int, 0)
	for _, val := range base {
		var result int = 0
		if val > 57 {
			result = int(val) - 87
		} else {
			result = int(val) - 48
		}
		tempArray = append(tempArray, result)
	}

	return tempArray[0]*36 + tempArray[1];
}

func (g *GeetLib) postValidate(url string, data string) (string, error) {

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	g.gtlog(url + data)
	response, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(response.Body)

	return string(body), nil
}

func (g *GeetLib) checkResultByPrivate(challenge string, validate string) bool {
	encodeStr := g.md5Encode(g.privateKey + "geetest" + challenge)
	return strings.Compare(validate, encodeStr) == 0
}

//输出debug信息，需要开启debugCode
func (g *GeetLib) gtlog(message string) {
	if g.debugCode {
		fmt.Println("gtlog: " + message)
	}
}

//md5
func (g *GeetLib) md5Encode(text string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(text))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
