Testing Api's

AppInfo:

	curl -X GET  -H 'Content-Type: application/json' shopping-cart.io/v1/info

Register:

	curl -X POST -d '{"name":"nikhil", "username":"nik", "password":"nik", "email":"nik@gamil.com"}' -H 'Content-Type: application/json' shopping-cart.io/v1/register

	curl -X POST -d '{"name":"admin", "username":"admin", "password":"admin", "email":"admin@gamil.com"}' -H 'Content-Type: application/json' shopping-cart.io/v1/register

Login: 

	curl -X POST -d '{"username":"nik","password":"nik"}'  -H 'Content-Type: application/json' shopping-cart.io/v1/login

	curl -X POST -d '{"username":"admin","password":"admin"}'  -H 'Content-Type: application/json' shopping-cart.io/v1/login

Categories: create categories for shopping list

	curl -X POST -d '{"name": "sport"}' -H 'Content-Type: application/json' shopping-cart.io/v1/categories?access_token=31192eda-6776-4763-a2cb-46e7ec64ba38

	curl -X GET -H 'Content-Type: application/json' shopping-cart.io/v1/categories?access_token=31192eda-6776-4763-a2cb-46e7ec64ba38

Inventory: add items in inventory

	curl -X POST -d '{"name":"shoe","price": 1200, "quantity": 10}' -H 'Content-Type: application/json' shopping-cart.io/v1/inventory?access_token=31192eda-6776-4763-a2cb-46e7ec64ba38

	curl -X GET -H 'Content-Type: application/json' shopping-cart.io/v1/inventory?access_token=31192eda-6776-4763-a2cb-46e7ec64ba38


Shopping List: create category specific shopping lists

	curl -X GET -H 'Content-Type: application/json' shopping-cart.io/v1/cart?access_token=8604a30f-84e3-4ee9-b5e6-53f0552f80a7

	curl -X POST -d '{"name": "mylist", "category": "sport"}' -H 'Content-Type: application/json' shopping-cart.io/v1/cart?access_token=8604a30f-84e3-4ee9-b5e6-53f0552f80a7

	curl -X DELETE -H 'Content-Type: application/json' shopping-cart.io/v1/cart/5fcac58500d81ef89f7bc078?access_token=8604a30f-84e3-4ee9-b5e6-53f0552f80a7
 
List Item:  add / delete items in shopping list

	curl -X POST -d '{"id": "5fcac60300d81ef89f7bc090"}' -H 'Content-Type: application/json' shopping-cart.io/v1/cart/5fcac58500d81ef89f7bc078?access_token=8604a30f-84e3-4ee9-b5e6-53f0552f80a7

	curl -X DELETE -H 'Content-Type: application/json'  shopping-cart.io/v1/cart/5fcac58500d81ef89f7bc078/5fcac60300d81ef89f7bc090?access_token=8604a30f-84e3-4ee9-b5e6-53f0552f80a7
 