// helper 自用常用函数工具包

package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// IsIPv4 判断是否是否 ipv4 地址
func IsIPv4(address string) bool {
	res, _ := regexp.MatchString(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`, address)
	return res
}

// IsIPv6 判断是否是 ipv6 地址
func IsIPv6(address string) bool {
	res, _ := regexp.MatchString(`^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$`, address)
	return res
}

// IsEmail 判断是否是一个邮箱号
func IsEmail(email string) bool {
	res, _ := regexp.MatchString(`^[0-9A-Za-zd]+([-_.][0-9A-Za-zd]+)*@([0-9A-Za-zd]+[-.])+[A-Za-zd]{2,5}$`,
		email)
	return res
}

// IsURL 判断是否是一个合法的 URL
func IsURL(url string) bool {
	result, _ := regexp.MatchString(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, url)
	return result
}

// IsMobile 是否是手机号码
func IsMobile(mobile string) bool {
	result, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobile)
	return result
}

// HaveChinese 判断字符串中有无中文
func HaveChinese(str string) bool {
	var a = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	//接受正则表达式的范围
	for _, v := range str {
		if a.MatchString(string(v)) {
			return true
		}
	}
	return false
}

// HaveLetter 判断字符串中是否有英文字母
func HaveLetter(str string) bool {
	for _, s := range str {
		if unicode.IsLetter(s) {
			return true
		}
	}
	return false
}

// IsOnlyNumber 是否是纯数字
func IsOnlyNumber(str string) bool {
	res, _ := regexp.MatchString(`^[0-9]*$`,
		str)
	return res
}

// HaveSpecialCharacters 是否有特殊字符
func HaveSpecialCharacters(str string) bool {
	res, _ := regexp.MatchString("[`~!@#$^&*()=|{}':;',\\[\\].<>《》/?~！@#￥……&*（）——|{}【】‘；：”“'。，、？ ]",
		str)
	return res
}

// 压缩字符串,将字符串去除各 tab 、空格等
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	return strings.Join(strings.Fields(str), "")
}

// ParseDate 解析一条日期 兼容多种格式
func ParseDate(str string) (time.Time, error) {
	var (
		date  time.Time
		err   error
		float float64
	)
	switch {
	case strings.Contains(str, "/"):
		date, err = time.ParseInLocation("2006/01/02", str, time.Local)
	case strings.Contains(str, "-"):
		date, err = time.ParseInLocation("2006-01-02", str, time.Local)
	default:
		// 防止Excel 出现浮点型数字如20200101.0必须要转换为浮点型先
		float, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return date, fmt.Errorf("日期解析失败 %s", str)
		}
		strTime := strconv.Itoa(int(float))
		date, err = time.ParseInLocation("20060102", strTime, time.Local)
	}
	if err != nil {
		return date, fmt.Errorf("日期解析失败 %s", str)
	}
	return date, nil //
}

// 生成唯一SessionID
func SessionID() string {
	unix := Time() //获取当前时间戳
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	res := unix + rndNum
	return base64.URLEncoding.EncodeToString([]byte(Md5(strconv.FormatInt(res, 10))))
}

// Time 获取时间戳
func Time() int64 {
	return time.Now().Unix()
}

// TimeStampToDate 时间戳 转日期	2020-02-06 19:36:32
func TimeStampToDate(timestamp int) string {
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02 15:04:05") // 2018-07-11 15:10:19
}

// PHP date函数,timeStamp -1 获取当前标准日期 Y-m-d H:i:s
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

// Date 方法的升级版,加入错误返回
// PHP date函数,timeStamp -1 获取当前标准日期 Y-m-d H:i:s
func MustDate(str string, timeStamp int) (string, error) {
	var tm time.Time
	if timeStamp != -1 {
		tm = time.Unix(int64(timeStamp), 0)
	} else {
		tm = time.Unix(Time(), 0)
	}
	switch str {
	case "Y-m-d H:i:s":
		return tm.Format("2006-01-02 15:04:05"), nil //2018-07-11 15:10:19
	case "Y-m-d h:i:s":
		return tm.Format("2006-01-02 03:04:05 PM"), nil //2020-02-06 11:15:53 PM
	case "Y-m-d H:i":
		return tm.Format("2006-01-02 15:04"), nil //2018-07-11 15:10:19
	case "Y-m-d h:i":
		return tm.Format("2006-01-02 03:04 PM"), nil
	case "Y-m-d H":
		return tm.Format("2006-01-02 15"), nil
	case "Y-m-d h":
		return tm.Format("2006-01-02 03 PM"), nil
	default:
		return "", fmt.Errorf("找不到匹配的格式")
	}
}

// Str2DEC 二进制字符串转十进制
func Str2DEC(s string) (num int) {
	l := len(s)
	for i := l - 1; i >= 0; i-- {
		num += (int(s[l-i-1]) - 48) << uint8(i)
	}
	return
}

// GetMax 两个值之间取最大值
func GetMax(i, j int) int {
	if i > j {
		return i
	}
	return j
}

// slice去重,传入slice,返回一个新的值
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

// slice去重,传入slice指针,直接修改原数据
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

// md5
func Md5(str string) string {
	Md5Inst := md5.New()
	_, _ = Md5Inst.Write([]byte(str))
	Result := Md5Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", Result)
}

// sha1
func Sha1(str string) string {
	Sha1Inst := sha1.New()
	_, _ = Sha1Inst.Write([]byte(str))
	Result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x\n\n", Result)
}

// 字符串转换为hash数字
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

// InArray 是否在数组内
func InArray(val, array interface{}) bool {
	if reflect.TypeOf(array).Kind() == reflect.Slice {
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
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

// ================ 身份证相关 =================

// GetAgeByID 凭身份证计算年龄
func GetAgeByID(identification_number string) int {
	reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)(\d{2})((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	//reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)`)
	params := reg.FindStringSubmatch(identification_number)
	birYear, _ := strconv.Atoi(params[1] + params[2])
	birMonth, _ := strconv.Atoi(params[3])
	age := time.Now().Year() - birYear
	if int(time.Now().Month()) < birMonth {
		age--
	}
	return age
}

var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var validValue = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var validProvince = []string{
	"11", // 北京市
	"12", // 天津市
	"13", // 河北省
	"14", // 山西省
	"15", // 内蒙古自治区
	"21", // 辽宁省
	"22", // 吉林省
	"23", // 黑龙江省
	"31", // 上海市
	"32", // 江苏省
	"33", // 浙江省
	"34", // 安徽省
	"35", // 福建省
	"36", // 山西省
	"37", // 山东省
	"41", // 河南省
	"42", // 湖北省
	"43", // 湖南省
	"44", // 广东省
	"45", // 广西壮族自治区
	"46", // 海南省
	"50", // 重庆市
	"51", // 四川省
	"52", // 贵州省
	"53", // 云南省
	"54", // 西藏自治区
	"61", // 陕西省
	"62", // 甘肃省
	"63", // 青海省
	"64", // 宁夏回族自治区
	"65", // 新疆维吾尔自治区
	"71", // 台湾省
	"81", // 香港特别行政区
	"91", // 澳门特别行政区
}

// IsValidCitizenNo 校验是否是有效的身份证号码
func IsValidCitizenNo(citizenNo *[]byte) bool {
	if citizenNo == nil || *citizenNo == nil {
		return false
	}
	nLen := len(*citizenNo)
	if nLen != 15 && nLen != 18 {
		return false
	}

	for i, v := range *citizenNo {
		n, _ := strconv.Atoi(string(v))
		if n >= 0 && n <= 9 {
			continue
		}

		if v == 'X' && i == 16 {
			continue
		}

		return false
	}

	if !CheckProvinceValid(*citizenNo) {
		return false
	}

	if nLen == 15 {
		*citizenNo = Citizen15To18(*citizenNo)
		if *citizenNo == nil {
			return false
		}
	} else if !IsValidCitizenNo18(citizenNo) {
		return false
	}

	nYear, _ := strconv.Atoi(string((*citizenNo)[6:10]))
	nMonth, _ := strconv.Atoi(string((*citizenNo)[10:12]))
	nDay, _ := strconv.Atoi(string((*citizenNo)[12:14]))

	return CheckBirthdayValid(nYear, nMonth, nDay)
}

// GetCitizenNoInfo 从身份证号码中获得信息。生日(时间戳)，性别，省份代码。
func GetCitizenNoInfo(citizenNo []byte) (err error, birthday int64, isMale bool, addrMask int) {
	err = nil
	birthday = 0
	isMale = false
	addrMask = 0
	if !IsValidCitizenNo(&citizenNo) {
		err = errors.New("Invalid citizen number.")
		return
	}

	// Birthday information.
	nYear, _ := strconv.Atoi(string(citizenNo[6:10]))
	nMonth, _ := strconv.Atoi(string(citizenNo[10:12]))
	nDay, _ := strconv.Atoi(string(citizenNo[12:14]))
	birthday = time.Date(nYear, time.Month(nMonth), nDay, 0, 0, 0, 0, time.Local).Unix()

	// Gender information.
	genderMask, _ := strconv.Atoi(string(citizenNo[16]))
	if genderMask%2 == 0 {
		isMale = false
	} else {
		isMale = true
	}

	// Address code mask.
	addrMask, _ = strconv.Atoi(string(citizenNo[:2]))

	return
}

// IsValidCitizenNo18 校验是否是 18 位身份证
func IsValidCitizenNo18(citizenNo18 *[]byte) bool {
	nLen := len(*citizenNo18)
	if nLen != 18 {
		return false
	}

	nSum := 0
	for i := 0; i < nLen-1; i++ {
		n, _ := strconv.Atoi(string((*citizenNo18)[i]))
		nSum += n * weight[i]
	}
	mod := nSum % 11
	return validValue[mod] == (*citizenNo18)[17]
}

// Citizen15To18 15 位身份证转 18 位
func Citizen15To18(citizenNo15 []byte) []byte {
	nLen := len(citizenNo15)
	if nLen != 15 {
		return nil
	}

	citizenNo18 := make([]byte, 0)
	citizenNo18 = append(citizenNo18, citizenNo15[:6]...)
	citizenNo18 = append(citizenNo18, '1', '9')
	citizenNo18 = append(citizenNo18, citizenNo15[6:]...)

	sum := 0
	for i, v := range citizenNo18 {
		n, _ := strconv.Atoi(string(v))
		sum += n * weight[i]
	}
	mod := sum % 11
	citizenNo18 = append(citizenNo18, validValue[mod])

	return citizenNo18
}

// IsLeapYear 是否是闰年
func IsLeapYear(nYear int) bool {
	if nYear <= 0 {
		return false
	}

	if (nYear%4 == 0 && nYear%100 != 0) || nYear%400 == 0 {
		return true
	}

	return false
}

// CheckBirthdayValid 验证是否是一个有效的生日
func CheckBirthdayValid(nYear, nMonth, nDay int) bool {
	if nYear < 1900 || nMonth <= 0 || nMonth > 12 || nDay <= 0 || nDay > 31 {
		return false
	}

	curYear, curMonth, curDay := time.Now().Date()
	if nYear == curYear {
		if nMonth > int(curMonth) {
			return false
		} else if nMonth == int(curMonth) && nDay > curDay {
			return false
		}
	}

	if 2 == nMonth {
		if IsLeapYear(nYear) && nDay > 29 {
			return false
		} else if nDay > 28 {
			return false
		}
	} else if 4 == nMonth || 6 == nMonth || 9 == nMonth || 11 == nMonth {
		if nDay > 30 {
			return false
		}
	}

	return true
}

// CheckProvinceValid 校验身份证中的省份代码
func CheckProvinceValid(citizenNo []byte) bool {
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, citizenNo[:2]...)
	provinceStr := string(provinceCode)

	for i := range validProvince {
		if provinceStr == validProvince[i] {
			return true
		}
	}

	return false
}
