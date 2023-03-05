package permanentdata

import (
	"context"
	"data-go-analyser/data"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/xitongsys/parquet-go-source/s3"
	"github.com/xitongsys/parquet-go/writer"
	"github.com/ztrue/tracerr"
)

func SaveProcessedData(datalist []data.ProcessedData, bucketName string, key string) error {
	var err error
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	s3Client := awss3.New(sess)

	fw, err := s3.NewS3FileWriterWithClient(context.Background(), s3Client, bucketName, key, "bucket-owner-full-control", nil)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return err
	}
	// create new parquet file writer
	pw, err := writer.NewParquetWriter(fw, new(data.ProcessedData), 4)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return err
	}
	// write 100 student records to the parquet file
	for _, data := range datalist {
		if err = pw.Write(data); err != nil {
			log.Println("Write error", err)
		}
	}
	// write parquet file footer
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop err", err)
	}

	err = fw.Close()
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	log.Println("Write Finished")
	return nil
}
