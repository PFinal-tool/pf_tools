package pak

import (
	"fmt"
	"os"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/8/6
 * @Desc:
 * @Project: pf_tools
 */

type Mkdir struct {
	Path       string
	Recursion  bool
	Seperation string
}

func (m *Mkdir) Done() {
	// 获取当前时间
	now := time.Now()
	// 格式化时间，获取年、月、日
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	// 构造文件夹路径
	dirPath := fmt.Sprintf("%s"+m.Seperation+"%s"+m.Seperation+"%s", year, month, day)
	if m.Recursion {
		dirPath = fmt.Sprintf("%s/%s/%s", year, month, day)
	}

	// 创建多级文件夹
	err := os.MkdirAll(m.Path+"/"+dirPath, 0755)
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
		return
	}

	fmt.Println("创建文件夹成功：", m.Path+"/"+dirPath)
}
