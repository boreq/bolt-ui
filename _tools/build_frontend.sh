#!/bin/bash
set -e

if ! git diff-index --quiet HEAD --; then
    echo "You have uncommited changes in your repository!"
    exit 1
fi

cd ./frontend

echo "Building..."
rm -rf dist
yarn build

echo "Copying files..."
cd ../ports/http/frontend
cp -r ../../../frontend/dist/. ./
cd ../../../

echo "Commiting..."
commitFile="/tmp/velo-frontend-commit.txt"
echo "Update frontend" > ${commitFile}
git add ports/http/frontend/
git commit -F ${commitFile}
