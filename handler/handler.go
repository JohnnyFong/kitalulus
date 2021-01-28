package handler

import (
	"context"
	"fmt"
	"kitalulus/helper"
	"kitalulus/model"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler struct
type Handler struct {
	R helper.Response
	C *mongo.Client
}

// CheckHeader func is to check the user header and return MyzoneUser interface{}
func (h Handler) CheckHeader(response http.ResponseWriter, request *http.Request) model.User {
	var u model.User
	response.Header().Set("content-type", "application/json")
	rh := request.Header.Get("Authorization")
	if rh == "" {
		h.R.StandardResponseNoMessage(response, 8)
		return u
	}
	// get token from header
	t := strings.Split(rh, "Bearer ")
	apiToken := t[1]
	fmt.Println(apiToken)
	collection := h.C.Database("kitalulus").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"apiToken": apiToken}).Decode(&u)
	if err != nil {
		h.R.StandardResponseWithMessage(response, 8, err.Error())
		return u
	}
	return u
}
