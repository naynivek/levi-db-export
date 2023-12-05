package copySnapshot

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func CopySnapshot(rdsClient *rds.Client, SourceSnapshotArn,TargetSnapshotName string) (*rds.CopyDBSnapshotOutput, error) {
	output, err := rdsClient.CopyDBSnapshot(context.TODO(),
		&rds.CopyDBSnapshotInput{
			SourceDBSnapshotIdentifier: aws.String(SourceSnapshotArn),
			TargetDBSnapshotIdentifier: aws.String(TargetSnapshotName),
		})
	if err != nil {
		log.Printf("Couldn't create the task %v: %v\n", SourceSnapshotArn, err)
		return nil, err
	} else {
		return output, nil
	}
	}
