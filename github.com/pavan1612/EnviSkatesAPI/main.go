package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"io/ioutil"
	"strconv"
)


type User struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Product struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Image string `json:"image,omitempty"`
	Price string `json:"price,omitempty"`
	Carbon_Footprint string `json:"carbon_footprint,omitempty"`
	Desc string `json:"desc,omitempty"`
}

type Order struct {
	ProductID string `json:"productid,omitempty"`
	UserID string `json:"userid,omitempty"`
	Status string `json:"status,omitempty"`
}

type ProductArray struct {
	Products []Product `json:"products,omitempty"`
}

type carbonEmission struct {
	CarbonEmissionVal int64 `json:"c_emission"`
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/register", RegisterEndPoint).Methods("POST")
	router.HandleFunc("/login",LoginEndPoint).Methods("POST")
	router.HandleFunc("/getAllProducts",GetAllProducts).Methods("GET")
	router.HandleFunc("/createOrder",CreateOrder).Methods("POST")
	router.HandleFunc("/getAllOrders",GetAllOrders).Methods("GET")
	router.HandleFunc("/getCarbonFootprint",GetCarbonEmmision).Methods("POST")
	fmt.Println("Server started")	
	log.Fatal(http.ListenAndServe(":8000",router))
}

func RegisterEndPoint(w http.ResponseWriter, r *http.Request){

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@enviskates-zpjae.gcp.mongodb.net/admin?retryWrites=true&w=majority"))
	if err != nil{
		log.Fatal(err)
	}
	UserCollection := client.Database("DevEnvi").Collection("User")
	res, err := UserCollection.InsertOne(ctx, user)
	fmt.Println(res)
	fmt.Println(bson.M{"user":"Pavan"})
}
func LoginEndPoint(w http.ResponseWriter, r *http.Request){
	var user User
	var result User
	_ = json.NewDecoder(r.Body).Decode(&user)
	ctx:= context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@enviskates-zpjae.gcp.mongodb.net/admin?retryWrites=true&w=majority"))
	if err != nil{
		log.Fatal(err)
	}
	
	UserCollection := client.Database("DevEnvi").Collection("User")
	ctx = context.TODO()
	filter := bson.M{"email" :user.Email}
	errfind := UserCollection.FindOne(ctx, filter).Decode(&result)
	if errfind != nil{
		log.Fatal(errfind)
	}
	
	if result.Password == user.Password { 
		json.NewEncoder(w).Encode(&result)
	} else {
		json.NewEncoder(w).Encode("Not valid")
	}

}
func GetAllProducts(w http.ResponseWriter, r *http.Request){
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@enviskates-zpjae.gcp.mongodb.net/admin?retryWrites=true&w=majority"))
	if err != nil{
		log.Fatal(err)
	}	
	UserCollection := client.Database("DevEnvi").Collection("test")
	cur ,err := UserCollection.Find(ctx,bson.D{{}})


	var results []Product

	for cur.Next(ctx){

		var product Product
		errProduct := cur.Decode(&product)
		if errProduct != nil{
			log.Fatal(errProduct)
		}

		results = append(results, product)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	products := ProductArray{
		Products: results,
	}

	json.NewEncoder(w).Encode(products)


}
func CreateOrder(w http.ResponseWriter, r *http.Request){}
func GetAllOrders(w http.ResponseWriter, r *http.Request){}
func GetCarbonEmmision(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)


	carEmissionFloatVal, err := strconv.ParseInt(keyVal["Car"],10,64);
	busEmissionFloatVal , err := strconv.ParseInt(keyVal["Bus"],10, 64)
	TaxiEmissionFloatVal , err := strconv.ParseInt(keyVal["Taxi"], 10,64)
	TrainEmissionFloatVal , err := strconv.ParseInt(keyVal["Train"], 10,  64)
	TramEmissionFloatVal , err := strconv.ParseInt(keyVal["Tram"], 10, 64)


	// 0.0001 meteric ton/km  ->  annually
	carEmission := int64(1)
	busEmission := int64(6)
	TaxiEmission := int64(4)
	TramEmission := int64(3)
	TrainEmission := int64(10)

	totalEmission := carEmission *  carEmissionFloatVal + busEmission * busEmissionFloatVal + TaxiEmission * TaxiEmissionFloatVal + TrainEmission * TrainEmissionFloatVal + TramEmission * TramEmissionFloatVal

	carbonData := carbonEmission{
		CarbonEmissionVal:  totalEmission,
	}

	json.NewEncoder(w).Encode(carbonData)
}

