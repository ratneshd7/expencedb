package apictrl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ratneshd7/expencedb/dboperations"
	"github.com/ratneshd7/expencedb/modal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Router Routing over Url
func Router() {
	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/api/", get).Methods("GET")
	r.HandleFunc("/api/{id}", getByID).Methods("GET")
	r.HandleFunc("/api/add", add).Methods("POST")
	r.HandleFunc("/api/update/{id}", update).Methods("PUT")
	r.HandleFunc("/api/delete/{id}", delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))
}

func get(w http.ResponseWriter, r *http.Request) {
	expense, err := dboperations.GetAll()
	fmt.Println("get endpoint")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	var expns []modal.ExpenseModal
	expense.All(context.TODO(), &expns)
	if expns == nil {
		w.Write([]byte(`{"message":"Empty Data"}`))
		return
	}
	json.NewEncoder(w).Encode(expns)
}

func getByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	expense, err := dboperations.GetByID(id["id"])
	fmt.Println("get endpoint")
	if err != nil {
		// json.NewEncoder(w).Encode("error Occured")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	var expns []modal.ExpenseModal
	expense.All(context.TODO(), &expns)
	if expns == nil {
		w.Write([]byte(`{"message": "Data Not Found"}`))
		return
	}
	json.NewEncoder(w).Encode(expns)
}

func add(w http.ResponseWriter, r *http.Request) {
	var expense modal.ExpenseModal
	_ = json.NewDecoder(r.Body).Decode(&expense)
	expense.ID = primitive.NewObjectID()
	exps, err := dboperations.Insert(expense)
	fmt.Println("add endpoint")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(exps.InsertedID)
}

func update(w http.ResponseWriter, r *http.Request) {
	var expense modal.ExpenseModal
	_ = json.NewDecoder(r.Body).Decode(&expense)
	expns, err := dboperations.Update(mux.Vars(r)["id"], expense)
	fmt.Println("update endpoint")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(expns.ModifiedCount)
}

func delete(w http.ResponseWriter, r *http.Request) {
	expense, err := dboperations.DeleteByID(mux.Vars(r)["id"])
	fmt.Println("delete endpoint")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(expense.DeletedCount)
}
