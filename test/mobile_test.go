package test

import (
	"fmt"
	"github.com/zheng-ji/gophone"
	"testing"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/10
 * @Desc:
 * @Project: pf_tools
 */

func TestGetinfo(t *testing.T) {
	pr, err := gophone.Find("18757753327")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pr)
}
