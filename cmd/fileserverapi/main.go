package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	// ctx := context.Background()
	// apiImpl := api.NewAPI(ctx)

	// restful.Add(apiImpl.WS)

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}
	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}

	// log.Println("Starting server on :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Fetch the image from the remote URL
	resp, err := http.Get("https://placehold.co/600x400/png")
	if err != nil {
		log.Printf("Failed to fetch image: %v", err)
		return
	}
	defer resp.Body.Close()

	// Try to detect the content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Write HTTP response headers
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: %s\r\n\r\n", contentType)

	// Write the image body to the connection
	// TODO: Instead of io.Copy, build my own response body and then serve it
	_, err = io.Copy(conn, resp.Body)
	if err != nil {
		log.Printf("Failed to write image to connection: %v", err)
	}

	conn.Write([]byte(" HTTP/1.1 200 OK"))
	log.Printf("New connection from %s", conn.RemoteAddr().String())
	// For now, just close the connection
}
