package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "go-dapr-grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type todoListManage struct {
	m *sync.RWMutex
	*pb.TodoListResponse
}

const (
	port = 4999
)

type server struct {
	pb.UnimplementedTodoListServer
	todoListMng todoListManage
}

func newServer() *server {
	initData := &pb.TodoListData{
		Title:     "init todolist",
		CreatedAt: uint64(time.Now().Unix()),
	}

	return &server{
		todoListMng: todoListManage{
			m: new(sync.RWMutex),
			TodoListResponse: &pb.TodoListResponse{
				TodoLists: []*pb.TodoListData{
					initData,
				},
				Size: 1,
			},
		},
		UnimplementedTodoListServer: pb.UnimplementedTodoListServer{},
	}
}

func (s *server) GetTodolist(ctx context.Context, empty *emptypb.Empty) (*pb.TodoListResponse, error) {
	s.todoListMng.m.RLock()
	defer s.todoListMng.m.RUnlock()
	return s.todoListMng.TodoListResponse, nil
}

func addTodoListHandle(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type bodyData struct {
			Title string `json:"title"`
		}
		body := &bodyData{}
		if err := readReq(r, body); err != nil {
			log.Fatalln(err)
		}

		s.todoListMng.m.Lock()
		defer s.todoListMng.m.Unlock()
		s.todoListMng.TodoListResponse.TodoLists = append(s.todoListMng.TodoListResponse.TodoLists, &pb.TodoListData{
			Title:     body.Title,
			CreatedAt: uint64(time.Now().Unix()),
		})
		s.todoListMng.TodoListResponse.Size = uint32(len(s.todoListMng.TodoListResponse.TodoLists))

		writeRes(w, body)
	}
}

func main() {
	serverStruct := newServer()

	go func(s *server) {
		log.Printf("starting http server on port %d", 8080)
		http.HandleFunc("/add-todo-list", addTodoListHandle(s))
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			writeRes(w, "pong")
		})
		http.ListenAndServe(":8080", nil)
	}(serverStruct)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTodoListServer(s, serverStruct)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// =============================== HTTP Functional ===============================

func writeRes(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		log.Fatalln(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(contentJson)
}

func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}
