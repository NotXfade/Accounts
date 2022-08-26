package routes

import (
	"net/http"
	"os"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/src/api"

	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Routes is a method in which all the service endpoints are defined
func Routes(router *gin.Engine) {

	//Healthz Check
	router.GET("/healthz", api.Healthz)

	router.StaticFile("/openapi.yaml", "./openapi.yaml")
	// developers help endpoint
	if config.Conf.Service.Environment != "production" {
		// endpoint to read variables
		router.GET("/end", checkToken, readEnv)
		router.GET("/logs", checkToken, readLogs)
	}
	//setting up middleware for protected apis
	authMiddleware := token.MwInitializer()
	// this group will contain all the accounts related api's
	v1 := router.Group("/v1")
	{
		v1.POST("/signup", api.Register)
		v1.POST("/admin-signup", api.RegisterAdmin)
		v1.POST("/login", api.LoginEndPoint)
		v1.POST("/resetpassword", api.ResetPassword)
		v1.PUT("/setpassword", api.SetPassword)
		//==== Google ====
		v1.GET("/google/callback", api.GoogleCallback)
		v1.GET("/google/login", api.GoogleLogin)
		//============== Protected routes ================
		v1.Use(authMiddleware.MiddlewareFunc())
		{
			v1.POST("/register", api.StoreEmployeeInfo)
			v1.GET("/refreshtoken", api.RefreshToken)
			v1.GET("/logout", api.Logout)
			v1.POST("/acceptpolicy", api.AcceptPolicy)
			v1.PUT("/changepassword", api.ChangePassword)
			//============= User Api's ===================
			v1.GET("/user/profile", api.ViewProfile)
			v1.PUT("/user/profile", api.UpdateProfile)
			v1.GET("/user/profilepicture", api.GetProfilePicture)
			v1.GET("/batch", api.GetBatch)
			//=============== Reviwer, BA,People and culture team,admin and superadmin view
			v1.Use(checkReviewer)
			{
				v1.GET("/adminlist", api.ListAdmin)
				v1.GET("/userprofile/:email", api.ViewUserProfile)
				v1.GET("/internslist/:level", api.InternsList)
				v1.POST("/google/invite", api.CalendarInvite)
				v1.GET("/reports/:level", api.GetReports)
			}
			//=============== People and culture team,admin and superadmin view ==============
			v1.Use(checkHr)
			{
				v1.POST("/department", api.CreateDepartment)
				v1.PUT("/department/:id", api.UpdateDepartment)
				v1.DELETE("/department/:id", api.DeleteDepartment)
				v1.GET("/department", api.ListDepartment)
				v1.POST("/changelevel", api.ChangeLevel)
				v1.PUT("/accountstatus/:status/userid/:email", api.AccountStatus)
				v1.POST("/invite", api.Invite)
				v1.POST("/assignreviewer", api.AssignReviewer)
				v1.GET("/listreviewer/:email", api.GetListOfAssignedReviewers)
				v1.DELETE("/deletereviewer/:email/user/:uemail", api.DeleteUserAssignedReviewer)
				v1.GET("/google/check", api.CheckCredentials)
				v1.GET("/intern-list", api.InternList)
			}
			//============== Admin and superadmin use Api's =============
			v1.Use(checkAdmin)
			{
				v1.POST("/admininvite", api.InviteAdmin)
				v1.DELETE("/delete/:email", api.Delete)
				v1.DELETE("/admin/:email", api.DeleteAdmin)
			}
		}
	}
	router.Use(token.MwInitializerWithQuery().MiddlewareFunc())
	{
		router.GET("/export", api.Export)
	}
}

//func checkSuperAdmin is used to check if role is super admin
func checkSuperAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := int(claims["id"].(float64))
	role := claims["role"].(string)
	if role != "admin" && id != 1 {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

//func checkAdmin is used to check user is admin or administrator
func checkAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	role := claims["role"].(string)
	if role != "admin" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

//checkHr is used to check user is admin,P&C,adminstrator
func checkHr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	role := claims["role"].(string)
	if role != "p&c" && role != "admin" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

//checkBA is used to check user is business analyst,admin,P&C,adminstrator
func checkBA(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	role := claims["role"].(string)
	if role != "business analyst" && role != "p&c" && role != "admin" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

//checkReviewer is used to check user is business analyst,admin,P&C,adminstrator,reviewer
func checkReviewer(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	role := claims["role"].(string)
	if role != "reviewer" && role != "business analyst" && role != "p&c" && role != "admin" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

// readLogs is a api handler for reading logs
func readLogs(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "info.txt")
}

// readEnv is api handler for reading configuration variables data
func readEnv(c *gin.Context) {
	if config.TomlFile == "" {
		// if configuration is done using environment variables
		env := make([]string, 0)
		for _, pair := range os.Environ() {
			env = append(env, pair)
		}
		c.JSON(200, gin.H{
			"environments": env,
		})
	} else {
		// if configuration is done using toml file
		http.ServeFile(c.Writer, c.Request, config.TomlFile)
	}
}

// checkToken is a middleware to check header is set or not for secured api
func checkToken(c *gin.Context) {
	xt := c.Request.Header.Get("XSOnboarding-token")
	if xt != "XSOnboarding" {
		c.Abort()
		c.JSON(404, gin.H{})
		return
	}
	c.Next()
}
