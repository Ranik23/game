package handlers

import (
	"game/internal/usecase"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WelcomeHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		session.Set("home_visited", true)
		session.Save()
		g.HTML(http.StatusOK, "welcome.html", gin.H{})
	}
}

func RoleHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		session.Set("roleSelection_visited", true)
		session.Save()
		g.HTML(http.StatusOK, "role.html", gin.H{})
	}
}

func LoginHandlerGET(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "login.html", gin.H{})
	}
}
func LoginHandlerPOST(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Структура для входных данных
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Привязка JSON к структуре
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_request",
			})
			return
		}

		// Получение данных из структуры
		// username := loginData.Username
		// password := loginData.Password

		// Проверка логина и пароля
		// if err := userOperator.CheckLoginInfo(username, password); err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "invalid_credentials",
		// 	})
		// 	return
		// }

		// Установка сессии
		session := sessions.Default(c)
		session.Set("login_visited", true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed_to_save_session",
			})
			return
		}

		// Отправка успешного ответа
		c.JSON(http.StatusOK, gin.H{
			"message":  "login success",
			"redirect": "/role/admin-panel",
		})
	}
}





func LoginLeaderHandlerGET(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "login-leader.html", gin.H{})
	}
}

func LoginPlayerHandlerGET(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "login-player.html", gin.H{})
	}
}


func CreateTeamHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Username string 	`json:"username"`
			TeamName string		`json:"teamName"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_request",
			})
			return
		}

		username := loginData.Username
		teamName := loginData.TeamName

		userOperator.CreateTeam(teamName)

		session := sessions.Default(c)
		session.Set("leader_login", true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed_to_save_session",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "team created",
			"redirect": "/role/leader-panel", 
			"username": username, 
			"teamName": teamName,
		})
	}
}


func LeaderPanelHanlder() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "leader-panel.html", gin.H{})
	}
}

func PlayerPanelHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "player-panel.html", gin.H{})
	}
}

func AdminPanelHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "admin-panel.html", gin.H{})
	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/role/login")
	}
}

func GetTeamsHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teams, err := userOperator.GetTeams()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"teams": teams,
		})
	}
}