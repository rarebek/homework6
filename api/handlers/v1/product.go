package v1

import (
	"EXAM3/api-gateway/api/model"
	pb "EXAM3/api-gateway/genproto/product_service"
	"EXAM3/api-gateway/pkg/logger"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateProduct
// @Summary create product
// @Tags Product
// @Description Insert a new product with provided details
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ProductDetails body model.Item true "Create product"
// @Success 201 {object} model.Item
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/product/create [post]
func (h *handlerV1) CreateProduct(c *gin.Context) {
	var (
		body       model.Item
		jspMarshal protojson.MarshalOptions
	)
	jspMarshal.UseProtoNames = true

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("failed to bind json", logger.Error(err))
		return
	}

	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: "amount cannot be smaller than zero",
		})
		return
	}

	if body.Price < 0 {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: "price cannot be smaller than zero",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	id := uuid.New().String()
	resp, err := h.serviceManager.MockProductService().CreateProduct(ctx, &pb.Product{
		Id:          id,
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Amount:      int64(body.Amount),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("error while creating product", logger.Error(err))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Update Product
// @Summary update product
// @Tags Product
// @Description Update ptoduct
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param UserInfo body model.Item true "Update Product"
// @Success 201 {object} model.Item
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/product/update/{id} [put]
func (h *handlerV1) UpdateProduct(c *gin.Context) {
	var (
		body        pb.Product
		jspbMarshal protojson.MarshalOptions
	)
	id := c.Param("id")

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	if body.Amount < 0 {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: "0 dan kichik amount kiritib bo`lmaydi",
		})
		return
	}

	if body.Price < 0 {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: "price cannot be smaller than zero",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockProductService().UpdateProduct(ctx, &pb.Product{
		Id:          id,
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Amount:      body.Amount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("error while updating product", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get Product By Id
// @Summary get product by id
// @Tags Product
// @Description Get product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} model.Item
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/product/get/{id} [get]
func (h *handlerV1) GetProductById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockProductService().GetProductById(ctx, &pb.ProductId{
		ProductId: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot get product", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete Product
// @Summary delete product
// @Tags Product
// @Description Delete product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} model.Status
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/product/delete/{id} [delete]
func (h *handlerV1) DeleteProduct(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockProductService().DeleteProduct(ctx, &pb.ProductId{
		ProductId: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})

		h.log.Error("cannot delete product", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get All Products
// @Summary get all products
// @Tags Product
// @Description get all products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page path string true "page"
// @Param limit path string true "limit"
// @Success 201 {object} model.ListItems
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/product/{page}/{limit} [get]
func (h *handlerV1) ListProducts(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	page := c.Param("page")
	intpage, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot parse page query param", logger.Error(err))
		return
	}

	limit := c.Param("limit")
	intlimit, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot parse limit query param", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockProductService().ListProducts(ctx, &pb.GetAllProductRequest{
		Page:  int64(intpage),
		Limit: int64(intlimit),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})

		h.log.Error("cannot list products", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// // Get All Purchased products by user
// // @Summary get all purchased products by user id
// // @Tags Product
// // @Description get all purchased products by user id
// // @Security BearerAuth
// // @Accept json
// // @Produce json
// // @Param id path string true "id"
// // @Success 201 {object} model.BoughtItemsList
// // @Failure 400 string Error models.ResponseError
// // @Failure 500 string Error models.ResponseError
// // @Router /v1/products/get/{id} [get]
// func (h *handlerV1) GetPurchasedProductsByUserId(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	userId := c.Param("id")
// 	if userId == "" {
// 		fmt.Println("=======")
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.MockProductService().GetBoughtProductsByUserId(ctx, &pb.UserId{
// 		UserId: userId,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, model.ResponseError{
// 			Code:  ErrorCodeInternalServerError,
// 			Error: err.Error(),
// 		})

// 		h.log.Error("cannot list products purchased by user", logger.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// Buy product
// @Summary buy a product
// @Tags Product
// @Description buy a product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param PurchaseInfo body model.BuyItemRequest true "Purchase a product"
// @Success 201 {object} model.BuyItemResponse
// @Failure 400 string Error models.ResponseError
// @Failure 500 string Error models.ResponseError
// @Router /v1/product/buy [post]
// func (h *handlerV1) BuyProduct(c *gin.Context) {
// 	var (
// 		res         model.BuyItemResponse
// 		body        model.BuyItemRequest
// 		jspbMarshal protojson.MarshalOptions
// 	)

// 	jspbMarshal.UseProtoNames = true
// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, model.ResponseError{
// 			Code:  ErrorCodeInvalidJSON,
// 			Error: err.Error(),
// 		})
// 		h.log.Error("cannot bind json", logger.Error(err))
// 		return
// 	}

// 	if body.Amount < 0 {
// 		c.JSON(http.StatusBadRequest, model.ResponseError{
// 			Code:  ErrorBadRequest,
// 			Error: "0 dan kichik amount kiritib bo`lmaydi",
// 		})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	status, err := h.serviceManager.MockProductService().CheckAmount(ctx, &pb.ProductId{ProductId: body.ProductId})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, model.ResponseError{
// 			Code:  ErrorCodeInternalServerError,
// 			Error: err.Error(),
// 		})

// 		h.log.Error("cannot list products purchased by user", logger.Error(err))
// 		return
// 	}
// 	if status.Amount == 0 {
// 		c.JSON(http.StatusBadRequest, model.ResponseError{
// 			Code:  ErrorBadRequest,
// 			Error: "the product is not currently available, sorry",
// 		})

// 		h.log.Error("not available product", logger.Error(err))
// 		return

// 	}
// 	if !(status.Amount < body.Amount) {
// 		c.JSON(http.StatusBadRequest, model.ResponseError{
// 			Code:  ErrorCodeInternalServerError,
// 			Error: fmt.Sprintf("we have only %d, sorry", status.Amount),
// 		})

// 		h.log.Error("not enough", logger.Error(err))
// 		return
// 	}
// 	buyResp, err := h.serviceManager.MockProductService().BuyProduct(ctx, &pb.BuyProductRequest{
// 		UserId:    body.UserId,
// 		ProductId: body.ProductId,
// 		Amount:    body.Amount,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, model.ResponseError{
// 			Code:  ErrorCodeInternalServerError,
// 			Error: err.Error(),
// 		})

// 		h.log.Error("cannot purchase the product", logger.Error(err))
// 		return
// 	}
// 	_, err = h.serviceManager.MockProductService().DecreaseAmount(ctx, &pb.ProductAmountRequest{
// 		ProductId: body.ProductId,
// 		Amount:    body.Amount,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, model.ResponseError{
// 			Code:  ErrorCodeInternalServerError,
// 			Error: err.Error(),
// 		})
// 		h.log.Error("cannot decrease the amount", logger.Error(err))
// 		return
// 	}

// 	res.Message = "successfully purchased"
// 	res.ProductId = body.ProductId
// 	res.UserId = body.UserId
// 	res.Amount = body.Amount
// 	res.ProductName = buyResp.Name

// 	c.JSON(http.StatusOK, res)
// }
