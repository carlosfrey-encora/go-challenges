package main

import (
	"context"
	"flag"
	"log"
	"time"

	"crud/internal/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newtask := &pb.Task{Name: "Go watch some film with my friends", Completed: false}

	tasks, err := c.CreateTask(ctx, &pb.CreateTaskRequest{Task: newtask})

	if err != nil {
		log.Fatalf("could not retrieve task: %v", tasks)
	}

	log.Printf("Task: %v", tasks)
}
