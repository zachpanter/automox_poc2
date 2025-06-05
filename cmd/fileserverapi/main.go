package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}
	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Ensure connection is closed when function exits
	defer conn.Close()

	// Buffer to read incoming request
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		return
	}

	// Get the first line of the request
	requestLine := strings.SplitN(string(buf[:n]), "\r\n", 2)[0]
	parts := strings.Split(requestLine, " ")

	// Check path parameter is present
	if len(parts) < 2 {
		fmt.Fprint(conn, "HTTP/1.1 400 Bad Request\r\n\r\n")
		return
	}

	// Get image filesys name
	path := parts[1]
	filename := strings.TrimPrefix(path, "/")

	// Use root tmp directory
	fullPath := filepath.Join("/tmp", filename)

	// Get file bytes
	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		return
	}
	defer file.Close()

	// Prep file-specific MIME
	ext := filepath.Ext(fullPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream" // Binary if unknown
	}

	// Write header
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: %s\r\n\r\n", contentType)

	// Write body
	_, err = io.Copy(conn, file) // Copy file contents to the connection
	if err != nil {
		log.Printf("Failed to write file to connection: %v", err)
	}
	log.Printf("Served file %s to %s", fullPath, conn.RemoteAddr().String())
}
