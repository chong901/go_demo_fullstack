package controllers

import (
	"net/http"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func GetAccountsView(c *gin.Context) {
	accounts := []models.Account{}

	role, _ := GetSessionRole(c)

	if err := db.DB.Preload("Role").Find(&accounts, "role_id >= ?", role.Id).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "accountList.html", gin.H{
		"data": accounts,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func UpdateAccountRole(c *gin.Context) {
	account := models.Account{}

	if err := c.BindJSON(&account); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Model(&account).Update("role_id", account.RoleId).Error; err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(account, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return nil
			}

			return nh.Insert(tx)
		},
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func CreateAccount(c *gin.Context) {
	af := models.AccountForm{}
	account := models.Account{}

	if err := c.BindJSON(&af); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := af.CheckField(); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(af.Password), models.BcryptCose)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	account.Account = af.Account
	account.Password = string(hashed)
	account.Status = models.Activated
	account.RoleId = models.DefaultRoleId

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Create(&account).Error; err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(account, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		},
	); err != nil {
		if constants.IsDuplicatedErr(err) {
			GoToErrorPage(http.StatusBadRequest, c, models.ErrAccountDuplicated)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func UpdateAccountPassword(c *gin.Context) {
	af := models.AccountForm{}

	if err := c.BindJSON(&af); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := af.CheckField(); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(af.Password), models.BcryptCose)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	account := models.Account{}
	account.Account = af.Account

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Model(&account).Update("password", string(hashed)).Error; err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(account, models.ChangePassword, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		},
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func Login(c *gin.Context) {
	af := models.AccountForm{}
	account := models.Account{}

	if err := c.BindJSON(&af); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := af.CheckField(); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := db.DB.Preload("Role.Functions.Function").Find(&account, "account = ?", af.Account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			GoToErrorPage(http.StatusUnauthorized, c, models.ErrNoAccount)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(af.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			GoToErrorPage(http.StatusUnauthorized, c, models.ErrWrongPassword)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	session := sessions.Default(c)
	session.Set(constants.SessionRole, account.Role)
	session.Set(constants.SessionLoginUser, account.Account)

	if err := session.Save(); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(nil, models.LoginAction, account.Account, constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := nh.Insert(db.DB); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
