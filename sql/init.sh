echo "STARTING COPY"

pushd /sql || exit
psql -f ./ddl/code.sql
popd || exit
