package main

import (
	"log"
    "os"
    "fmt"
	"strconv"
	"strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
	
	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	msg "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {


	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

func InitDatabase() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    

	var combStr string
    combStr = user + ":" + password + "@tcp(" + host +":" + port + ")/" + dbname 
    
    db, dberr := sql.Open("mysql", combStr)
    
	if dberr != nil {
		return nil, dberr //error
    }

	return db, nil
}

func PushMessage (m *msg.Message) error {
	db, err := InitDatabase()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to init MySQL client, err: %v", err)) //log error
	}
	defer db.Close()

	chat,err := arrangeChat(m.GetChat())

	if err != nil {
		return err
	}

	text := m.GetText()
	sender := m.GetSender()
	sendtime := m.GetSendTime()

	query := `INSERT INTO chat_logs (chat, message, sender, date_time) VALUES (?, ?, ?, ?)`

	_, err = db.Exec(query, chat, text, sender, sendtime)
	if err != nil {
		return err
	}

	return nil
}

func PullMessages(pull *msg.PullRequest) ([]*msg.Message, error) {
	db, err := InitDatabase()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to init MySQL client, err: %v", err)) //log error
	}
	defer db.Close()

	chat,err := arrangeChat(pull.GetChat())

	if err != nil {	
		return nil, err
	}

	cursor := strconv.Itoa(int(pull.GetCursor()))
	limit := strconv.Itoa(int(pull.GetLimit()))
	reverse := pull.GetReverse()

	var sort string
	if reverse {
		sort = `DESC`
	} else {
		sort = `ASC`
	}

	query := `SELECT chat, date_time, sender, message FROM chat_logs WHERE chat = '` + chat + `' AND date_time >= ` + cursor + ` ORDER BY date_time ` + sort + ` LIMIT ` + limit
	fmt.Println(query)
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chat_logs []*msg.Message
	for rows.Next(){
		var log msg.Message
		err := rows.Scan(&log.SendTime, &log.Sender, &log.Text )
		if err != nil {
			return nil, err
		}
		chat_logs = append(chat_logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chat_logs, nil
}

func arrangeChat(chat string) (string, error) {
	
	var newChat string

	lowercase := strings.ToLower(chat)
	members := strings.Split(lowercase, ":")
	if len(members) != 2 {
		err := fmt.Errorf("Invalid Chat format '%s', please arrange in format of user1:user2", chat)
		return "", err
	}

	mem1, mem2 := members[0], members[1]

	if mem1 <= mem2 {
		newChat = fmt.Sprintf("%s:%s", mem1, mem2)
	} else {
		newChat = fmt.Sprintf("%s:%s", mem2, mem1)
	}

	return newChat, nil

}