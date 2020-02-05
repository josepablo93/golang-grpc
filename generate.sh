#!/bin/bash

echo "Running $1"
if [ $1 == "calc" ]; then  
  protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
elif [ $1 == "greet" ]; then
  protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
elif [ $1 == "blog" ]; then
  protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.
else 
  echo "Running all proto files"
  protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
  protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
  protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.
fi