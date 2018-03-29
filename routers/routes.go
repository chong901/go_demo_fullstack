package routers

import (
	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/controllers"
	"github.com/aaa59891/go_demo_fullstack/middlewares"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {

	view := r.Group("")
	view.Use(middlewares.I18nTranslate)
	view.GET("/", controllers.Index)

	r.GET("/logout", controllers.Logout)
	r.POST("/login", controllers.Login)

	authorized := r.Group("/")

	authorized.Use(middlewares.AuthRequired)
	{
		authorizedView := view.Group("")
		authorizedView.Use(middlewares.AuthRequired)
		{
			authorizedView.GET("/welcome", controllers.Welcome)
			authorizedView.GET("/inventory", controllers.GetInventoriesView)
			authorizedView.GET("/inventory/history", controllers.GetInventoryHistoryView)
			authorizedView.GET("/skuM", controllers.GetSkusView)
			authorizedView.GET("/skuM/save", controllers.EditSkuView)
			authorizedView.GET("/jobM/save", controllers.EditJobView)
			authorizedView.GET("/jobM", controllers.JobListView)
			authorizedView.GET("/jobSchedule", controllers.JobScheduleView)
			authorizedView.GET("/poM", controllers.GetPosView)
			authorizedView.GET("/poM/save", controllers.EditPoView)
			authorizedView.GET("/function", controllers.GetFuncsView)
			authorizedView.GET("/machineM", controllers.MachineListView)
			authorizedView.GET("/machineM/save", controllers.EditMachineView)
			authorizedView.GET("/role", controllers.GetRoleSettingView)
			authorizedView.GET("/account", controllers.GetAccountsView)
			authorizedView.GET("/jobIdRuleList", controllers.IdRuleListView(models.RuleCategory_Job))
			authorizedView.GET("/jobIdRule", controllers.IdRuleView(models.RuleCategory_Job))
			authorizedView.GET("/configuration/skuParameters", controllers.GetConfigurationByType(models.Parameters))
			authorizedView.GET("/configuration/unit", controllers.GetConfigurationByType(models.Unit))
			authorizedView.GET("/reportOutput", controllers.GetRecordsView)
			authorizedView.GET("/reportOutput/history", controllers.GetRecordHistoryView)
		}

		authorized.PUT("/inventory", controllers.UpdateInventoryLevel)

		authorized.POST("/skuM/save", controllers.SaveSku)
		authorized.POST("/skuM/uploadFile", controllers.UploadFile)
		authorized.DELETE("/skuParameter", controllers.DeleteParam)
		authorized.DELETE("/skuM", controllers.DeleteSku)

		authorized.POST("/jobM/save", controllers.SaveJob)
		authorized.DELETE("/jobM", controllers.DeleteJob)

		authorized.PUT("/jobSchedule", controllers.SaveJobSchedule)

		authorized.POST("/poM/save", controllers.SavePo)
		authorized.DELETE("/poM", controllers.DeletePo)
		authorized.DELETE("/poSku", controllers.DeletePoSku)

		authorized.POST("/function", controllers.SaveFunc)
		authorized.DELETE("/function", controllers.DeleteFunc)

		authorized.POST("/machineM/save", controllers.CreateMachine)
		authorized.PUT("/machineM/save", controllers.UpdateMachine(models.UpdateColumnsFromBrowser, constants.FromBrowser))
		authorized.DELETE("/machineM", controllers.DeleteMachine)

		authorized.POST("/role", controllers.CreateNewRole)
		authorized.DELETE("/role", controllers.DeleteRole)
		authorized.PUT("/role/name", controllers.UpdateRoleName)
		authorized.PUT("/role/functions", controllers.UpdateRoleFunctions)

		authorized.PUT("/account/role", controllers.UpdateAccountRole)
		authorized.POST("/account", controllers.CreateAccount)
		authorized.PUT("/account/password", controllers.UpdateAccountPassword)

		authorized.POST("/configuration", controllers.AddConfiguration)
		authorized.DELETE("/configuration", controllers.DeleteConfiguration)

		authorized.DELETE("/jobIdRule", controllers.DeleteIdRule(models.RuleCategory_Job))
		authorized.PUT("/idRule", controllers.UpdateIdRule)
		authorized.POST("/idRule", controllers.CreateIdRule)

		authorized.POST("/reportOutput", controllers.CreateRecord)
		authorized.PUT("/reportOutput", controllers.UpdateRecord)
		authorized.DELETE("/reportOutput", controllers.DeleteRecord)
	}

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/configuration/prodParameters", controllers.GetConfigurationByTypeApi(models.Parameters))
		apiV1.GET("/configuration/units", controllers.GetConfigurationByTypeApi(models.Unit))
		apiV1.GET("/sku", controllers.GetAllSkusApi)
		apiV1.GET("/function/funcDict", controllers.GetFuncsDictApi)
		apiV1.GET("/function", controllers.GetFuncByIdApi)
		apiV1.GET("/user/functions", controllers.GetFunctionsSessionApi)
		apiV1.PUT("/machine", controllers.UpdateMachine(models.UpdateColumnsFromPi, constants.FromPi))
		apiV1.GET("/machine/summary", controllers.GetMachineSummaryApi)

		apiV1.GET("/jobIdRuleList", controllers.IdRuleListApi(models.RuleCategory_Job))

		apiV1.GET("/jobM/partQtyList", controllers.GetPartQtyApi)
		apiV1.GET("/jobPlanned", controllers.GetPlannedJobsApi)
		apiV1.GET("/jobUndone", controllers.GetUndoneJobsApi)
		apiV1.GET("/lang", controllers.ChangeLang)

		apiV1.POST("/recordOutput", controllers.CreateRecord)

		authorizedApi := apiV1.Group("")
		authorizedApi.Use(middlewares.AuthRequired)
		{
			authorizedApi.GET("/role/all", controllers.GetRolesApi)
			authorizedApi.GET("/role", controllers.GetRoleApi)
		}
	}
}
