//以Json格式保存配置数据
package jsonconfig

import (
	//	"bytes"
	//"encoding/gob"
	"encoding/json"
	"github.com/seefan/gosee/utils"
	"io/ioutil"
	"os"
	//	"reflect"
)

//配置主体
type JsonConfig struct {
	//主体数据
	data interface{}
}

//新建的个Config配置
func New() *JsonConfig {
	return &JsonConfig{
		data: make(map[string]interface{}),
	}
}

//新建的个Config配置，并加载指定文件
func NewAndLoad(path string) (*JsonConfig, error) {
	jc := New()
	return jc, jc.Read(path)
}

//从一个指定文件读取配置
func (this *JsonConfig) Read(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if this.data == nil {
		this.data = make(map[string]interface{})
	}
	if err := json.Unmarshal(content, &this.data); err != nil {
		return err
	}
	return nil
}

//保存文件到指定路径
func (this *JsonConfig) Write(path string) error {
	bt, err := json.Marshal(this.data)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path, bt, os.ModePerm); err != nil {
		return err
	}
	return nil
}

//保存文件到指定路径，并对文件进行格式化，方便阅读
func (this *JsonConfig) WriteIndent(path string) error {
	bt, err := json.MarshalIndent(this.data, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path, bt, os.ModePerm); err != nil {
		return err
	}
	return nil
}

//将原数据按map取出
func (this *JsonConfig) Map() (map[string]interface{}, bool) {
	if this.data != nil {
		m, ok := (this.data).(map[string]interface{})
		return m, ok
	} else {
		return nil, false
	}
}

//取指定的取节点
func (this *JsonConfig) Get(key string) *JsonConfig {
	if m, ok := this.Map(); ok {
		if val, ok := m[key]; ok {
			return &JsonConfig{val}
		}
	}
	return &JsonConfig{nil}
}

//取指定的取节点
func (this *JsonConfig) GetArray(key string) []*JsonConfig {
	jss := []*JsonConfig{}
	arrs := this.Get(key)
	if ms, ok := arrs.data.([]interface{}); ok {
		for _, d := range ms {
			jss = append(jss, &JsonConfig{d})
		}
	}
	return jss
}

//设置一个节点的值，可以是简单类型，也可以是字典或其它类型
func (this *JsonConfig) Set(key string, v interface{}) *JsonConfig {
	if m, ok := this.Map(); ok {
		m[key] = v
	}
	return this
}

//转为指定的interface，注意vp必须是特定类型的指针
func (this *JsonConfig) Interface(vp interface{}) error {
	b, err := json.Marshal(this.data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, vp)
	if err != nil {
		return err
	}
	return nil
}

//取指定节点的值，并指定类型
func (this *JsonConfig) Int(key string, defaultValue ...int) int {
	if v := this.Get(key); v.data != nil {
		switch v.data.(type) {
		case json.Number:
			if i, err := v.data.(json.Number).Int64(); err == nil {
				return int(i)
			}
		default:
			return utils.AsInt(v.data, defaultValue...)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

//取指定节点的值，并指定类型
func (this *JsonConfig) Int64(key string, defaultValue ...int64) int64 {
	if v := this.Get(key); v.data != nil {
		switch v.data.(type) {
		case json.Number:
			if i, err := v.data.(json.Number).Int64(); err == nil {
				return int64(i)
			}
		default:
			return utils.AsInt64(v.data, defaultValue...)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

//取指定节点的值，并指定类型
func (this *JsonConfig) Float64(key string, defaultValue ...float64) float64 {
	if v := this.Get(key); v.data != nil {
		switch v.data.(type) {
		case json.Number:
			if i, err := v.data.(json.Number).Float64(); err == nil {
				return float64(i)
			}
		default:
			return utils.AsFloat64(v.data, defaultValue...)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

//取指定节点的值，并指定类型
func (this *JsonConfig) String(key string, defaultValue ...string) string {
	if v := this.Get(key); v.data != nil {
		if s := utils.AsString(v.data); len(s) > 0 {
			return s
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

//取指定节点的值，并指定类型
func (this *JsonConfig) Array(key string) []interface{} {
	if v := this.Get(key); v.data != nil {
		if re, ok := v.data.([]interface{}); ok {
			return re
		}
	}
	return []interface{}{}
}
