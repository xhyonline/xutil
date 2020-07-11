package xtype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Strings ===========字符串列表===========
type Strings []string

// String 转换为string类型
func (t Strings) String() string {
	if t == nil {
		return ""
	}
	var s = []string(t)
	return strings.Join(s, ",")
}

// Exists 元素是否存在
// TODO: 待废弃，改名 Contains
func (t Strings) Exists(s string) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// Contains 是否包含元素
func (t Strings) Contains(s string) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// Intersectant 是否有交集
func (t Strings) Intersectant(s Strings) bool {
	for _, ss := range s {
		if t.Contains(ss) {
			return true
		}
	}
	return false
}

// SAdd 添加不重复的元素
func (t Strings) SAdd(s string) Strings {
	if t.Contains(s) {
		return t
	}
	return append(t, s)
}

// Remove 移除所有值为 s 的元素
func (t Strings) Remove(s string) Strings {
	tmp := make(Strings, 0, len(t))
	for _, tt := range t {
		if tt != s {
			tmp = append(tmp, tt)
		}
	}
	return tmp
}

// Union 返回并集，t若有重复不会被去重
func (t Strings) Union(s Strings) Strings {
	for _, ss := range s {
		t = t.SAdd(ss)
	}
	return t
}

// Sub 返回集合 t 减掉 s 后的集合
func (t Strings) Sub(s Strings) Strings {
	for _, ss := range s {
		t = t.Remove(ss)
	}
	return t
}

// MarshalJSON 转换为json类型
func (t Strings) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(t))
}

// UnmarshalJSON 不做处理
func (t *Strings) UnmarshalJSON(data []byte) error {
	var tmp []string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = Strings(tmp)
	return nil
}

// MarshalYAML 转换为json类型
func (t Strings) MarshalYAML() ([]byte, error) {
	if t == nil {
		return []byte{}, nil
	}
	return yaml.Marshal([]string(t))
}

// UnmarshalYAML 不做处理
func (t *Strings) UnmarshalYAML(data []byte) error {
	var tmp []string
	if err := yaml.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = Strings(tmp)
	return nil
}

// Scan implements the Scanner interface.
func (t *Strings) Scan(src interface{}) error {
	*t = make([]string, 0)
	if src == nil {
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("Read tags from DB failed")
	}
	if len(tmp) == 0 {
		return nil
	}
	*t = strings.Split(string(tmp), ",")
	return nil
}

// Value implements the driver Valuer interface.
func (t Strings) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.String(), nil
}

// Numbers ===========数字列表===========
type Numbers []int

// String 转换为string类型
func (t Numbers) String() string {
	if t == nil {
		return ""
	}
	var s = []int(t)
	var ns []string
	for _, i := range s {
		ns = append(ns, strconv.Itoa(i))
	}
	return strings.Join(ns, ",")
}

// Exists 元素是否存在
// TODO: 待废弃 改为Contains
func (t Numbers) Exists(s int) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// Contains 是否包含元素
func (t Numbers) Contains(s int) bool {
	for _, item := range t {
		if item == s {
			return true
		}
	}
	return false
}

// MarshalJSON 转换为json类型 加域名
func (t Numbers) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]int(t))
}

// UnmarshalJSON 不做处理
func (t *Numbers) UnmarshalJSON(data []byte) error {
	var tmp []int
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*t = Numbers(tmp)
	return nil
}

// Scan implements the Scanner interface.
func (t *Numbers) Scan(src interface{}) error {
	*t = make([]int, 0)
	if src == nil {
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("Read tags from DB failed")
	}
	if len(tmp) == 0 {
		return nil
	}
	ts := strings.Split(string(tmp), ",")
	for _, i := range ts {
		n, _ := strconv.Atoi(i)
		*t = append(*t, n)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t Numbers) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.String(), nil
}
