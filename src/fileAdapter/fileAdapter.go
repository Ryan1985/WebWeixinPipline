package fileAdapter

import (
	"fmt"
	"os"
)

func WriteFile(fileName string, fileBuffer []byte) (bool, error) {
	//打开文件，返回文件指针
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(file)
	defer file.Close()

	//写入byte的slice数据
	n, err := file.Write(fileBuffer)

	return n > 0, err
}
