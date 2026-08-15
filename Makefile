PORT=8080

ifneq (,$(wildcard .env))
include .env
export
endif

run:
	go run ./cmd/server

port:
	netstat -ano | findstr :$(PORT)

stop:
	for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8080') do taskkill /PID %%a /F

restart:
	make stop
	make run
