package core_test

import (
	"log"
	"testing"

	"github.com/godcong/wego/core"
)

func TestLocalAddress(t *testing.T) {
	core.GetServerIp()
}

func TestXmlToMap(t *testing.T) {
	m := core.XmlToMap([]byte(`<xml>
<return_code><![CDATA[FAIL]]></return_code>
<return_msg><![CDATA[CERT_ERR]]></return_msg>
</xml>`))
	log.Println(m)
}

func TestSignatureSHA1(t *testing.T) {
	// s := "jsapi_ticket=sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg&noncestr=Wm3WZYTPz0wzccnW×tamp=1414587457&url=http://mp.weixin.qq.com?params=value"
	p := core.Map{
		"noncestr":     "Wm3WZYTPz0wzccnW",
		"jsapi_ticket": "9KwiourQPRN3vx3Nn1c_iX9qGaI3Cf8dwVy7qqYeYKcd3BK4Zd_jSlol7E7baUfgOY0E2ybaw2OrlhkChKaS7w",
		"timestamp":    1414587457,
		"url":          "http://mp.weixin.qq.com?params=value",
	}
	s := core.SignatureSHA1(core.Map{
		"noncestr":     "Wm3WZYTPz0wzccnW",
		"jsapi_ticket": "9KwiourQPRN3vx3Nn1c_iX9qGaI3Cf8dwVy7qqYeYKcd3BK4Zd_jSlol7E7baUfgOY0E2ybaw2OrlhkChKaS7w",
		"timestamp":    1414587457,
		"url":          "http://mp.weixin.qq.com?params=value",
	})
	if s != "32eb8ad0c84b65a8c5c73674e15f47ebdee48b13" {
		t.Error("SignatureSHA1", s)
	}

	s1 := core.SHA1("jsapi_ticket=9KwiourQPRN3vx3Nn1c_iX9qGaI3Cf8dwVy7qqYeYKcd3BK4Zd_jSlol7E7baUfgOY0E2ybaw2OrlhkChKaS7w&noncestr=Wm3WZYTPz0wzccnW&timestamp=1414587457&url=http://mp.weixin.qq.com?params=value")
	if s1 != "32eb8ad0c84b65a8c5c73674e15f47ebdee48b13" {
		t.Error("SHA1", s1)
	}
	// 32eb8ad0c84b65a8c5c73674e15f47ebdee48b13
}
