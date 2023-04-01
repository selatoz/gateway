package amazonSvc

import (
	"context"
	"fmt"
	"net/http"
	"os"

	// "gopkg.me/selling-partner-api-sdk"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"selatoz/gateway/pkg/cfglib"
)

// Define constants
const (
	// Errors
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
func NewSvc() (Svc, error) {
	// Set up AWS credentials

	// Set up AWS session

	// Handle return
	return nil, nil
}

func (s *svc) GetProductInfo(ctx context.Context, asin string) (*ProductInfo, error) {
	// Create a new Amazon MWS client

	// Call the appropriate API methods to get the product information you need

	// Parse the response and extract the relevant information
	// You can do the same for the other pieces of information you need

	return ProductInfo{}, nil
}

func getProductCosts(ctx context.Context) (float64, float64, error) {
	// Call the appropriate API methods to get the product information you need
	// You can use the client to call other Amazon MWS API methods as well
	// Here's an example of how to get the product information you need:

	// Parse the response and extract the relevant information
	// You can do the same for the other pieces of information you need
	var fbaCost, fbmCost float64

	// Handle return
	return fbaCost, fbmCost, nil
}