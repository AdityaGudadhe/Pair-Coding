package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	conn, err := initMQ()
	if err != nil {
		log.Fatalf("failed to initialize MQ: %v", err)
	}
	defer conn.Close()

	const consumerCount = 5
	var wg sync.WaitGroup
	wg.Add(consumerCount)

	for i := 1; i <= consumerCount; i++ {
		consumerID := i
		go func() {
			defer wg.Done()
			consumeMsg(conn, consumerID)
		}()
	}

	wg.Wait()
	c, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var cfg S3Config = S3Config{
		Endpoint:     "http://localhost:9000",
		Region:       "us-east-1",
		AccessKey:    "miniouser",
		SecretKey:    "miniosecret",
		UsePathStyle: true,
	}
	store, err := NewStore(c, cfg)
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 0)
	for i := range 100 {
		data = append(data, byte(i))
	}
	store.PutObject(c, Test, "/test.txt", data, "text/plain")
	d, metadata, err := store.GetObject(c, Test, "test.txt")
	fmt.Print(d)
	fmt.Print("Metadata:")
	fmt.Print(metadata)
}
