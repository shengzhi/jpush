// 用于在服务器端查询、设置、更新、删除设备的 tag,alias 信息

package jpush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TagAndAliasView struct {
	Tags   []string `json:"tags"`
	Alias  string   `json:"alias"`
	Mobile string   `json:"mobile"`
}

// ViewTagAndAlias 查看指定注册ID的标签和别名
func (j *JPushClient) ViewTagAndAlias(regid string) TagAndAliasView {
	uri := fmt.Sprintf("%s/v3/devices/%s", apigateway, regid)
	req, _ := http.NewRequest("GET", uri, nil)
	var reply TagAndAliasView
	j.do(req, &reply)
	return reply
}

type SetTagAndAliasParam struct {
	RegID string `json:"-"`
	Tags  struct {
		Add    []string `json:"add,omitempty"`
		Remove []string `json:"remove,omitempty"`
	} `json:"tags"`
	Alias  string `json:"alias,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}

// SetTagAndAlias 设置标签和别名
func (j *JPushClient) SetTagAndAlias(r SetTagAndAliasParam) error {
	uri := fmt.Sprintf("%s/v3/devices/%s", apigateway, r.RegID)
	data, _ := json.Marshal(r)
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	return j.do(req, nil)
}

func (j *JPushClient) do(httpReq *http.Request, reply interface{}) error {
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.basicauth))
	res, err := j.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if reply != nil {
		return json.NewDecoder(res.Body).Decode(reply)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("%d-%s", res.StatusCode, res.Status)
}
