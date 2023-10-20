package getinfo

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func GetSnapshot(rdsClient *rds.Client, DBInstanceIdentifier string) (*types.DBSnapshot, error) {
	output, err := rdsClient.DescribeDBSnapshots(context.TODO(),
		&rds.DescribeDBSnapshotsInput{
			DBInstanceIdentifier: aws.String(DBInstanceIdentifier),
			SnapshotType:         aws.String("automated"),
			MaxRecords:           aws.Int32(20),
		})
	if err != nil {
		log.Printf("Couldn't get snapshot for instance %v: %v\n", DBInstanceIdentifier, err)
		return nil, err
	} else {
		index := len(output.DBSnapshots) - 1
		return &output.DBSnapshots[index], nil
	}
}

func ExportSnapshot(rdsClient *rds.Client, SourceArn,ExportTaskIdentifier,IamRoleArn,KmsKeyId,AWS_S3_PREFIX,BUCKET_NAME string) (*rds.StartExportTaskOutput, error) {
	output, err := rdsClient.StartExportTask(context.TODO(),
		&rds.StartExportTaskInput{
			SourceArn:			  aws.String(SourceArn),
			ExportTaskIdentifier: aws.String(ExportTaskIdentifier),
			IamRoleArn:           aws.String(IamRoleArn),
			KmsKeyId:             aws.String(KmsKeyId),
			S3Prefix:			  aws.String(AWS_S3_PREFIX),
			S3BucketName:		  aws.String(BUCKET_NAME),
		})
	if err != nil {
		log.Printf("Couldn't create the task %v: %v\n", SourceArn, err)
		return nil, err
	} else {
		return output, nil
	}
}