.PHONY: build test up down

build: ayd-smb-probe

ayd-smb-probe: main.go go.sum
	go build .

test:
	go test -tags=integration -race .

up:
	docker compose -f testdata/docker-compose.yaml up -d

down:
	docker compose -f testdata/docker-compose.yaml down
