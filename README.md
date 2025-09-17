token = token used to auth for notification service

temperature endpoint = port 3003 with url encoded value -> fastapi?

Running the producer

#### 1. Build the image
`docker build -t client -f Dockerfile .`

#### 2. Create the network
`docker network create producer`

#### 3. Run this line to start the consumer
`docker run --rm --name client --network producer --volume $(pwd)/auth:/usr/src/app/auth client group2 group2-group`

#### 4. Run the test producer
`docker run --rm --name producer -v ./auth:/experiment-producer/auth -v ./loads/2.json:/config.json -it --network producer dclandau/cec-experiment-producer -b kafka1.dlandau.nl:19092 --config-file /config.json --topic group2`