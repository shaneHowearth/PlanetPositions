all: rest sun julian

rest:
	docker build -t rest -f restServer/Dockerfile .
	docker run -d -p 5055:5055 -t rest

.PHONY: sun julian
sun:
	docker build -t sun -f sun/Dockerfile .
	docker run -d -p 5055 -t sun

julian:
	docker build -t julian -f julian/Dockerfile .
	docker run -d -p 5055 -t julian

clean:
	docker stop $(shell docker ps -q)
	docker rm $(shell docker ps -aq)
