# generate
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative modelopr.proto

# Presto
./presto --server localhost:9000 --execute "
CREATE SCHEMA IF NOT EXISTS minio.dadatalake
WITH (location = 's3a://dadatalake/');

CREATE TABLE IF NOT EXISTS minio.dadatalake.iris_parquet (
  sepal_length DOUBLE,
  sepal_width  DOUBLE,
  petal_length DOUBLE,
  petal_width  DOUBLE,
  class        VARCHAR
)
WITH (
  external_location = 's3a://dadatalake/iris_parquet',
  format = 'PARQUET'
);"