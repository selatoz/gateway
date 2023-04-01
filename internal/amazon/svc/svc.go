package amazonSvc

import (
	"context"
	"fmt"
	"github.com/amazon/aws-sdk-go/aws"
	"github.com/amazon/aws-sdk-go/aws/credentials"
	"github.com/amazon/aws-sdk-go/aws/session"
	"github.com/amazon/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	"selatoz/pkg/cfglib"
)

// Define constants
const (
	// Errors
	ErrInvalidPassword 			= "invalid password"
	ErrFailedToHashPassword 	= "failed to hash password"
	ErrFailedToCreateUser 		= "failed to create user"
	ErrFailedToGenerateToken 	= "failed to generate token"
)

// Seller type defines the structure of information for a seller
type Seller struct {
	Name  string
	Stock int
}

type ProductCost struct {
	FBA	float64
	FBM	float64
}

type ProductRank struct {
	BSARank	int
}

// ProductInfo type defines the structure of information for a product
type ProductInfo struct {
	Cost			ProductCost
	Rank      	ProductRank
	Sellers   	[]Seller
	Stock     	int
}

// Svc is an interface for defining the methods that the user service will provide.
type Svc interface {
	GetProductInfo(ctx context.Context, asin string) (*ProductInfo, error)
	// Add more methods here as needed
}

// svc is an implementation of the UserSvc interface that handles the business logic for user-related operations.
type svc struct {
	session *session.Session
}

// NewSvc creates a new instance of svc and returns it as a Svc interface.
func NewSvc() (*AmazonService, error) {
	// Set up AWS credentials
	creds := credentials.NewStaticCredentials(cfglib.DefaultConf.AwsAccessKeyId, cfglib.DefaultConf.AwsSecretAccessKey, "")

	// Set up AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv(cfglib.DefaultConf.AwsRegion)),
		Credentials: creds,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %w", err)
	}

	return &AmazonService{
		session: sess,
	}, nil
}

func (s *svc) GetProductInfo(ctx context.Context, asin string) (*ProductInfo, error) {
	// Create a new Amazon MWS client
	client := mws.New(s.session)

	// Call the appropriate API methods to get the product information you need
	// You can use the client to call other Amazon MWS API methods as well
	// Here's an example of how to get the cost to sell on FBA:
	feesEstimateRequest := mws.FeesEstimateRequest{
		FeesEstimateRequestList: mws.FeesEstimateRequestList{
			FeesEstimateRequest: []mws.FeesEstimateRequest{
				{
					MarketplaceId: 		aws.String(cfglib.DefaultConf.AmazonMarketplaceId),
					IdType:        		aws.String("ASIN"),
					IdValue:       		aws.String(asin),
					IsAmazonFulfilled: 	aws.Bool(true),
				},
			},
		},
	}
	feesEstimateResponse, err := client.GetFeesEstimate(ctx, &feesEstimateRequest)
	if err != nil {
		return nil, fmt.Errorf("error getting FBA fees estimate: %w", err)
	}
	// Parse the response and extract the relevant information
	// You can do the same for the other pieces of information you need

	productInfo := &Product{
		CostToFBA: // extract the cost to sell on FBA from the response,
		CostToFBM: // extract the cost to sell on FBM from the response,
		Rank:      // extract the product ranking from the response,
		Sellers:   // extract the sellers selling the product from the response,
		Stock:     // extract the stock amount from the response,
	}

	return productInfo, nil
}

func getProductCosts(ctx context.Context) (float64, float64, error) {
		// Call the appropriate API methods to get the product information you need
	// You can use the client to call other Amazon MWS API methods as well
	// Here's an example of how to get the product information you need:
	getMyPriceForSKURequest := mws.GetMyPriceForSKUInput{
		SellerSKUList: &mws.SellerSKUListType{
			SellerSKU: []string{asin},
		},
		MarketplaceId: aws.String(os.Getenv("AMAZON_MARKETPLACE_ID")),
	}
	getMyPriceForSKUResponse, err := client.GetMyPriceForSKU(ctx, &getMyPriceForSKURequest)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting product information: %w", err)
	}
	// Parse the response and extract the relevant information
	// You can do the same for the other pieces of information you need
	var fbaCost, fbmCost float64
	for _, product := range getMyPriceForSKUResponse.GetMyPriceForSKUResult {
		if product.Product != nil {
			for _, offer := range product.Product.Offers {
				if aws.StringValue(offer.FulfillmentChannel) == "AMAZON" {
					fbaCost = aws.Float64Value(offer.ListingPrice.Amount)
				} else {
					fbmCost = aws.Float64Value(offer.ListingPrice.Amount)
				}
			}
		}
	}

	return fbaCost, fbmCost, nil
}