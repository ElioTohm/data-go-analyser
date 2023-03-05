# Process uploaded file on S3 upload

### Infrastructure
#
simple solution utilizing s3+lambda trigger
upon upload on the file it will trigger the lambda function which will process the data save it as parquet file in a seperate bucket and send the json formated svc as a message in sqs

### Listen to S3 events
#
the infra is setup to trigger a lambda function on upload

#### Process the information
#
the information is read from the s3 file transformed into a `struct`
### Send to queue
#
the processed data is than sent to the queue created for later processing

### Store the information
#
In a parallel thread the processed data is being saved on a permanent storage

# Scaling the Solution

### Data Structure
#
the client information should be moved to a Database and cached into a redis cluster for faster read

### Infrastructure modifications
#
Uploading directly to the S3 shouldn't be an option given to the client
we should create a microservice that will handle the upload.
We can keep sqs as a messaging queue or for faster live processing we can move to kafka however this decision should be based on how and how fast we want to process this data
we can also store the information as parquet files on S3 and use presto to query this data.