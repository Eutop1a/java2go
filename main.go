package main

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"java2go/config"
	"java2go/controller"
	"log"
	"time"
)

func main() {
	r := gin.Default()
	gob.Register(time.Time{})
	store := cookie.NewStore([]byte("secret"))
	// 设置会话过期时间为 1 小时
	store.Options(sessions.Options{
		MaxAge: 3600,
	})
	r.Use(sessions.Sessions("mysession", store))
	r.Use(config.Cors())

	r.GET("/permission_denied", controller.PermissionDenied)
	r.GET("/getLoginStatus", controller.GetLoginStatus)
	r.POST("/login", controller.Login)
	r.POST("/logout", controller.Logout)
	r.POST("/registered", controller.Registered)
	r.GET("/getApplyUser", controller.GetApplyUser)
	r.GET("/getAllUser", controller.GetAllUser)
	r.GET("/deleteUser", controller.DeleteUser)
	r.GET("/passApply", controller.PassApply)
	r.GET("/deleteApply", controller.DeleteApply)

	r.GET("/getQuestionGenHistoriesByTestPaperUid", controller.GetQuestionGenHistoriesByTestPaperUid)
	r.GET("/deleteQuestionGenHistoryByTestPaperUid", controller.DeleteQuestionGenHistoryByTestPaperUid)
	r.POST("/updateQuestionGenHistory", controller.UpdateQuestionGenHistory)
	r.GET("/reExportTestPaper", controller.ReExportTestPaper)
	r.GET("/exportAnswer", controller.ExportAnswer)
	r.GET("/getAllQuestionLabels", controller.GetAllQuestionLabels)
	r.GET("/getDistinctChapter1", controller.GetDistinctChapter1)
	r.GET("/getDistinctChapter2", controller.GetDistinctChapter2)
	r.GET("/getChapter2ByChapter1", controller.GetChapter2ByChapter1)
	r.GET("/getDistinctLabel1", controller.GetDistinctLabel1)
	r.GET("/getDistinctLabel2", controller.GetDistinctLabel2)

	qBan := controller.NewQuestionBankController()
	r.GET("/getAllQuestionBank", qBan.GetQuestionBank)
	r.GET("/getQuestionBank", qBan.GetQuestionBank)
	r.GET("/getTopicType", qBan.GetTopicType)
	r.GET("/searchQuestionByTopic", qBan.SearchQuestionByTopic)
	r.POST("/insertSingleQuestionBank", qBan.InsertSingleQuestionBank)
	r.GET("/deleteSingleQuestionBank", qBan.DeleteSingleQuestionBank)
	r.GET("/getQuestionBankById", qBan.GetQuestionBankById)
	r.POST("/updateQuestionBankById", qBan.UpdateQuestionBankById)
	r.POST("/upload", qBan.UploadFile)
	r.GET("/getEachChapterCount", qBan.GetEachChapterCount)
	r.GET("/getEachScoreCount", qBan.GetEachScoreCount)

	r.GET("/getAllTestPaperGenHistory", controller.GetAllTestPaperGenHistory)

	r.POST("/randomSelect", controller.RandomSelect)
	r.POST("/geneticSelect", controller.GeneticSelect)
	r.POST("/questionGen", controller.QuestionGen)
	r.POST("/questionGen2", controller.QuestionGen2)
	r.POST("/getFile", controller.GetFile)

	log.Println("Server started on :8081")
	r.Run(":8081")
}
