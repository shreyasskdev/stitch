package handler

import (
	"fmt"
	"net/http"

	"github.com/FulgurCode/stitch/models"
	"github.com/FulgurCode/stitch/pkg/mysql"
	"github.com/FulgurCode/stitch/utils"
	"github.com/FulgurCode/stitch/view/admin"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Admin handler
func Admin(c echo.Context) error {
	var component = admin.Admin()

	return utils.Render(c, component)
}

// Admin Login Handler
func AdminLogin(c echo.Context) error {
	var component = admin.Login()
	var loggedIn = utils.GetSessionValue(c, "auth", "isLoggedIn")
	if loggedIn != nil && loggedIn.(bool) {
		return c.Redirect(http.StatusSeeOther, "/admin")
	}

	return utils.Render(c, component)
}

// Admin Login request
func AdminLoginPost(c echo.Context) error {
	var body models.Admin
	var err = c.Bind(&body)
	if err != nil {
		fmt.Println(err)
	}

	user, err := mysql.GetAdminUser(body.Username)
	var correct = utils.BcryptCompare(body.Password, user.Password)

	if !correct {
		return utils.Render(c, admin.Login())
	}

	utils.CreateSession(c, "auth")
	utils.AddSessionValue(c, "auth", "username", user.Username)
	utils.AddSessionValue(c, "auth", "isLoggedIn", true)

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin")
}

// Admin Logout
func AdminLogout(c echo.Context) error {
	utils.DeleteSession(c, "auth")

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin/login")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/login")
}

func AdminChangePassword(c echo.Context) error {
	var component = admin.AdminChangePassword()

	return utils.Render(c, component)
}

func AdminChangePasswordPost(c echo.Context) error {
	var body models.AdminPass
	var err = c.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusSeeOther, "/admin")
	}

	body.Username = utils.GetSessionValue(c, "auth", "username").(string)

	a, err := mysql.GetAdminUser(body.Username)
	if err != nil {
		fmt.Println(err)
		return c.Redirect(http.StatusSeeOther, "/admin")
	}

	if !utils.BcryptCompare(body.OldPassword, a.Password) {
		var component = admin.AdminChangePassword()

		return utils.Render(c, component)
	}

	a.Password = utils.BcryptGenerateHash(body.NewPassword)

	err = mysql.UpdateAdminPassword(a)
	if err != nil {
		fmt.Println(err)
	}

	return c.Redirect(http.StatusSeeOther, "/admin")
}

// Admin Products Handler
func AdminProducts(c echo.Context) error {
	var products, err = mysql.GetProducts()
	if err != nil {
		var component = admin.AdminProducts([]models.Product{})

		return utils.Render(c, component)
	}

	var component = admin.AdminProducts(products)

	return utils.Render(c, component)
}

// Add Product GET
func AddProductGet(c echo.Context) error {
	var component = admin.AdminAddProduct()

	return utils.Render(c, component)
}

func DeleteProduct(c echo.Context) error {
	var id = c.Param("productId")

	var err = mysql.DeleteProduct(id)
	if err != nil {
		fmt.Print(err.Error())
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.NoContent(200)
	}

	var products, _ = mysql.GetProducts()

	var component = admin.AdminProducts(products)
	return utils.Render(c, component)
}

// Add Product POST
func AddProductPost(c echo.Context) error {
	var product models.Product
	c.Bind(&product)

	product.Id = uuid.NewString()

	form, _ := c.MultipartForm()

	var file = form.File["main-image"][0]
	utils.StoreFile(file, product.Id+"-main")

	var files = form.File["images"]
	utils.StoreFiles(files, product.Id)

	var err = mysql.AddProduct(product)
	if err != nil {
		fmt.Println(err)
	}

	var stock models.Stock = models.Stock{
		ProductId: product.Id,
		S:         0,
		M:         0,
		L:         0,
		XL:        0,
		XXL:       0,
		XXXL:      0,
	}

	err = mysql.AddStock(stock)
	if err != nil {
		fmt.Println(err)
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin/products")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/products")
}

func EditProduct(c echo.Context) error {
	var productId = c.Param("productId")
	var product models.Product

	c.Bind(&product)
	product.Id = productId

	form, _ := c.MultipartForm()

	if len(form.File["main-image"]) > 0 {
		var file = form.File["main-image"][0]

		var err = utils.StoreFile(file, product.Id+"-main")
		if err != nil {
			fmt.Println(err)
		}
	}

	var files = form.File["images"]
	utils.StoreFiles(files, product.Id)

	var err = mysql.EditProduct(product)
	if err != nil {
		fmt.Println(err)

		var product, _ = mysql.GetProductById(productId)

		var component = admin.AdminItem(product)
		return utils.Render(c, component)
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin/products")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/products")
}

func AdminUpdateStockPost(c echo.Context) error {
	var stock models.Stock
	c.Bind(&stock)

	var err = mysql.UpdateStock(stock)
	if err != nil {
		fmt.Println(err)
	}
	stocks, err := mysql.GetStocks()
	if err != nil {
		fmt.Println(err.Error())
	}

	var component = admin.AdminStock(stocks)

	return utils.Render(c, component)
}

// Admin Item Handler
func AdminItem(c echo.Context) error {
	var id = c.Param("productId")

	var product, err = mysql.GetProductById(id)
	if err != nil {
		fmt.Print(err.Error())
	}

	var component = admin.AdminItem(product)
	return utils.Render(c, component)
}

// Admin Orders Handler
func AdminOrders(c echo.Context) error {
	var orders, err = mysql.GetOrdersWithStatus("ordered")
	if err != nil {
		fmt.Println(err.Error())
	}

	var component = admin.AdminOrders(orders)
	return utils.Render(c, component)
}

func AdminChangeOrderStatus(c echo.Context) error {
	var status = c.QueryParam("status")
	var productId = c.QueryParam("productId")
	var id = c.Param("orderId")
	var err = mysql.ChangeOrderStatus(id, productId, status)
	if err != nil {
		fmt.Println(err)
		return c.JSON(501, "Failed")
	}

	return c.NoContent(200)
}

// Admin Stock Handler
func AdminStock(c echo.Context) error {
	var stocks, err = mysql.GetStocks()
	if err != nil {
		fmt.Println(err.Error())
	}
	var component = admin.AdminStock(stocks)

	return utils.Render(c, component)
}

// Admin Shipped Handler
func AdminShipped(c echo.Context) error {
	var orders, err = mysql.GetOrdersWithStatus("shipped")
	if err != nil {
		fmt.Println(err.Error())
	}

	var component = admin.AdminShipped(orders)
	return utils.Render(c, component)
}

// Admin Delivered Handler
func AdminDelivered(c echo.Context) error {
	var orders, err = mysql.GetOrdersWithStatus("delivered")
	if err != nil {
		fmt.Println(err.Error())
	}

	var component = admin.AdminDelivered(orders)
	return utils.Render(c, component)
}

func DeleteOrder(c echo.Context) error {
	var id = c.Param("orderId")
	var productId = c.QueryParam("productId")
	var err = mysql.DeleteOrder(id, productId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(501, "Failed")
	}

	return c.NoContent(200)
}

// Admin Setting Handler
func AdminSettings(c echo.Context) error {
	var settings = mysql.GetSettings()
	var component = admin.AdminSettings(settings)

	return utils.Render(c, component)
}

func AdminHomeBanner(c echo.Context) error {
	var form, _ = c.MultipartForm()

	if img := form.File["home_main_banner"]; len(img) > 0 {
		utils.StoreFile(img[0], "home_main_banner")
	}

	var homeDescription = form.Value["home_main_description"][0]
	var homeTitle = form.Value["home_main_title"][0]
	mysql.UpdateHome(homeTitle, homeDescription)

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin/settings")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}

func AdminHeroOne(c echo.Context) error {
	form, _ := c.MultipartForm()

	if img := form.File["hero_one_banner"]; len(img) > 0 {
		utils.StoreFile(img[0], "hero_one_banner")
	}

	var heroDescription = form.Value["hero_one_description"][0]
	var heroTitle = form.Value["hero_one_title"][0]
	mysql.UpdateHeroOne(heroTitle, heroDescription)

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin/settings")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}

func AdminHeroTwo(c echo.Context) error {
	form, _ := c.MultipartForm()

	if img := form.File["hero_two_banner"]; len(img) > 0 {
		utils.StoreFile(img[0], "hero_two_banner")
	}

	var heroDescription = form.Value["hero_two_description"][0]
	var heroTitle = form.Value["hero_two_title"][0]
	mysql.UpdateHeroTwo(heroTitle, heroDescription)

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Location", "/admin/settings")
		return c.NoContent(200)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}
