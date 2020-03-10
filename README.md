# bookstore_oauth-api
OAuth API


```
docker run --name cassandra-1 -p 7000:7000 -p 9042:9042 cassandra:latest
docker exec -it cassandra-1 bash
cqlsh
CREATE KEYSPACE oauth WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor': 1};
CREATE TABLE access_tokens( access_token varchar PRIMARY KEY, user_id bigint, client_id bigint, expires bigint);
```