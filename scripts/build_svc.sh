set -e

cd `dirname "$0"`/../ # change to root of ducker project

rm -rf output # clean previous builds
mkdir -p output/bin output/conf
cp scripts/run_svc*.sh output/
cp -r config/ output/conf

cd cmd/backend/
go build -v -o ducker-svc
cd -
mv cmd/backend/ducker-svc output/bin/ducker-svc
