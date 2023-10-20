- [Description](#description)
- [Cloud Support](#cloud-support)
- [Getting Started](#getting-started)
  - [How to install](#how-to-install)
  - [Prerequisites](#prerequisites)
    - [Environment Variables](#environment-variables)
    - [Enable executable permission](#enable-executable-permission)
    - [Enable binary into your default PATH](#enable-binary-into-your-default-path)
  - [Simple run the app or decalre the ENV\_VARS during the execution](#simple-run-the-app-or-decalre-the-env_vars-during-the-execution)


# Description
The app levi-db-export aims to facilitade automation when you need to export database to a bucket, with the correct environment variables
is possible to start a job that will export database to a specific bucket

# Cloud Support
This app only supports database export on AWS which is most widely used Cloud Platform

# Getting Started

## How to install
Install the app with this command, pay attention into your architecture first
```bash
curl -sL -o levi-db-export "https://github.com/naynivek/levi-db-export/releases/latest/download/levi-db-export-macos-arm64"
```
You can check more releases on the [releases link](https://github.com/naynivek/levi-db-export/releases/latest/)


## Prerequisites
### Environment Variables
Before run the app you need to have this environment variables into your terminal session.
```bash
export BUCKET_NAME="YOUR_BUCKET_NAME"
export DB_NAME="YOUR_DATABASE_NAME"
export AWS_ROLE_ARN="YOUR_AWS_ROLE_ARN"
export AWS_KMS_KEY_ID="YOUR KMS KEY"
```
### Enable executable permission
Enable the executable permission in the binary
```bash
chmod +x levi-db-export
```
### Enable binary into your default PATH
Move the binary to your default PATH 
```bash
mv levi-db-export /usr/bin/levi-db-export
```
## Simple run the app or decalre the ENV_VARS during the execution
```bash
export BUCKET_NAME="YOUR_BUCKET_NAME" \
export DB_NAME="YOUR_DATABASE_NAME" \
export AWS_ROLE_ARN="YOUR_AWS_ROLE_ARN" \
export AWS_KMS_KEY_ID="YOUR KMS KEY" \
levi-db-export
```