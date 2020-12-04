Project: Shopping-Cart

Description 
shopping-cart is microservice includes following features
1. User Registration and Authentication
2. Add / Update / Delete inventory (enabled for admin only)
3. Add / Update / Delete / View Cart 

Endpoints 

	App Info :
	"GET" "/v1/info", user.AppDetails" 
	
	RegisterUser: 
	"POST" "/v1/register", user.RegisterUser).Methods("POST")
			
	Login:
	"POST" "/v1/login" 
	
	Logout:
	"POST" "/v1/logout" 
	
	Inventory: [admin only]
	 "GET" "/v1/inventory" 
	 "POST" "/v1/inventory" 
	 "DELETE" "/v1/inventory" 
	 "DELETE" "/v1/inventory{itemid} 
	
	Cart:
	 "GET" "/v1/cart"
	 "PATCH" "/v1/cart" 
	 "DELETE" "/v1/cart" 
	 "DELETE" "/v1/cart/{itemid}" 
	 
Build:
 update repository details in ./Makefile and use following command to generate build and push image to docker repo
 command: make all

Deployment:
 user can use helm chart to deploy shopping-cart app on cluster.
 example: 
    navigate to deploy/shopping-cart and use following command to deploy shopping-cart app
	helm install shopingcart -n=cart ./ 
			
