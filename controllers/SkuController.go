package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/db"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetSkusView(c *gin.Context) {
	currentPafe, err := GetCurrentPage(c)
	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	sf := models.SkuFilter{}
	if err := c.BindQuery(&sf); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	skus := []models.Sku{}
	db := db.DB.Preload("Parameters")
	db = GetFilterCriteria(sf, db)
	pg, err := Pagination(10, currentPafe, db, &skus)
	//if err := db.DB.Preload("Parameters").Find(&skus).Error; err != nil {
	//	GoToErrorPage(http.StatusInternalServerError, c, err)
	//	return
	//}
	c.HTML(http.StatusOK, "skuMList.html", gin.H{
		"data": skus,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		"pg":     pg,
		"filter": sf,
	})
}

func EditSkuView(c *gin.Context) {
	id := c.Query("id")
	sku := models.Sku{}
	c.BindQuery(&sku)

	unitOptions := []models.Configuration{}
	if err := db.DB.Find(&unitOptions, "type = ?", 2).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	parameterConfig := []models.Configuration{}
	if len(id) != 0 {
		if err := db.DB.Preload("Parameters").Find(&sku, "id = ?", id).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
		if err := db.DB.Find(&parameterConfig, "type = ?", models.Parameters).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
	}

	c.HTML(http.StatusOK, "skuMEdit.html", gin.H{
		"data":                    sku,
		"paramConfig":             parameterConfig,
		"units":                   unitOptions,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func SaveSku(c *gin.Context) {
	sku := models.Sku{}

	if err := c.BindJSON(&sku); err != nil {
		log.Panic(err)
	}

	sku.User = GetSessionLoginUser(c)

	if err := models.Transactional(sku.Save(constants.FromBrowser)); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": sku.Id,
	})
}

func UploadFile(c *gin.Context) {
	id := c.PostForm("id")

	gsUri, err := saveFile(c.Request, "gsFile", id+"_gsFile")
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	sopUri, err := saveFile(c.Request, "sopFile", id+"_sopFile")
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}
	sku := models.Sku{}
	sku.Id = id

	update := make(map[string]interface{})

	if len(gsUri) != 0 {
		update["gsUri"] = gsUri
	}
	if len(sopUri) != 0 {
		update["sopUri"] = sopUri
	}

	if len(update) > 0 {
		newSku := sku
		newSku.GsUri = gsUri
		newSku.SopUri = sopUri

		if err := models.Transactional(
			func(tx *gorm.DB) error {
				if err := tx.Model(&sku).Update(update).Error; err != nil {
					return err
				}

				nh, err := models.CreateNormalHistoryByModel(newSku, models.SkuUploadFile, GetSessionLoginUser(c), constants.FromBrowser)

				if err != nil {
					return err
				}

				return nh.Insert(tx)
			},
		); err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, nil)
}

func saveFile(req *http.Request, key, filename string) (string, error) {
	file, header, err := req.FormFile(key)
	if err != nil {
		if err.Error() != "http: no such file" {
			log.Fatal(err)
			return "", err
		}
		return "", nil
	}

	oriFilename := header.Filename
	splitByDot := strings.Split(oriFilename, ".")
	ext := splitByDot[len(splitByDot)-1]

	newFilename := fmt.Sprintf("./static/uploads/%s.%s", filename, ext)

	out, err := os.Create(newFilename)

	defer out.Close()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	_, err = io.Copy(out, file)

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return strings.TrimLeft(newFilename, "."), nil
}

func DeleteParam(c *gin.Context) {
	sp := models.SkuParameter{}

	if err := c.BindQuery(&sp); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := db.DB.Find(&sp).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(sp, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(
		nh.Insert,
		models.DeleteById(&models.SkuParameter{}, sp.Id),
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetAllSkusApi(c *gin.Context) {
	skus := []models.Sku{}

	if err := db.DB.Find(&skus).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": skus,
	})
}

func DeleteSku(c *gin.Context) {
	sku := models.Sku{}

	if err := c.BindQuery(&sku); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := db.DB.Find(&sku).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(sku, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(nh.Insert, sku.Remove); err != nil {
		if err == models.ErrInventoryNotEmpty {
			GoToErrorPage(http.StatusBadRequest, c, err)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}
