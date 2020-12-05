# Shopping-Cart: Microservice Demo Application

## Description:

shopping-cart offers REST api's to create category wise shopping list(s). 
user (admin) can define categories for the list, and can define items which can be included in user list.  
user can create create personalized shopping-list by selecting categories which are pre-defined and can add items to his lists.
user can view, modify or delete items added in list.

It is built using Go and packaged in Docker containers.
It can be deployed easily using helm chart included.

## Build:

update repository details in ./Makefile and use following command to generate build and push image to docker repo
 	
	command: make all

## Deployment:

Navigate to "shopping-cart/deploy/shopping-cart/" and use following command to deploy shopping-cart app

	deploy example: 
 		helm install shopingcart -n=cart ./ 
 	uninstall example:
        	helm delete shopingcart -n=cart

## Test 

For unit testing it uses [testify](https://github.com/stretchr/testify) library.

## API Guide

Please refer [Api Test Guide](API_Guide.md) includes all api call's with examples.
