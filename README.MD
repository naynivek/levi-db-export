# Table of contents

- [Table of contents](#table-of-contents)
- [Description](#description)
- [Cloud Support](#cloud-support)
- [Getting Started](#getting-started)
  - [How to install](#how-to-install)
  - [Prerequisites](#prerequisites)
    - [Environment Variables](#environment-variables)
    - [Enable executable permission](#enable-executable-permission)
    - [Enable binary into your default PATH](#enable-binary-into-your-default-path)
  - [Authorization](#authorization)
  - [Simple run the app or declare the ENV\_VARS during the execution](#simple-run-the-app-or-declare-the-env_vars-during-the-execution)
    - [Check the s3 export task](#check-the-s3-export-task)
    - [Check the snapshot copy task](#check-the-snapshot-copy-task)


# Description
The app levi-db-export aims to facilitade automation when you need to:
1. Export database to a bucket, with the correct environment variables
2. Copy your database snapshot to a different snapshot region

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
1. YOUR_BUCKET_NAME -> The bucket that will receive the database export (For s3 export)
2. DB_NAME -> Your database name (For both snapshot copy and s3 export)
3. AWS_ROLE_ARN -> The ARN used to create the export task (For s3 export)
4. AWS_KMS_KEY_ID -> The KMS Key id to encrypt the database files on s3
5. AWS_KMS_KEY_ID_DST -> The KMS Key id to encrypt the snapshot copied (for snapshot copy, it should be the same AWS_KMS_KEY_ID with same regions or different if the destination snapshot will be stored in a different region. Please use an KMS Key from destination region)
6. AWS_REGION_DST -> The destination region in case of snapshot copied will be stored on a different region (for snapshot copy)
7. CREDS -> for using WebIdentityToken as an authentication method (optional)
8. COPY -> for enable the snapshot copy feature
9. S3_EXPORT -> for enable the s3 export feature
Here's the basic command
```bash
export BUCKET_NAME="YOUR_BUCKET_NAME"
export DB_NAME="YOUR_DATABASE_NAME"
export AWS_ROLE_ARN="YOUR_AWS_ROLE_ARN"
export AWS_KMS_KEY_ID="YOUR_KMS_KEY_ID"
export AWS_KMS_KEY_ID_DST="YOUR_KMS_KEY_ID_DST"
export COPY="true"
export S3_EXPORT="true"

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
## Authorization
The executor of this tool should have some AWS permissions in order to accomplish the task
Here is a sample policy statement:
```json
{
	"Statement": [
		{
			"Action": [
				"rds:DescribeDBSnapshots",
				"rds:DescribeDBSnapshotAttributes",
				"rds:StartExportTask",
				"rds:CopyDBSnapshot"
			],
			"Effect": "Allow",
			"Resource": [
				"*"
			],
			"Sid": "WriteOnRDS"
		},
		{
			"Action": [
				"iam:PassRole"
			],
			"Effect": "Allow",
			"Resource": [
				"YOUR_AWS_ROLE_ARN"
			],
			"Sid": "PassRole"
		},
		{
			"Action": [
				"kms:*"
			],
			"Effect": "Allow",
			"Resource": [
				"YOUR_KMS_KEY_ID",
				"YOUR_KMS_KEY_ID_DST"
			],
			"Sid": "kmsReadWrite"
		}
	],
	"Version": "2012-10-17"
}
```
## Simple run the app or declare the ENV_VARS during the execution
```bash
export BUCKET_NAME="YOUR_BUCKET_NAME" \
export DB_NAME="YOUR_DATABASE_NAME" \
export AWS_ROLE_ARN="YOUR_AWS_ROLE_ARN" \
export AWS_KMS_KEY_ID="YOUR_KMS_KEY" \
export AWS_KMS_KEY_ID_DST="YOUR_KMS_KEY_ID" \
export COPY="true" \
export S3_EXPORT="true" \
levi-db-export
```
### Check the s3 export task
If you're logged on AWS, you can check your job here, remember to change your region if needed:
https://us-east-1.console.aws.amazon.com/rds/home?region=us-east-1#exports-in-s3:

### Check the snapshot copy task
If you're logged on AWS, you can check your job here, remember to change your region if needed:
https://us-east-1.console.aws.amazon.com/rds/home?region=us-east-1#snapshots-list: