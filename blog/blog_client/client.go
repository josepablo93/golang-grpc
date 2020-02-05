package main

import (
	"context"
	"fmt"
	"grpc_go_course/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm a blog client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "JP",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// read blog function
	fmt.Println("Read the blog")
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "198231uj213hk"})

	if err2 != nil {
		fmt.Printf("Error happended while reading %v\n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBLogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBLogErr != nil {
		fmt.Printf("Error happended while reading %v\n", readBLogErr)
	}

	fmt.Println("Blog was read ", readBlogRes)

	// update blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Jose Flores",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog (edited)",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog})

	if updateErr != nil {
		fmt.Printf("Error happended while updating %v\n", err2)
	}

	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	})

	if deleteErr != nil {
		fmt.Printf("Error happended while deleting %v\n", deleteErr)
	}

	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list blogs
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
