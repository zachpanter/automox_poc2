run:
	go run cmd/fileserverapi/main.go & \
	sleep 2 && open http://localhost:8080/waves.jpg