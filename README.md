# GCP Key Rotator

## What?
This is a tool which can replace google cloud authentication json's pretty damn fast.

## How?
 - Make sure you are authenticated with GCP locally
 - `go run main.go <path/to/file.json>`
 - Tool will generate a new key, remove the old key and write results to original json file
