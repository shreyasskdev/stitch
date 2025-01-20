package server

import (
	"log"
	"os"

	"github.com/FulgurCode/stitch/internal/handler"
	"github.com/FulgurCode/stitch/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(port string) {
	var app = echo.New()

	// Remove trailing slash to avoid 404 errors
	app.Pre(middleware.RemoveTrailingSlash())

	app.Use(middleware.Recover())

	// Creating session
	var sessionSecret = os.Getenv("SESSION_SECRET")
	app.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))

	// add static files
	app.Static("/static", "assets")

	// echo routes
	app.GET("/", handler.Home)
	app.GET("/products", handler.Products)
	app.GET("/item/:productId", handler.Item)
	app.GET("/order/:productId", handler.OrderGet)
	app.POST("/order/:productId", handler.OrderPost)
	app.GET("/about", handler.About)
	app.GET("/collections", handler.CommingSoon)
	app.GET("/cart", handler.Cart)
	app.GET("/cart/order", handler.CartOrderGet)
	app.POST("/cart/order", handler.CartOrderPost)
	app.GET("add-to-cart/:productId", handler.AddToCart)
	app.DELETE("delete-cart/:productId", handler.DeleteFromCart)

	app.GET("/admin", utils.CheckLogin(handler.Admin))

	app.GET("/admin/login", handler.AdminLogin)
	app.POST("/admin/login", handler.AdminLoginPost)
	app.GET("/admin/logout", utils.CheckLogin(handler.AdminLogout))
	app.POST("/admin/change-password", utils.CheckLogin(handler.AdminChangePasswordPost))
	app.GET("/admin/change-password", utils.CheckLogin(handler.AdminChangePassword))

	app.GET("/admin/products", utils.CheckLogin(handler.AdminProducts))
	app.GET("/admin/add-product", utils.CheckLogin(handler.AddProductGet))
	app.POST("/admin/product", utils.CheckLogin(handler.AddProductPost))
	app.DELETE("/admin/product/:productId", utils.CheckLogin(handler.DeleteProduct))
	app.PUT("/admin/product/:productId", utils.CheckLogin(handler.EditProduct))

	app.GET("/admin/stock", utils.CheckLogin(handler.AdminStock))
	app.POST("/admin/update-stock", utils.CheckLogin(handler.AdminUpdateStockPost))

	app.GET("/admin/item/:productId", utils.CheckLogin(handler.AdminItem))
	app.GET("/admin/orders", utils.CheckLogin(handler.AdminOrders))
	app.GET("/admin/shipped", utils.CheckLogin(handler.AdminShipped))
	app.GET("/admin/delivered", utils.CheckLogin(handler.AdminDelivered))
	app.GET("/admin/change-order-status/:orderId", utils.CheckLogin(handler.AdminChangeOrderStatus))
	app.DELETE("/admin/delete-order/:orderId", utils.CheckLogin(handler.DeleteOrder))
	app.GET("/admin/settings", utils.CheckLogin(handler.AdminSettings))

	app.POST("/admin/home-main", utils.CheckLogin(handler.AdminHomeBanner))
	app.POST("/admin/hero-one", utils.CheckLogin(handler.AdminHeroOne))
	app.POST("/admin/hero-two", utils.CheckLogin(handler.AdminHeroTwo))

	app.GET("/*", handler.NotFound)

	log.Fatal(app.Start(":" + port))
}
