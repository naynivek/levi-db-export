package getInfo

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
