all: rest

rest:
	docker build -t rest -f restServer/Dockerfile .
	docker run -d -p 5055:5055 -t rest

clean:
	docker stop $(shell docker ps -q)
	docker rm $(shell docker ps -aq)
