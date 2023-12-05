package exportSnapshot

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

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