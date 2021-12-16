all:
	go build -o velo_e2e ./cmd/

windows:
	GOOS=windows go build -o velo_e2e.exe ./cmd/
