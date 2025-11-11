#!/bin/bash
set -e

echo "Starting the Check for Frontend..."

cd /frontend
npm ci
npm run lint

echo "Done Successfully!"