package Reader

import (
	"bufio"
	"context"
	"data-go-analyser/data"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ztrue/tracerr"
)

func DownloadFile(s3Client *s3.Client, bucketName string, objectKey string) (io.Reader, error) {
	getObjectOutput, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
	return getObjectOutput.Body, err
}

func ReadFile(readFile io.Reader) (customers []*data.Customer, orders []*data.Order, items []*data.Item) {
	fileScanner := bufio.NewScanner(readFile)
	flag := 0

	for fileScanner.Scan() {
		if fileScanner.Text() == "" {
			continue
		}
		line := strings.Split(fileScanner.Text(), ",")
		if line[0] != "id" {
			if flag == 1 {
				customers = append(customers, data.ParseCustomer(line))
				continue
			}
			if flag == 2 {
				orders = append(orders, data.ParseOrders(line))
				continue
			}
			if flag == 3 {
				items = append(items, data.ParseItems(line))
				continue
			}
			continue
		}
		flag = flag + 1
	}
	return customers, orders, items
}

func ConstructData(customers []*data.Customer, orders []*data.Order, items []*data.Item) ([]*data.Customer, []*data.ErrorMessage) {
	// link Items to Orders
	for _, item := range items {
		for _, order := range orders {
			if order.OrderReference == item.OrderReference {
				order.Items = append(order.Items, item)
			}
		}
	}
	// link Orders to Customers
	for _, order := range orders {
		for _, customer := range customers {
			if customer.Orders == nil {
				customer.Orders = make(map[string]*data.Order)
			}
			if customer.CustomerReference == order.CustomerReference {
				customer.Orders[order.OrderReference] = order
			}
		}
	}
	return customers, nil
}

func AnalyseData(customers []*data.Customer) (processedData []data.ProcessedData) {
	for customer_reference, customer := range customers {
		processedData = append(processedData, data.ProcessedData{
			Type:              "customer_message",
			CustomerReference: customer.CustomerReference,
			NumberOfOrders:    int32(len(customer.Orders)),
		})
		for _, order := range customer.Orders {
			for _, item := range order.Items {
				processedData[customer_reference].TotalAmountSpent = processedData[customer_reference].TotalAmountSpent + item.TotalPrice
			}
		}
	}
	return processedData
}
