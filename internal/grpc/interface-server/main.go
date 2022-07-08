package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	gormdb "crud/internal/gorm-db"
	"crud/internal/grpc/pb"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedDatabaseServer
	db *gormdb.OrmCrudOperations
}

func (s *server) RetrieveTask(ctx context.Context, in *pb.TaskRequestById) (*pb.TaskResponseById, error) {

	task, err := s.db.GetTaskById(int(in.GetId()))

	if err != nil {

		return nil, fmt.Errorf("RetrieveTask: %v", task)
	}
	return &pb.TaskResponseById{Task: task}, nil
}

func (s *server) DeleteTask(ctx context.Context, in *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {

	_, err := s.db.DeleteTask(int(in.GetId()))

	if err != nil {

		return nil, fmt.Errorf("DeleteTask: %v", err)
	}

	return &pb.DeleteTaskResponse{}, nil
}

func (s *server) RetrieveByCompletion(ctx context.Context, in *pb.TaskRequestByCompletion) (*pb.TaskResponseByCompletion, error) {

	tasks, err := s.db.GetTaskByCompletion(in.GetCompleted())

	if err != nil {

		return nil, fmt.Errorf("RetrieveByCompletion: %v", err)
	}

	return &pb.TaskResponseByCompletion{Tasks: tasks}, nil
}

func (s *server) RetrieveAllTasks(ctx context.Context, in *pb.TaskRequestAll) (*pb.TaskResponseAll, error) {

	tasks, err := s.db.ListAll()

	if err != nil {

		return nil, fmt.Errorf("RetrieveAllTasks: %v", err)
	}

	return &pb.TaskResponseAll{Tasks: tasks}, nil
}

func (s *server) PutTask(ctx context.Context, in *pb.PutTaskRequest) (*pb.PutTaskResponse, error) {

	_, err := s.db.UpdateTask(int64(in.GetId()), *in.GetTask())

	if err != nil {

		return nil, fmt.Errorf("PutTask: %v", err)
	}

	return &pb.PutTaskResponse{}, nil
}

func (s *server) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {

	id, err := s.db.CreateTask(*in.GetTask())

	if err != nil {

		return nil, fmt.Errorf("PutTask: %v", err)
	}

	return &pb.CreateTaskResponse{Id: int32(id)}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterDatabaseServer(s, &server{db: gormdb.Connect()})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
