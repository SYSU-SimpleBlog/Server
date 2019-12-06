/*
 * simple blog
 *
 * A Simple Blog
 *
 * API version: 1.0.0
 * Contact: apiteam@swagger.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"log"
	"net/http"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
	sw "github.com/Server/go"
)

func main() {
	log.Printf("Server started")
	//t.TestArticle()
	router := sw.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

}
