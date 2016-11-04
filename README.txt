

docker-compose  up gobin

to test all (including tarantool) with existing binary



docker-compose  up go_cut

to test with custom binary from source code


adjust envs PARALLEL and ITERATIONS to change workload in docker-compose.yml


docker-compose exec tnt console

open tarantool's console