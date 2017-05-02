package jpush

import (
	"bytes"
	"fmt"
)

// AudienceObject 推送目标对象
type AudienceObject struct {
	op    string
	value []string
}

// Audience 推送目标
type Audience struct {
	isall   bool
	objects []AudienceObject
}

// MarshalJSON 自定义JSON序列化内容
func (a Audience) MarshalJSON() ([]byte, error) {
	if a.isall || len(a.objects) <= 0 {
		return []byte(`"all"`), nil
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	for _, obj := range a.objects {
		fmt.Fprintf(&buf, `"%s":[`, obj.op)
		for _, v := range obj.value {
			fmt.Fprintf(&buf, `"%s",`, v)
		}
		buf.Truncate(buf.Len() - 1)
		buf.WriteString("],")
	}
	buf.Truncate(buf.Len() - 1)
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// AddObject 增加推送目标
func (a *Audience) AddObject(obj AudienceObject) *Audience {
	a.objects = append(a.objects, obj)
	return a
}

// AudienceAll 全平台推送
func AudienceAll() Audience {
	return Audience{isall: true}
}

// AudienceTag 多个标签之间是 OR 的关系，即取并集
func AudienceTag(tag ...string) AudienceObject {
	return AudienceObject{
		op:    "tag",
		value: tag,
	}
}

// AudienceTagAnd 多个标签之间是 AND 关系，即取交集
func AudienceTagAnd(tag ...string) AudienceObject {
	return AudienceObject{
		op:    "tag_and",
		value: tag,
	}
}

// AudienceAlias 多个别名之间是 OR 关系，即取并集
func AudienceAlias(alias ...string) AudienceObject {
	return AudienceObject{
		op:    "alias",
		value: alias,
	}
}

// AudienceRegID 多个注册ID之间是 OR 关系，即取并集
func AudienceRegID(ids ...string) AudienceObject {
	return AudienceObject{
		op:    "registration_id",
		value: ids,
	}
}
