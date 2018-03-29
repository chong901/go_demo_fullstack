package inits

import (
	"encoding/gob"
	"log"

	"github.com/aaa59891/mosi_demo_go/db"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/googollee/go-socket.io"
)

var SocketServer *socketio.Server

func init() {
	var err error
	SocketServer, err = socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	SocketServer.On("connection", func(so socketio.Socket) {
		log.Println("socket connected.")
		so.Join("mosi")
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	SocketServer.On("error", func(so socketio.Socket, err error) {
		log.Panic("Socket server had an error, ", err)
	})
}

func CreateTable() {
	modelArr := make([]interface{}, 0)
	modelArr = append(modelArr, models.Account{})
	modelArr = append(modelArr, models.Function{})
	modelArr = append(modelArr, models.InventoryHistory{})
	modelArr = append(modelArr, models.Inventory{})
	modelArr = append(modelArr, models.Job{})
	modelArr = append(modelArr, models.Po{})
	modelArr = append(modelArr, models.PoSku{})
	modelArr = append(modelArr, models.RoleFunction{})
	modelArr = append(modelArr, models.Role{})
	modelArr = append(modelArr, models.SkuParameter{})
	modelArr = append(modelArr, models.Sku{})
	modelArr = append(modelArr, models.Configuration{})
	modelArr = append(modelArr, models.Machine{})
	modelArr = append(modelArr, models.MachineStatus{})
	modelArr = append(modelArr, models.NormalHistory{})
	modelArr = append(modelArr, models.IdRule{})
	modelArr = append(modelArr, models.RecordOutput{})

	for _, model := range modelArr {
		db.DB.Set("gorm:table_options", "CHARACTER SET = utf8").AutoMigrate(model)
	}
}

func RegisterStruct() {
	gob.Register(&models.Role{})
}
