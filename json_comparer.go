package json_comparer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type JsonComparer interface {
	CompareJson(example, resp string) (bool, error)
	IgnoreListSequence(bool) JsonComparer
}

type Comparer struct {
	ignore             []string
	logger             string
	ignoreListSequence bool
}

func NewComparer(ignore ...string) JsonComparer {
	return &Comparer{ignore: ignore}
}

func (c *Comparer) IgnoreListSequence(b bool) JsonComparer {
	c.ignoreListSequence = b
	return c
}

func (c *Comparer) errLog() error {
	if c.logger == "" {
		return nil
	}
	log := "\n"
	split := strings.Split(c.logger, "\n")
	for i := 0; i < len(split); i++ {
		for j := 0; j < i; j++ {
			log += " "
		}
		log += split[i] + "\n"
	}
	c.logger = ""
	return fmt.Errorf("%v", log)
}

func (c *Comparer) CompareJson(a, b string) (bool, error) {
	var m1, m2 map[string]interface{}
	if err := json.Unmarshal([]byte(a), &m1); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(b), &m2); err != nil {
		return false, err
	}
	return c.CompareMap(m1, m2), c.errLog()
}

func (c *Comparer) CompareValue(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		c.logger = fmt.Sprintf("%v != %v", a, b)
		return false
	}
	k1, k2 := reflect.TypeOf(a).Kind(), reflect.TypeOf(b).Kind()
	if k1 != k2 {
		c.logger = fmt.Sprintf("type: %s != %s", k1.String(), k2.String())
		return false
	}
	switch k1 {
	case reflect.Map:
		return c.CompareMap(a.(map[string]interface{}), b.(map[string]interface{}))
	case reflect.Slice:
		return c.CompareSlice(a, b)
	default:
		ok := reflect.DeepEqual(a, b)
		if !ok {
			c.logger = fmt.Sprintf("value: %s != %s", a, b)
		}
		return ok
	}
}

func (c *Comparer) CompareSliceIgnoreSequence(v1, v2 reflect.Value) bool {
	for i := 0; i < v1.Len(); i++ {
		j := 0
		for ; j < v2.Len(); j++ {
			if c.CompareValue(v1.Index(i).Interface(), v2.Index(j).Interface()) {
				break
			}
		}
		if j == v2.Len() {
			value, _ := json.Marshal(v1.Index(i).Interface())
			c.logger = fmt.Sprintf("index %d not found:\n%s", i, value)
			return false
		}
	}
	c.logger = ""
	return true
}

func (c *Comparer) CompareSlice(a, b interface{}) bool {
	v1, v2 := reflect.ValueOf(a), reflect.ValueOf(b)
	if v1.Len() != v2.Len() {
		c.logger = fmt.Sprintf("length: %d != %d", v1.Len(), v2.Len())
		return false
	}
	if c.ignoreListSequence {
		return c.CompareSliceIgnoreSequence(v1, v2)
	}
	for i := 0; i < v1.Len(); i++ {
		if !c.CompareValue(v1.Index(i).Interface(), v2.Index(i).Interface()) {
			c.logger = fmt.Sprintf("index %d:\n%s", i, c.logger)
			return false
		}
	}
	return true
}

func (c *Comparer) shouldIgnore(k string) bool {
	i := 0
	for ; i < len(c.ignore); i++ {
		if k == c.ignore[i] {
			return true
		}
	}
	return false
}

func (c *Comparer) CompareMap(a, b map[string]interface{}) bool {
	for k, va := range a {
		vb, ok := b[k]
		if !ok {
			c.logger = fmt.Sprintf("field: `%s` not found", k)
			return false
		}
		if c.shouldIgnore(k) {
			continue
		}
		ok = c.CompareValue(va, vb)
		if !ok {
			c.logger = fmt.Sprintf("field `%s`:\n%s", k, c.logger)
			return false
		}
	}
	flag := true
	for k := range b {
		_, ok := a[k]
		if !ok {
			c.logger = fmt.Sprintf("unexpected field: `%s`\n%s", k, c.logger)
			flag = false
		}
	}
	return flag
}

func (c *Comparer) DeepCopy(obj interface{}) interface{} {
	temp, _ := json.Marshal(obj)
	var res interface{}
	_ = json.Unmarshal(temp, &res)
	return res
}
