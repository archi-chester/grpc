package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/archi-chester/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client\n")

	// tls := true
	opts := grpc.WithInsecure()
	// if tls {
	// 	certFile := "ssl/ca.crt"
	// 	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	// 	if sslErr != nil {
	// 		log.Fatalf("Failed loading CA certificate: %v", sslErr)
	// 		return
	// 	}
	// 	opts = grpc.WithTransportCredentials(creds)
	// }

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect :%v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create blog
	fmt.Println("Creating the blog\n")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("unexpected error :%v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// read Blog
	fmt.Println("Reading the blog\n")

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: ""})
	if err != nil {
		fmt.Printf("unexpected error :%v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("unexpected error :%v\n", err)
	}
	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update Blog
	fmt.Println("Updating the blog\n")
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Stephane",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v\n", updateErr)
	}
	fmt.Printf("Blog updating: %v\n", updateRes)

	// delete blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v\n", deleteErr)
	}

	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list Blog
	fmt.Println("List blog")
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
