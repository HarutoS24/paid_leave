#!/bin/bash

cd /frontend
npm ci
npm run dev -- --host 0.0.0.0 --port 5173