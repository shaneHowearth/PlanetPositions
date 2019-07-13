# PlanetPositions-WIP
Collection of microservices that can be used to determine information about various heavenly bodies for a given point on earth.

# Installation
Clone the repository then within the cloned directory run `make`.
This will create the required docker containers.

# Simple usage
A RESTful API is listening on localhost:5055, only one endpoint is currently active.
The URL has the following format localhost:5055/v1/api/Sunrise/{Longitude}/{Latitude}/{Year}/{Month}/{Day}

# Examples
`curl localhost:5055/v1/api/Sunrise/174.7633/36.8485/1994/09/03`
or
with your browser at http://localhost:5055/v1/api/Sunrise/174.7633/36.8485/1994/09/03

# Note:
This is example code, it was built to demonstrate simple gRPC connections between microservices behind a RESTful API. 

It was also used as a vehicle to learn how to use modules in go projects.
