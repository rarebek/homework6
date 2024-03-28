package handlers1

import (
	"EXAM3/api-gateway/api_test/storage"
	pbp "EXAM3/api-gateway/genproto/product_service"
	pbu "EXAM3/api-gateway/genproto/user_service"
	mocktest "EXAM3/api-gateway/mock-test"
	"context"
	"github.com/gin-gonic/gin"
	_ "log"
	"net/http"
)

type Handler struct {
	UserService    *mocktest.UserServiceClient
	ProductService *mocktest.ProductServiceClient
}

func NewHandler(userService *mocktest.UserServiceClient, productService *mocktest.ProductServiceClient) *Handler {
	return &Handler{
		UserService:    userService,
		ProductService: productService,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser storage.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//
	//userJson, err := json.Marshal(newUser)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}

	//if err := kv.Set(newUser.Id, string(userJson), 1000); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}
	h.UserService.CreateUser(context.Background(), &pbu.User{
		Id:       newUser.Id,
		Name:     newUser.FirstName,
		Age:      newUser.Age,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	})
	c.JSON(http.StatusOK, newUser)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Query("id")
	//userString, err := kv.Get(userID)
	//pp.Println(userID)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//	return
	//}

	//var resp storage.User
	//if err := json.Unmarshal([]byte(userString), &resp); err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//	return
	//}
	us, err := h.UserService.GetUserById(context.Background(), &pbu.UserId{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, us)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userId := c.Query("id")
	h.UserService.DeleteUser(context.Background(), &pbu.UserId{UserId: userId})

	c.JSON(http.StatusOK, gin.H{
		"message": "user was deleted successfully",
	})
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.UserService.ListUser(context.Background(), &pbu.GetAllUserRequest{
		Page:  1,
		Limit: 99,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Product crud
func (h *Handler) CreateProduct(c *gin.Context) {
	var newProduct storage.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//userJson, err := json.Marshal(newProduct)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}
	//
	//if err := kv.Set(cast.ToString(newProduct.Id), string(userJson), 1000); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}

	h.ProductService.CreateProduct(context.Background(), &pbp.Product{
		Id:          newProduct.Id,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Amount:      newProduct.Amount,
	})
	c.JSON(http.StatusOK, newProduct)
}

func (h *Handler) GetProduct(c *gin.Context) {
	productID := c.Query("id")
	//productString, err := kv.Get(productID)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//var resp storage.Product
	//if err := json.Unmarshal([]byte(productString), &resp); err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//	return
	//}

	resp, err := h.ProductService.GetProductById(context.Background(), &pbp.ProductId{ProductId: productID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	productId := c.Query("id")
	//if err := kv.Delete(productId); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}

	h.ProductService.DeleteProduct(context.Background(), &pbp.ProductId{
		ProductId: productId,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "product was deleted successfully",
	})
}

func (h *Handler) ListProducts(c *gin.Context) {
	//productsStrings, err := kv.List()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}

	//var products []*storage.Product
	//for _, productString := range productsStrings {
	//	var product storage.Product
	//	if err := json.Unmarshal([]byte(productString), &product); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"error": err.Error(),
	//		})
	//		return
	//	}
	//	products = append(products, &product)
	//}
	products, err := h.ProductService.ListProducts(context.Background(), &pbp.GetAllProductRequest{
		Page:  1,
		Limit: 99,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}
