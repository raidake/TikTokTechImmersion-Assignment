package main

import (
	"context"
	//"math/rand"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()
	err := PushMessage(req.GetMessage())
	if err != nil {
		resp.Code = 500
		resp.Msg = "oops" 
		return resp, err
	}

	resp.Code = 0
	resp.Msg = "success"
	return resp, nil

}	

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	messages, err := PullMessages(req)
	if err != nil {
		return nil, err
	}

	resp.SetMessages(messages)

	resp.Code = 0
	resp.Msg = "success"
	return resp, nil
}

// func areYouLucky() (int32, string) {
// 	if rand.Int31n(2) == 1 {
// 		return 0, "success"
// 	} else {
// 		return 500, "oops"
// 	}
// }
