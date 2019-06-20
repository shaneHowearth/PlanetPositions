all: network rest sun julian

network:
	docker network create --driver bridge planet_positions

rest:
	docker build -t rest -f restServer/Dockerfile .
	docker run -d -p 5055:5055 --name rest --net planet_positions -t rest

.PHONY: sun julian
sun:
	docker build -t sun -f sun/Dockerfile .
	docker run -d -p 5055 --name sun --net planet_positions -t sun

julian:
	docker build -t julian -f julian/Dockerfile .
	docker run -d -p 5055 --name julian --net planet_positions -t julian

clean:
	docker stop $(shell docker ps -q)
	docker rm $(shell docker ps -aq)
	docker network rm planet_positions
