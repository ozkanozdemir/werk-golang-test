package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var url string
var hasuraAdminSecret string
var method string = "POST"

type User struct {
	ID         int
	Email      string
	Name       string
	Surname    string
	created_at string
	updated_at string
}

type Company struct {
	ID         int
	Name       string
	created_at string
	updated_at string
}

type Ralationship struct {
	User    string
	Company string
}

func main() {
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	url = os.Getenv("HASURA_URL")
	hasuraAdminSecret = os.Getenv("HASURA_ADMIN_SECRET")

	// create a new router
	router := mux.NewRouter()

	// specify endpoints, handler functions and HTTP method
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", GetUser).Methods("GET")
	router.HandleFunc("/users", StoreUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}", DeleteUser).Methods("DELETE")
	router.HandleFunc("/companies", GetCompanies).Methods("GET")
	router.HandleFunc("/companies/{id:[0-9]+}", GetCompany).Methods("GET")
	router.HandleFunc("/companies", StoreCompany).Methods("POST")
	router.HandleFunc("/companies/{id:[0-9]+}", UpdateCompany).Methods("PUT")
	router.HandleFunc("/companies/{id:[0-9]+}", DeleteCompany).Methods("DELETE")

	router.HandleFunc("/users-companies/add-relationship", AddRelationship).Methods("POST")
	router.HandleFunc("/users-companies/remove-relationship", RemoveRelationship).Methods("DELETE")
	http.Handle("/", router)

	// start and listen to requests
	fmt.Println("Listening :8080")
	http.ListenAndServe(":8080", router)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	graphqlData := map[string]string{
		"query": `{
			users {
				id
				email
				firstname
				lastname
				created_at
				updated_at
				companies {
				  company {
					id
					name
					created_at,
					updated_at
				  }
				}
			}
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	// w.Write(jsonResponse)
	w.Write(body)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	graphqlData := map[string]string{
		"query": `{
			users_by_pk(id: ` + vars["id"] + `) {
				id
				email
				firstname
				lastname
				created_at
				updated_at
				companies {
				  company {
					id
					name
					created_at,
					updated_at
				  }
				}
			}
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	// w.Write(jsonResponse)
	w.Write(body)
}

func StoreUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	graphqlData := map[string]string{
		"query": `
		mutation
		{
			insert_users(objects: [{
				email: "` + user.Email + `",
				firstname: "` + user.Name + `",
				lastname: "` + user.Surname + `"
			}]) {
				returning {
				  	id
				  	email
				  	firstname
				  	lastname
				  	created_at
				  	updated_at
				}
			  }
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	graphqlData := map[string]string{
		"query": `
		mutation {
		  update_users_by_pk(
			pk_columns: {id: ` + vars["id"] + `}, 
			_set: {
			  	email: "` + user.Email + `",
				firstname: "` + user.Name + `",
				lastname: "` + user.Surname + `"
			}) {
			id
			email
			firstname
			lastname
			created_at
			updated_at
		  }
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	graphqlData := map[string]string{
		"query": `
		mutation {
		  delete_users_by_pk(id: ` + vars["id"] + `) {
			id
		  }
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	graphqlData := map[string]string{
		"query": `{
			companies_by_pk(id: ` + vars["id"] + `) {
				id
				name
				created_at
				updated_at
				users {
				  user {
					id
					email
					firstname
					lastname
					created_at,
					updated_at
				  }
				}
			}
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func GetCompanies(w http.ResponseWriter, r *http.Request) {
	graphqlData := map[string]string{
		"query": `{
			companies {
				id
				name
				created_at
				updated_at
				users {
				  user {
					id
					email
					firstname
					lastname
					created_at,
					updated_at
				  }
				}
			}
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	// w.Write(jsonResponse)
	w.Write(body)
}

func StoreCompany(w http.ResponseWriter, r *http.Request) {
	var company Company

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		fmt.Println(err)
		return
	}

	graphqlData := map[string]string{
		"query": `
		mutation
		{
			insert_companies(objects: [{
				name: "` + company.Name + `"
			}]) {
				returning {
				  	id
				  	name
				  	created_at
				  	updated_at
				}
			  }
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var company Company

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		fmt.Println(err)
		return
	}

	graphqlData := map[string]string{
		"query": `
		mutation {
		  update_companies_by_pk(
			pk_columns: {id: ` + vars["id"] + `}, 
			_set: {
			  	name: "` + company.Name + `"
			}) {
			id
			name
			created_at
			updated_at
		  }
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	graphqlData := map[string]string{
		"query": `
		mutation {
		  delete_companies_by_pk(id: ` + vars["id"] + `) {
			id
		  }
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func AddRelationship(w http.ResponseWriter, r *http.Request) {
	var relationship Ralationship

	err := json.NewDecoder(r.Body).Decode(&relationship)
	if err != nil {
		fmt.Println(err)
		return
	}

	graphqlData := map[string]string{
		"query": `
		mutation {
			insert_user_company(objects: {user_id: ` + relationship.User + `, company_id: ` + relationship.Company + `}) {
				returning {
				  id
				  company {
					id
					name
					created_at
					updated_at
				  }
				  user {
					id
					email
					firstname
					lastname
					created_at
					updated_at
				  }
				}
			}
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}

func RemoveRelationship(w http.ResponseWriter, r *http.Request) {
	var relationship Ralationship

	err := json.NewDecoder(r.Body).Decode(&relationship)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(relationship)

	graphqlData := map[string]string{
		"query": `
		mutation {
			delete_user_company(where: {_and: {user_id: {_eq: ` + relationship.User + `}, company_id: {_eq: ` + relationship.Company + `}}}) {
				returning {
				  company_id
				  user_id
				  id
				}
			}
		}`,
	}
	jsonData, _ := json.Marshal(graphqlData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// update content type
	w.Header().Set("Content-Type", "application/json")

	// specify HTTP status code
	w.WriteHeader(http.StatusOK)

	// update response
	w.Write(body)
}
