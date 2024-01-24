package pak

import (
	"fmt"
	"os"
	"path/filepath"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/1/24
 * @Desc:
 * @Project: pf_tools
 */

type ClearPath struct {
	Path string
}

func (c *ClearPath) removeAllFilesWithFilename(dirPath, filename string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			err := os.Remove(path)
			if err != nil {
				return err
			}
			fmt.Println("Deleted:", path)
		}
		return nil
	})
	return err
}

func (c *ClearPath) ClearDotDSStore() {
	filename := ".DS_Store"
	err := c.removeAllFilesWithFilename(c.Path, filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
