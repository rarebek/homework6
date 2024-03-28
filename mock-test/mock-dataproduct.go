package mocktest

import (
	pbp "EXAM3/api-gateway/genproto/product_service"
	"context"

	"google.golang.org/grpc"
)

type ProductServiceClientI interface {
	CreateProduct(ctx context.Context, in *pbp.Product, opts ...grpc.CallOption) (*pbp.Product, error)
	GetProductById(ctx context.Context, in *pbp.ProductId, opts ...grpc.CallOption) (*pbp.Product, error)
	UpdateProduct(ctx context.Context, in *pbp.Product, opts ...grpc.CallOption) (*pbp.Product, error)
	DeleteProduct(ctx context.Context, in *pbp.ProductId, opts ...grpc.CallOption) (*pbp.Status, error)
	ListProducts(ctx context.Context, in *pbp.GetAllProductRequest, opts ...grpc.CallOption) (*pbp.GetAllProductResponse, error)
}

type ProductServiceClient struct {
}

func NewProductServiceClient() *ProductServiceClient {
	return &ProductServiceClient{}
}

func (c *ProductServiceClient) CreateProduct(ctx context.Context, in *pbp.Product, opts ...grpc.CallOption) (*pbp.Product, error) {

	return in, nil
}

func (c *ProductServiceClient) GetProductById(ctx context.Context, in *pbp.ProductId, opts ...grpc.CallOption) (*pbp.Product, error) {

	return &pbp.Product{
		Id:          "1",
		Name:        "Nodirbek's Product",
		Description: "Nodirbek's Description",
		Price:       99.9,
		Amount:      99,
	}, nil
}

func (c *ProductServiceClient) UpdateProduct(ctx context.Context, in *pbp.Product, opts ...grpc.CallOption) (*pbp.Product, error) {
	return &pbp.Product{
		Id:          "yewruoe",
		Name:        "Product name",
		Description: "Product description",
		Price:       10.1,
		Amount:      5,
	}, nil
}

func (c *ProductServiceClient) DeleteProduct(ctx context.Context, in *pbp.ProductId, opts ...grpc.CallOption) (*pbp.Status, error) {
	return &pbp.Status{
		Success: true,
	}, nil
}

func (c *ProductServiceClient) ListProducts(ctx context.Context, in *pbp.GetAllProductRequest, opts ...grpc.CallOption) (*pbp.GetAllProductResponse, error) {
	pr := pbp.Product{
		Id:          "yewruoe",
		Name:        "Product name",
		Description: "Product description",
		Price:       10.1,
		Amount:      5,
	}
	return &pbp.GetAllProductResponse{
		Count: 3,
		Products: []*pbp.Product{
			&pr,
			&pr,
			&pr,
		},
	}, nil
}
