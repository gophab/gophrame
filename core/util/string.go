package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func StringAddr(s string) *string {
	if s == "" {
		return nil
	} else {
		return &s
	}
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func DefaultString(s string, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}

func NotNullString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func Nullable(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func ConditionString(condition bool, f1 interface{}, f2 interface{}) string {
	if condition {
		if f, ok := f1.(func() string); ok {
			return f()
		} else {
			return f1.(string)
		}
	} else {
		if f, ok := f2.(func() string); ok {
			return f()
		} else {
			return f2.(string)
		}
	}
}

func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		return ""
	}

	end := begin + length
	if end > lth {
		end = lth
	}

	return string(rs[begin:end])
}

// 字符首字母大写转换
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// 驼峰转C
func Camel2C(str string) string {
	var cstr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if vv[i] >= 65 && vv[i] <= 90 { // 后文有介绍
			vv[i] += 32 // string的码表相差32位
			if i > 0 {
				cstr += "_"
			}
			cstr += string(vv[i])
		} else {
			cstr += string(vv[i])
		}
	}
	return cstr
}

// C转驼峰
func C2Camel(str string) string {
	var cstr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if vv[i] == 45 || vv[i] == 95 {
			i++
			if i >= len(vv) {
				break
			}
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				if cstr != "" {
					vv[i] -= 32 // string的码表相差32位
				}
				cstr += string(vv[i])
			}
		} else {
			cstr += string(vv[i])
		}
	}
	return cstr
}

var variableNameReg = regexp.MustCompile(`\w+`)

func VariableName(str string) string {
	return variableNameReg.FindString(str)
}

func DbFieldName(s string) string {
	return Camel2C(VariableName(s))
}

func DbFields(kv map[string]interface{}) map[string]interface{} {
	var tkv = make(map[string]interface{})
	for k, v := range kv {
		tkv[Camel2C(VariableName(k))] = v
	}
	return tkv
}

func FormatParamterContent(content string, params map[string]string) string {
	reg, err := regexp.Compile("\\$\\{([\u4E00-\u9FA5A-Za-z0-9_]+.)*\\}")
	if err != nil {
		return content
	}
	return reg.ReplaceAllStringFunc(content, func(s string) string {
		// ${name:default}
		part, _ := strings.CutPrefix(s, "${")
		part, _ = strings.CutSuffix(part, "}")

		segs := strings.Split(part, ":")

		txt, b := params[segs[0]]
		if b {
			return txt
		}

		if len(segs) < 2 {
			return segs[0]
		} else {
			return segs[1]
		}
	})
}

func FormatParamterContentEx(content string, params map[string]interface{}) string {
	reg, err := regexp.Compile("\\$\\{([\u4E00-\u9FA5A-Za-z0-9_]+.)*\\}")
	if err != nil {
		return content
	}
	return reg.ReplaceAllStringFunc(content, func(s string) string {
		// ${name:default}
		part, _ := strings.CutPrefix(s, "${")
		part, _ = strings.CutSuffix(part, "}")

		segs := strings.Split(part, ":")

		txt := GetRecordFieldValue(params, segs[0], "")
		if txt != "" {
			return txt
		}

		if len(segs) < 2 {
			return segs[0]
		} else {
			return segs[1]
		}
	})
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子

	// 定义字符串字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func GenerateRandomNumeric(length int) string {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子

	// 定义字符串字符集
	charset := "0123456789"

	// 生成字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

// 定义script的正则表达式，去除js可以防止注入
var scriptRegex = regexp.MustCompile(`<script[^>]*?>[\s\S]*?<\\/script>`)

// 定义style的正则表达式，去除style样式，防止css代码过多时只截取到css样式代码
var styleRegex = regexp.MustCompile(`<style[^>]*?>[\s\S]*?<\/style>`)

// 定义HTML标签的正则表达式，去除标签，只提取文字内容
var htmlRegex = regexp.MustCompile(`<[^>]+>`)

// 定义空格,回车,换行符,制表符
var spaceRegex = regexp.MustCompile(`\s*|\t|\r|\n`)

func RemoveHTMLTags(htmlStr string) string {
	//定义字符串
	// htmlStr = "<script type>var i=1; alert(i)</script><style> .font1{font-size:12px}</style><span>少年中国说。</span>红日初升，其道大光。<h3>河出伏流，一泻汪洋。</h3>潜龙腾渊， 鳞爪&nbsp;&nbsp;飞扬。乳 虎啸  谷，百兽震惶。鹰隼试翼，风尘吸张。奇花初胎，矞矞皇皇。干将发硎，有作&nbsp其芒。天戴其苍，地履其黄。纵有千古，横有" +
	// 					"八荒。<a href=\"www.baidu.com\">前途似海，来日方长</a>。<h1>美哉我少年中国，与天不老！</h1><p>壮哉我中国少年，与国无疆！</p>";
	// 过滤script标签
	htmlStr = scriptRegex.ReplaceAllString(htmlStr, "")
	// 过滤style标签
	htmlStr = styleRegex.ReplaceAllString(htmlStr, "")
	// 过滤html标签
	htmlStr = htmlRegex.ReplaceAllString(htmlStr, "")
	// 过滤空格等
	htmlStr = spaceRegex.ReplaceAllString(htmlStr, "")
	// 过滤&nbsp;
	htmlStr = strings.ReplaceAll(htmlStr, "&nbsp;", "")
	// 过滤&nbsp
	htmlStr = strings.ReplaceAll(htmlStr, "&nbsp", "")
	// 返回文本字符串
	htmlStr = strings.Trim(htmlStr, " \t\r\n")

	return htmlStr
}
