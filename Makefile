all: rest sun

rest:
	docker build -t rest -f restServer/Dockerfile .
	docker run -d -p 5055:5055 -t rest

.PHONY: sun
sun:
	docker build -t sun -f sun/Dockerfile .
	docker run -d -p 5055 -t sun

clean:
	docker stop $(shell docker ps -q)
	docker rm $(shell docker ps -aq)
