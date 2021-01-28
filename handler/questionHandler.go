package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"kitalulus/model"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/gobeam/mongo-go-pagination"
	guuid "github.com/google/uuid"
	"github.com/gorilla/mux"
)

// QuestionCreate - to create the question document
func (h Handler) QuestionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// get user from header
		u := h.CheckHeader(w, r)

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var q model.Question

		err := dec.Decode(&q)
		q.CreatedAt = time.Now()
		q.UpdatedAt = time.Now()
		q.CreatedBy = u.ID.Hex()
		q.UpdatedBy = u.ID.Hex()
		q.IsActive = true
		q.Uuid = guuid.New()
		fmt.Println(q.Uuid)

		if err != nil {
			h.R.StandardResponseWithMessage(w, 99, err.Error())
			return
		}
		collection := h.C.Database("kitalulus").Collection("question")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := collection.InsertOne(ctx, q)
		if err != nil {
			h.R.StandardResponseWithMessage(w, 99, err.Error())
			return
		}
		h.R.StandardResponse(w, 0, result)
	}
}

// QuestionList - list question with pagination
func (h Handler) QuestionList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		collection := h.C.Database("kitalulus").Collection("question")

		filter := bson.M{"isActive": true}

		var limit int64 = 10
		var page int64 = 1

		l := r.FormValue("itemsPerPage")
		p := r.FormValue("page")

		if l != "" {
			i, err := strconv.ParseInt(l, 10, 64)
			if err != nil {
				h.R.StandardResponseWithMessage(w, 99, err.Error())
				return
			}
			limit = i
		}

		if p != "" {
			i, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				h.R.StandardResponseWithMessage(w, 99, err.Error())
				return
			}
			page = i
		}

		paginatedData, err := New(collection).Limit(limit).Page(page).Filter(filter).Find()
		if err != nil {
			if err != nil {
				h.R.StandardResponseWithMessage(w, 99, err.Error())
				return
			}
		}

		var lists []model.Question
		for _, raw := range paginatedData.Data {
			var q *model.Question
			if marshallErr := bson.Unmarshal(raw, &q); marshallErr == nil {
				lists = append(lists, *q)
			} else {
				h.R.StandardResponseWithMessage(w, 99, marshallErr.Error())
				return
			}

		}
		h.R.StandardResponsePagination(w, 0, lists, paginatedData.Pagination)
	}
}

// QuestionShow - show details using UUID
func (h Handler) QuestionShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var q model.Question

		uuid, _ := guuid.Parse(mux.Vars(r)["uuid"])
		data, err := uuid.MarshalBinary()

		if err != nil {
			h.R.StandardResponseWithMessage(w, 8, err.Error())
			return
		}

		collection := h.C.Database("kitalulus").Collection("question")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		errr := collection.FindOne(ctx, bson.M{"uuid": data, "isActive": true}).Decode(&q)
		if errr != nil {
			h.R.StandardResponseWithMessage(w, 8, errr.Error())
			return
		}

		h.R.StandardResponse(w, 0, q)
	}
}

// QuestionRemove - to delete from collection
func (h Handler) QuestionRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// get user from header
		u := h.CheckHeader(w, r)

		objID, errr := primitive.ObjectIDFromHex(r.FormValue("id"))
		if errr != nil {
			h.R.StandardResponseWithMessage(w, 8, errr.Error())
			return
		}
		collection := h.C.Database("kitalulus").Collection("question")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{
				"$set": bson.M{
					"isActive": false,
					"updateAt": time.Now(),
					"updateBy": u.ID.Hex(),
				},
			},
		)
		if err != nil {
			h.R.StandardResponseWithMessage(w, 8, err.Error())
			return
		}
		h.R.StandardResponse(w, 0, result.ModifiedCount)
	}
}

// QuestionUpdate - to delete from collection
func (h Handler) QuestionUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// get user from header
		u := h.CheckHeader(w, r)

		objID, errr := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
		if errr != nil {
			h.R.StandardResponseWithMessage(w, 8, errr.Error())
			return
		}
		collection := h.C.Database("kitalulus").Collection("question")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{
				"$set": bson.M{
					"question": r.FormValue("question"),
					"updateAt": time.Now(),
					"updateBy": u.ID.Hex(),
				},
			},
		)
		if err != nil {
			h.R.StandardResponseWithMessage(w, 8, err.Error())
			return
		}
		h.R.StandardResponse(w, 0, result.ModifiedCount)
	}
}
