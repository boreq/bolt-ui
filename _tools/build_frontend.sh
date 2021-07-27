#!/bin/bash
set -e

if ! git diff-index --quiet HEAD --; then
    echo "You have uncommited changes in your repository!"
    exit 1
fi

versionFile="ports/http/frontend/version.go"
previousCommit=$(cat ${versionFile} | grep FrontendCommit | awk ' {gsub(/"/, "", $4); print $4 }')

cd ./frontend

commit=$(git rev-parse HEAD)
echo "Previous frontend commit: ${previousCommit}"
echo "Current frontend commit: ${commit}"

if [ "$commit" = "$previousCommit" ]; then
    echo "Frontend is already up to date (${commit} == ${previousCommit})"
    exit 0
fi

echo "Building..."
rm -rf dist
yarn build

echo "Copying files..."
cd ../ports/http/frontend
cp -r ../../../frontend/dist/. ./
cd ../../../

echo "Persisting frontend version..."
echo "package frontend" > "${versionFile}"
echo "" >> "${versionFile}"
echo "const FrontendCommit = \"${commit}\"" >> "${versionFile}"

echo "Commiting..."
commitFile="/tmp/velo-frontend-commit.txt"
echo "Update frontend" > ${commitFile}
git add ports/http/frontend/
git commit -F ${commitFile}
