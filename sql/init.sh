echo "STARTING COPY"

pushd /sql
psql -f ./ddl/code.sql
popd
