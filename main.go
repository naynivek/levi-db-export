package main

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/google/uuid"
	"os"
	"time"
	"github.com/naynivek/levi-db-export/getinfo"
)

func main() {
//Set variables
	BUCKET_NAME := os.Getenv("BUCKET_NAME")
	DB_NAME :=os.Getenv("DB_NAME")
	EXPORT_JOB_NAME :=os.Getenv("EXPORT_JOB_NAME")
	AWS_ROLE_ARN :=os.Getenv("AWS_ROLE_ARN")
	AWS_KMS_KEY_ID :=os.Getenv("AWS_KMS_KEY_ID")
	AWS_S3_PREFIX :=os.Getenv("AWS_S3_PREFIX")
//Verify variables
	if BUCKET_NAME == "" {
		log.Fatal("Missing BUCKET_NAME environment variable")
	}
	if DB_NAME == "" {
		log.Fatal("Missing DB_NAME environment variable")
	}
	if EXPORT_JOB_NAME == "" {
		log.Println("Missing EXPORT_JOB_NAME environment variable")
		log.Println("Creating a UUID for it")
		id := uuid.New()
		EXPORT_JOB_NAME=DB_NAME+"-"+id.String()
    	log.Println("The EXPORT_JOB_NAME will be: ",EXPORT_JOB_NAME)
	}
	if AWS_ROLE_ARN == "" {
		log.Fatal("Missing AWS_ROLE_ARN environment variable")
	}
	if AWS_KMS_KEY_ID == "" {
		log.Fatal("Missing AWS_KMS_KEY_ID environment variable")
	}
	if AWS_S3_PREFIX == "" {
		log.Println("Missing AWS_S3_PREFIX environment variable")
		log.Println("Creating a database and date for it")
		currentTime := time.Now()
		AWS_S3_PREFIX=DB_NAME+"-"+currentTime.Format("01-02-2006")
    	log.Println("The EXPORT_JOB_NAME will be: ",AWS_S3_PREFIX)
	}

// Configure SDK Client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	credentialSet := cfg.Credentials
	log.Println("Start the app using this credential: ",credentialSet)

// Configure rds client
	rdsClient := rds.NewFromConfig(cfg)
	snapshotName, err := getinfo.GetSnapshot(rdsClient, DB_NAME)
	log.Println("Snapshot name: ", *snapshotName.DBSnapshotArn)
	exportTask, err := getinfo.ExportSnapshot(rdsClient,*snapshotName.DBSnapshotArn,EXPORT_JOB_NAME,AWS_ROLE_ARN,AWS_KMS_KEY_ID,AWS_S3_PREFIX,BUCKET_NAME)
	log.Println("Task name: ", *exportTask)
}
// adicionar testIdentifier dinâmico e prefixo do s3
// *exportTask.SourceType
// *exportTask.Status
// *exportTask.FailureCause