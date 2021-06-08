package Utils

import (
	"encoding/csv"
	"os"
)

func WriteCsv(filename string, data [][]string) error {
	//创建文件
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 写入UTF-8 BOM
	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return err
	}
	//创建一个新的写入文件流
	w := csv.NewWriter(f)
	//写入数据
	err = w.WriteAll(data)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
