package main

import (
	"context"
	"log"
	"time"
	"strconv"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/kitex_gen/rpc"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/kitex_gen/rpc/imservice"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/proto_gen/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var cli imservice.Client

func main() {
	r, err := etcd.NewEtcdResolver([]string{"etcd:2379"})
	if err != nil {
		log.Fatal(err)
	}
	cli = imservice.MustNewClient("demo.rpc.server",
		client.WithResolver(r),
		client.WithRPCTimeout(1*time.Second),
		client.WithHostPorts("rpc-server:8888"),
	)

	h := server.Default(server.WithHostPorts("0.0.0.0:8080"))

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.POST("/api/send", sendMessage)
	h.GET("/api/pull", pullMessage)

	h.Spin()
}

func sendMessage(ctx context.Context, c *app.RequestContext) {
	var req api.SendRequest
	err := c.Bind(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body: %v", err)
		return
	}
    
    sender := c.Query("sender")
    receiver := c.Query("receiver")
    text := c.Query("text")
    sendTime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))

	resp, err := cli.Send(ctx, &rpc.SendRequest{
		Message: &rpc.Message{
			Chat:   sender+":"+receiver ,
			Text:   text,
			Sender: sender,
            SendTime: sendTime.Unix(),
		},
	})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
	} else {
		c.Status(consts.StatusOK)
	}
}

func pullMessage(ctx context.Context, c *app.RequestContext) {
	var req api.PullRequest
	err := c.Bind(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body: %v", err)
		return
	}

	chat := c.Query("chat")
	tempCursor, err := strconv.ParseInt(c.Query("cursor"),10,64)
	var cursor int64
	if tempCursor == 0 {
		cursor = 0
	} else {
		cursor = tempCursor
	}
	
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, "Failed to parse cursor: %v", err)
	// 	return
	// }


	tempLimit,err  := strconv.ParseInt(c.Query("limit"),10,32)
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, "Failed to parse limit: %v", err)
	// 	return
	// }


	var limit int32
	if tempLimit == 0 {
		limit = 10
	} else {
		limit = int32(tempLimit)
	}
	

	var reverse bool
	tempReverse,err := strconv.ParseBool(c.Query("reverse"))
	if c.Query("reverse") == "" {
		reverse = false
		err = nil
	} else {
		reverse = tempReverse
	}
	
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse reverse: %v", err)
		return
	}

	resp, err := cli.Pull(ctx, &rpc.PullRequest{
		Chat:    chat,
		Cursor:  cursor,
		Limit:   limit,
		Reverse: &reverse,
	})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
		return
	}

	messages := make([]*api.Message, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, &api.Message{
			Chat:     msg.Chat,
			Text:     msg.Text,
			Sender:   msg.Sender,
			SendTime: msg.SendTime,
		})
	}
	c.JSON(consts.StatusOK, &api.PullResponse{
		Messages:   messages,
		HasMore:    resp.GetHasMore(),
		NextCursor: resp.GetNextCursor(),
	})
}
