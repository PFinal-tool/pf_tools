package pak

import (
	"encoding/base64"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/14
 * @Desc:
 * @Project: pf_tools
 */

type Cryption struct {
	Quiet   bool
	URLSafe bool
}

func splitStringByLength(s string, length int) []string {
	var result []string
	for i := 0; i < len(s); i += length {
		if i+length > len(s) {
			result = append(result, s[i:])
		} else {
			result = append(result, s[i:i+length])
		}
	}
	return result
}

func (c *Cryption) encoding() *base64.Encoding {
	if c.URLSafe {
		return base64.URLEncoding
	}
	return base64.StdEncoding
}

func (c *Cryption) EncodeToString(str string) {
	sEnc := c.encoding().EncodeToString([]byte(str))
	view := ResultView{Quiet: c.Quiet}
	title := "加密结果"
	if c.URLSafe {
		title = "Base64URL 加密结果"
	}
	view.ShowText(title, sEnc, sEnc, true)
}

func (c *Cryption) DecodeString(sEnc string) {
	view := ResultView{Quiet: c.Quiet}
	sDec, err := c.encoding().DecodeString(sEnc)
	if err != nil {
		// 兼容：URLSafe 输入有时混用标准 base64
		if c.URLSafe {
			sDec, err = base64.StdEncoding.DecodeString(sEnc)
		}
		if err != nil {
			view.ShowText("解密结果", "解密失败：不是一个正确的base64编码", "", false)
			return
		}
	}
	text := string(sDec)
	view.ShowText("解密结果", "解密成功："+text, text, true)
}
