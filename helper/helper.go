package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// 工具包

//HTTP Post请求发送表单形式
//请求参数key_val:	url.Values{"name": {"小明"},"age": {"23"}}
func HttpPostForm(path string, val url.Values) (string, error) {
	resp, err := http.PostForm(path, val)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("读取失败")
	}
	return string(body), nil
}

//HTTP	get请求
func HttpGet(url string) (string, error) {
	var res, err = http.Get(url) //res是响应
	if err != nil {
		return "", errors.New("请求出错")
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("读取出错:", err)
		return "", errors.New("读取出错")
	}
	return string(bytes), nil
}

//生成唯一Session_Id
func Session_Id() string {
	unix := Time() //获取当前时间戳
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	res := unix + rndNum
	return base64.URLEncoding.EncodeToString([]byte(Md5(strconv.FormatInt(res, 10))))
}

//获取时间戳
func Time() int64 {
	return time.Now().Unix()
}

//时间戳 转日期	2020-02-06 19:36:32
func TimeStampToDate(timestamp int) string {
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02 15:04:05") //2018-07-11 15:10:19
}

//PHP date函数,timeStamp -1 获取当前标准日期 Y-m-d H:i:s
func Date(str string, timeStamp int) string {
	var tm time.Time
	if timeStamp != -1 {
		tm = time.Unix(int64(timeStamp), 0)
	} else {
		tm = time.Unix(Time(), 0)
	}
	switch str {
	case "Y-m-d H:i:s":
		return tm.Format("2006-01-02 15:04:05") //2018-07-11 15:10:19
	case "Y-m-d h:i:s":
		return tm.Format("2006-01-02 03:04:05 PM") //2020-02-06 11:15:53 PM
	case "Y-m-d H:i":
		return tm.Format("2006-01-02 15:04") //2018-07-11 15:10:19
	case "Y-m-d h:i":
		return tm.Format("2006-01-02 03:04 PM")
	case "Y-m-d H":
		return tm.Format("2006-01-02 15")
	case "Y-m-d h":
		return tm.Format("2006-01-02 03 PM")
	default:
		return "Date() error"
	}
}

//二进制字符串转十进制
func Str2DEC(s string) (num int) {
	l := len(s)
	for i := l - 1; i >= 0; i-- {
		num += (int(s[l-i-1]) - 48) << uint8(i)
	}
	return
}

//二叉树结构
type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Val   int
}

//判断两棵树是否相等,传入两棵树的根
func IsSameTree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p != nil && q != nil && p.Val == q.Val {
		return IsSameTree(p.Left, q.Left) && IsSameTree(p.Right, q.Right)
	} else {
		return false
	}
}

//获取树的最大深度
func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return GetMax(MaxDepth(root.Left), MaxDepth(root.Right)) + 1
}

//两个值之间取最大值
func GetMax(i, j int) int {
	if i > j {
		return i
	}
	return j
}

//slice去重,传入slice,返回一个新的值
func duplicate(data interface{}) interface{} {
	inArr := reflect.ValueOf(data)
	if inArr.Kind() != reflect.Slice && inArr.Kind() != reflect.Array {
		return data
	}
	existMap := make(map[interface{}]bool)
	outArr := reflect.MakeSlice(inArr.Type(), 0, inArr.Len())

	for i := 0; i < inArr.Len(); i++ {
		iVal := inArr.Index(i)

		if _, ok := existMap[iVal.Interface()]; !ok {
			outArr = reflect.Append(outArr, inArr.Index(i))
			existMap[iVal.Interface()] = true
		}
	}

	return outArr.Interface()
}

//slice去重,传入slice指针,直接修改原数据
func DuplicateSlice(data interface{}) {
	dataVal := reflect.ValueOf(data)
	if dataVal.Kind() != reflect.Ptr {
		fmt.Println("input data.kind is not pointer")
		return
	}

	tmpData := duplicate(dataVal.Elem().Interface())
	tmpDataVal := reflect.ValueOf(tmpData)

	dataVal.Elem().Set(tmpDataVal)
}

//生成随机数
func GetRandom(num int) int {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(num)
	return x
}

//md5加密
func Md5(str string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(str))
	Result := Md5Inst.Sum([]byte(""))
	return fmt.Sprintf("%x\n\n", Result)
}

//sha1加密
func Sha1(str string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(str))
	Result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x\n\n", Result)
}

//字符串转换为hash数字
func HashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return v
}

// 解析一个日期
func ParseDate(str string) (time.Time, error) {
	float, err := strconv.ParseFloat(str, 64)
	var date time.Time
	// 不能转代表填写的不是纯数字格式的日期
	if err != nil {
		date, err = time.ParseInLocation("2006/01/02", str, time.Local)
	} else {
		str := strconv.Itoa(int(float))
		date, err = time.ParseInLocation("20060102", str, time.Local)
	}
	// 如果都解析失败了,这条日期就真错了
	if err != nil {
		return date, fmt.Errorf("日期解析失败 %s", str)
	}
	return date, nil
}

// CleanStruct 清洗结构体中的每一个字段,当且结构体字段的值 与 old 全等时才会替换
// 当结构体字段值与 old 参数全等,则替换为 replace , omit 忽略某些结构体字段
func CleanStruct(structName interface{}, old, replace string, omit ...string) error {
	t := reflect.TypeOf(structName)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("首个参数必须为指针,且是一个结构体的指针")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("首个参数必须为指针,且是一个结构体的指针")
	}
	fieldNum := t.NumField()
ROW:
	for i := 0; i < fieldNum; i++ {
		name := t.Field(i).Name // 获取结构体字段名
		// 是否是需要忽略的字段
		if len(omit) != 0 {
			for _, item := range omit {
				if item == name {
					continue ROW
				}
			}
		}
		immutable := reflect.ValueOf(structName).Elem()
		// 获取该字段名的值
		val := immutable.FieldByName(name).String()
		if val == old {
			immutable.FieldByName(name).SetString(replace) // 替换
		}
	}
	return nil
}
