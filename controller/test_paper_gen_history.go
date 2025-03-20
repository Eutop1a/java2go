package controller

import (
	"github.com/gin-gonic/gin"
	"java2go/entity"
	"java2go/mapper"
	"java2go/utils"
	"log"
	"time"
)

// TestPaperGenHistoryMapper 定义测试试卷生成历史映射器接口
type TestPaperGenHistoryMapper interface {
	InsertTestPaperGenHistory(testPaperGenHistory entity.TestPaperGenHistory) (int64, error)
	QueryAllTestPaperGenHistory() ([]entity.TestPaperGenHistory, error)
	GetTestPaperNameByTestPaperUid(testPaperUID string) ([]string, error)
	UpdateTestPaperTime(testPaperUID string, date time.Time) (int64, error)
	DeleteTestPaperGenHistoryByTestPaperUid(testPaperUID string) (int64, error)
}

var testMapper, _ = mapper.NewTestPaperGenHistoryGormMapper()

func GetAllTestPaperGenHistory(c *gin.Context) {
	info, err := testMapper.QueryAllTestPaperGenHistory()
	if err != nil {
		log.Println(err)
		utils.Make500Resp(err.Error())
	}
	utils.Make200Resp("successfully get all test paper gen history", info)
}
