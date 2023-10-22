test:
	go test ./... -v
	
fugazi:
	go test -run ^TestFugazi$ . -v