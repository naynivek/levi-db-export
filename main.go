package main

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/google/uuid"
	"os"
	"time"
	"github.com/naynivek/levi-db-export/getSnapshot"
	"github.com/naynivek/levi-db-export/exportSnapshot"
	"github.com/naynivek/levi-db-export/copySnapshot"
	"github.com/slack-go/slack"
)

func main() {
//Set variables
	COPY := "false"
	S3_EXPORT := "false"
	BUCKET_NAME := os.Getenv("BUCKET_NAME")
	DB_NAME := os.Getenv("DB_NAME")
	EXPORT_JOB_NAME := os.Getenv("EXPORT_JOB_NAME")
	AWS_ROLE_ARN := os.Getenv("AWS_ROLE_ARN")
	AWS_KMS_KEY_ID := os.Getenv("AWS_KMS_KEY_ID")
	AWS_KMS_KEY_ID_DST := os.Getenv("AWS_KMS_KEY_ID_DST")
	AWS_REGION_DST := os.Getenv("AWS_REGION_DST")
	AWS_S3_PREFIX := os.Getenv("AWS_S3_PREFIX")
	CREDS := os.Getenv("CREDS")
	COPY = os.Getenv("COPY")
	S3_EXPORT = os.Getenv("S3_EXPORT")
	SLACK_NOTIFY := os.Getenv("SLACK_NOTIFY")
	SLACK_TOKEN := os.Getenv("SLACK_TOKEN")
	SLACK_CHANNEL_ID := os.Getenv("SLACK_CHANNEL_ID")
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
	if AWS_KMS_KEY_ID_DST == "" {
		log.Fatal("Missing AWS_KMS_KEY_ID_DST environment variable")
	}
	if CREDS == "" {
		log.Println("Missing CREDS environment variable, using default configuration only")
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
	if CREDS == "web" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithWebIdentityRoleCredentialOptions(func(options *stscreds.WebIdentityRoleOptions) {
			options.RoleSessionName = "levi-db-export-go"
		}))
		if err != nil {
			log.Fatal(err)
		}
		credentialSet := cfg.Credentials
		log.Println("Start the app using this web credential: ",credentialSet)
	}

// Configure RDS client
	rdsClient := rds.NewFromConfig(cfg)
// Get Snapshot info
	snapshotName, err := getSnapshot.GetSnapshot(rdsClient, DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Snapshot name: ", *snapshotName.DBSnapshotArn)
// Run the copy to another region
	if COPY == "true" {
		cfg2, err := config.LoadDefaultConfig(context.TODO(),
					config.WithRegion(AWS_REGION_DST))
		if err != nil {
		log.Fatal(err)
		}
		rdsClient2 := rds.NewFromConfig(cfg2)
		currentTime := time.Now()
		destinationSnapshotName := *snapshotName.DBInstanceIdentifier+"-snapshot-destination-"+currentTime.Format("01-02-2006")
		copyTask, err := copySnapshot.CopySnapshot(rdsClient2, *snapshotName.DBSnapshotArn, destinationSnapshotName, AWS_KMS_KEY_ID_DST)
		log.Println("Copy is enabled")
		log.Println("Snapshot ARN: ", *copyTask.DBSnapshot.DBSnapshotArn)
		// Run the Slack notification
		if SLACK_NOTIFY == "true" {
			log.Printf("Sending the notification message to channel %s", SLACK_CHANNEL_ID)
			api := slack.New(SLACK_TOKEN)
			message := "The snapshot "+*snapshotName.DBSnapshotIdentifier+" has been copied to "+AWS_REGION_DST+" : https://us-east-2.console.aws.amazon.com/rds/home?region=us-east-2#db-snapshot:engine=postgres;id="+*copyTask.DBSnapshot.DBSnapshotIdentifier
			channelID, timestamp, err := api.PostMessage(
				SLACK_CHANNEL_ID,
				slack.MsgOptionText(message, false),)
			if err != nil {
				log.Printf("%s\n", err)
				return
			}
			log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
		}
	} else {
		log.Println("Copy is disabled")
	}
// Run the export to s3 task
	if S3_EXPORT == "true" {
		exportTask, err := exportSnapshot.ExportSnapshot(rdsClient,*snapshotName.DBSnapshotArn,EXPORT_JOB_NAME,AWS_ROLE_ARN,AWS_KMS_KEY_ID,AWS_S3_PREFIX,BUCKET_NAME)
		if err != nil {
			log.Fatal(err)
			}
		log.Println("Export to s3 is enabled")
		log.Println("Task identifier: ", *exportTask.ExportTaskIdentifier)
		if SLACK_NOTIFY == "true" {
			log.Printf("Sending the notification message to channel %s", SLACK_CHANNEL_ID)
			api := slack.New(SLACK_TOKEN)
			message := "The database "+DB_NAME+" it has been exported to "+BUCKET_NAME+" : https://us-east-1.console.aws.amazon.com/rds/home?region=us-east-1#export-task:id="+*exportTask.ExportTaskIdentifier+";source-location=exports-in-s3"
			channelID, timestamp, err := api.PostMessage(
				SLACK_CHANNEL_ID,
				slack.MsgOptionText(message, false),)
			if err != nil {
				log.Printf("%s\n", err)
				return
			}
			log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
		}
	} else {
		log.Println("Export to s3 is disabled")
	}
}
// adicionar testIdentifier dinâmico e prefixo do s3
// *exportTask.SourceType
// *exportTask.Status
// *exportTask.FailureCause