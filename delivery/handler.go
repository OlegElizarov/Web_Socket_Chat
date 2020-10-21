package delivery

import (
	"Web_Socket_Chat/models"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/pgxpool"
	"github.com/labstack/echo"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	upgrader = websocket.Upgrader{}
)

type Handler struct {
	Db *pgxpool.Pool
}

func (h *Handler) Receive(c echo.Context) error {
	//Receive from client
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	//fmt.Println("body: ", string(body))
	var mes models.Message
	err = mes.UnmarshalJSON(body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	mes.Shown = false
	fmt.Println(mes)
	//write to DB
	WriteMessage(mes, h.Db)
	return c.JSON(200, "status:Gotcha")
}

func (h *Handler) Send(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Read from client
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return err
			//fmt.Println(err)
		}
		fmt.Printf("%s\n", msg)
		//
		data := strings.Fields(string(msg))
		// Write to client
		var mess models.Messages
		//var dialog models.Dialog
		sender, _ := strconv.Atoi(data[0])
		receiver, _ := strconv.Atoi(data[1])
		//read from db
		//dial.UnmarshalJSON(msg)
		fmt.Println(sender, receiver)
		mess = GetUnshownMessages(receiver, sender, h.Db)
		for _, val := range mess {
			if !val.Shown {
				err = ws.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(val.Sender)+": "+val.Message_line))
			} else {
				err = ws.WriteMessage(websocket.TextMessage, []byte{})
			}
		}
		if err != nil {
			//fmt.Println(err)
			return err
		}

	}
}

func WriteMessage(mes models.Message, db *pgxpool.Pool) {
	//fmt.Println(mes)
	sql := `insert into message values (default,$1,$2,$3,$4)`
	queryRes, err := db.Exec(context.Background(), sql, mes.Sender, mes.Receiver, mes.Message_line, mes.Shown)
	if err != nil {
		fmt.Println(err)
		return
	}
	if queryRes.RowsAffected() == 0 {
		fmt.Println("Not affected")
		return
	}
}

func GetMessages(id int, db *pgxpool.Pool) models.Messages {
	messages := models.Messages{}
	sql := `select * from message where receiver = $1`
	queryRes, err := db.Query(context.Background(), sql, id)
	for queryRes.Next() {
		mes := &models.Message{}
		err = queryRes.Scan(mes.ID, mes.Sender, mes.Receiver, mes.Message_line, mes.Shown)
		if err != nil {
		}
		messages = append(messages, *mes)
	}
	return messages
}

func GetUnshownMessages(id_sender, id_receiver int, db *pgxpool.Pool) models.Messages {
	messages := models.Messages{}
	fmt.Println("params: ", id_receiver, id_sender)
	sql := `select * from message where receiver = $1 AND sender = $2 AND shown = false`
	queryRes, err := db.Query(context.Background(), sql, id_receiver, id_sender)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for queryRes.Next() {
		mes := models.Message{}
		err = queryRes.Scan(&mes.ID, &mes.Sender, &mes.Receiver, &mes.Message_line, &mes.Shown)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		//fmt.Println(mes)
		messages = append(messages, mes)
	}
	MakeShown(id_receiver, id_sender, db)
	return messages
}

func MakeShown(id_receiver, id_sender int, db *pgxpool.Pool) {
	sql := "update message set shown = true where shown = false and receiver = $1 and sender = $2"
	queryRes, _ := db.Exec(context.Background(), sql, id_receiver, id_sender)
	if queryRes.RowsAffected() != 1 {
		return
	}
}
