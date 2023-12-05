package copySnapshot

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func CopySnapshot(rdsClient *rds.Client, SourceSnapshotArn,TargetSnapshotName,src_region,dst_region string) (*rds.CopyDBSnapshotOutput, error) {
	output, err := rdsClient.CopyDBSnapshot(context.TODO(),
		&rds.CopyDBSnapshotInput{
			SourceRegion: aws.String(src_region),
			DestinationRegion: aws.String(dst_region),
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
