package services

import (
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
)

// ExcelReader 结构体用于读取 Excel 文件
type ExcelReader struct {
	inputStream io.Reader
}

// NewExcelReader 创建一个新的 ExcelReader 实例
func NewExcelReader(inputStream io.Reader) *ExcelReader {
	return &ExcelReader{
		inputStream: inputStream,
	}
}

// ReadExcel 读取 Excel 文件并返回数据
func (er *ExcelReader) ReadExcel() ([]map[string]interface{}, error) {
	f, err := excelize.OpenReader(er.inputStream)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取第一个工作表的名称
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	// 从第二行开始读取数据
	for _, row := range rows[1:] {
		r := make(map[string]interface{})
		if len(row) > 1 {
			r["topic"] = row[0]
		}
		if len(row) > 2 {
			r["topic_material_id"] = row[1]
		}
		if len(row) > 3 {
			r["answer"] = row[2]
		}
		if len(row) > 4 {
			r["topic_type"] = row[3]
		}
		if len(row) > 5 {
			score, err := strconv.ParseFloat(row[4], 64)
			if err == nil {
				r["score"] = score
			}
		}
		if len(row) > 6 {
			difficulty, err := strconv.ParseInt(row[5], 10, 64)
			if err == nil {
				r["difficulty"] = difficulty
			}
		}
		if len(row) > 7 {
			r["chapter_1"] = row[6]
		}
		if len(row) > 8 {
			r["chapter_2"] = row[7]
		}
		if len(row) > 9 {
			r["label_1"] = row[8]
		}
		if len(row) >= 10 {
			r["label_2"] = row[9]
			r["update_time"] = row[10]
		}
		result = append(result, r)
	}
	return result, nil
}

//func main() {
//	// 示例文件路径，需要替换为实际路径
//	filePath := "/Users/xiaobocai/Downloads/QuestionBank.xlsx"
//	f, err := excelize.OpenFile(filePath)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer f.Close()
//
//	// 获取第一个工作表的名称
//	sheetName := f.GetSheetName(0)
//	rows, err := f.GetRows(sheetName)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	var result []map[string]interface{}
//	// 从第二行开始读取数据
//	for _, row := range rows[1:] {
//		r := make(map[string]interface{})
//		if len(row) > 1 {
//			r["topic"] = row[1]
//		}
//		if len(row) > 2 {
//			r["topic_material_id"] = row[2]
//		}
//		if len(row) > 3 {
//			r["answer"] = row[3]
//		}
//		if len(row) > 4 {
//			r["topic_type"] = row[4]
//		}
//		if len(row) > 5 {
//			score, err := strconv.ParseFloat(row[5], 64)
//			if err == nil {
//				r["score"] = score
//			}
//		}
//		if len(row) > 6 {
//			difficulty, err := strconv.ParseFloat(row[6], 64)
//			if err == nil {
//				r["difficulty"] = difficulty
//			}
//		}
//		if len(row) > 7 {
//			r["chapter_1"] = row[7]
//		}
//		if len(row) > 8 {
//			r["chapter_2"] = row[8]
//		}
//		if len(row) > 9 {
//			r["label_1"] = row[9]
//		}
//		if len(row) > 10 {
//			r["label_2"] = row[10]
//		}
//		result = append(result, r)
//	}
//	fmt.Println(result)
//}
