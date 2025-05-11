## HashRing implements consistent hashing with partitioning </br></br>

# **Flow Diagram:**
</br>

## | Client    +------->+ HashRing      +------->+ Node Responsible|

## | PUT/GET   |  -      | hashKey(key)  |        | for that hash   |          | Store/Retrieve KV|

</br>                                  

# Description:
### 1. The client issues a PUT or GET request with a key.
### 2. HashRing hashes the key using SHA-1 (first 4 bytes).
### 3. Sorted node hashes determine which node the key maps to.
### 4. Data is routed to the appropriate node.
### 5. Partitioning is achieved by consistent hashing, which distributes keys evenly across nodes by mapping them to a circular hash space. </br>
###    When a node is added or removed, only a subset of keys needs to be remapped, minimizing disruption and ensuring scalability.

## Run code 
> go run main.go
#### If you encounter an error while parsing the environment file

> set -a && source env && set +a && go run main.go

### Deploy in Docker
> docker build -t kvservice . </br>
> docker run -p 8080:8080 kvservice

#### To run multiple nodes   
> docker compose -f 'docker-compose.yaml' up -d --build