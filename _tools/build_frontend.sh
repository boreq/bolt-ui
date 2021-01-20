#!/bin/bash
set -e

# Build frontend
echo "Running yarn build"
cd ../velo-frontend
rm -rf dist
yarn build

# Build backend
cd ../velo/ports/http/frontend
echo "Running statik"
statik -f -src=../../../../velo-frontend/dist
