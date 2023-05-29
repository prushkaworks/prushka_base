package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"prushka/internal/db"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

// func MainHandler(w http.ResponseWriter, r *http.Request) {
// 	var users []db.User
// 	user := db.User{
// 		ID: 23, Name: "Dmitriy", Email: "safr@328392mail.ru", Password: "3232323",
// 	}
// 	usr_model := db.DB.Db.Model(db.User{})
// 	usr_model.Create(&user)
// 	usr_model.Find(&users)
// 	respBody, _ := json.Marshal(users)
// 	fmt.Fprint(w, string(respBody))
// }

func UserHandler(w http.ResponseWriter, r *http.Request) {
	var users []db.User
	var user db.User
	usrModel := db.DB.Model(db.User{})

	is_all := r.URL.Query()["all"]
	limit := r.URL.Query()["limit"]
	id := mux.Vars(r)["id"]

	if id == "" { // url without id
		if r.Method == http.MethodPost {
			r.ParseForm()
			id, _ := strconv.Atoi(r.PostFormValue("id"))
			password := r.PostFormValue("password")
			randomWord := seedGenerator(10)
			hashedPassword := sha256.Sum256([]byte(password + randomWord))
			result := db.RDB.RPush(ctx, fmt.Sprint(id), randomWord)
			if _, err := result.Result(); err != nil {
				log.Fatal(err)
			}

			user = db.User{
				ID:           int32(id),
				Name:         r.PostFormValue("name"),
				Email:        r.PostFormValue("email"),
				Password:     base64.URLEncoding.EncodeToString(hashedPassword[:]),
				IsAuthorized: false,
			}

			usrModel.Create(&user)
			usrModel.Find(&user, id)
			if user.ID == int32(id) {
				sendRegisterEmail(user.Email, user.ID)
			}
			respBody, _ := json.Marshal(prepareData(user))
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, string(respBody))
			return
		}
		if len(limit) != 0 {
			limit_int, _ := strconv.Atoi(limit[0])
			WriteModelToJson(usrModel, users, w, limit_int)
			return
		}
		if len(is_all) != 0 {
			if all, _ := strconv.ParseBool(is_all[0]); all {
				WriteModelToJson(usrModel, users, w, 0)
				return
			} else {
				WriteModelToJson(usrModel, users, w, 20)
				return
			}
		} else {
			anotherParams := r.URL.Query()
			if len(anotherParams) != 0 {
				filteredParams := make(map[string]string)
				for param, value := range anotherParams {
					if len(value[0]) > 20 || len(param) > 20 {
						continue
					}
					filteredParams[param] = value[0]
				}

				db.DB.Where(filteredParams).Find(&users)
				respBody, _ := json.Marshal(prepareData(users))
				fmt.Fprint(w, string(respBody))
				return
			}
			WriteModelToJson(usrModel, users, w, 20)
			return
		}
	} else { // url with id

		id_int, _ := strconv.Atoi(id)

		if r.Method == http.MethodPost {
			r.ParseForm()
			usr := usrModel.Where("id = ?", id_int)
			for key, val := range r.PostForm {
				usr.Update(key, val[0])
			}
		} else if r.Method == http.MethodDelete {
			usrModel.Delete(&user, id_int)
			WriteDeletedToJson(w)
			return
		}

		usrModel.Find(&user, id_int)
		respBody, _ := json.Marshal(prepareData(user))
		fmt.Fprint(w, string(respBody))
		return
	}
}

func ModelHandler(model interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := reflect.TypeOf(model)
		modelList := reflect.New(reflect.SliceOf(t)).Elem().Interface()
		modelInstance := reflect.New(t).Interface()
		mainModel := db.DB.Model(model)

		is_all := r.URL.Query()["all"]
		limit := r.URL.Query()["limit"]
		id := mux.Vars(r)["id"]

		if id == "" { // url without id
			if r.Method == http.MethodPost {
				r.ParseForm()
				id, _ := strconv.Atoi(r.PostFormValue("id"))
				data := prepareDataParams(r.PostForm)

				mainModel.Create(data)
				mainModel.Find(&modelInstance, id)
				respBody, _ := json.Marshal(prepareData(modelInstance))
				fmt.Fprint(w, string(respBody))
				return
			}
			if len(limit) != 0 {
				limit_int, _ := strconv.Atoi(limit[0])
				WriteModelToJson(mainModel, modelList, w, limit_int)
				return
			}
			if len(is_all) != 0 {
				if all, _ := strconv.ParseBool(is_all[0]); all {
					WriteModelToJson(mainModel, modelList, w, 0)
					return
				} else {
					WriteModelToJson(mainModel, modelList, w, 20)
					return
				}
			} else {
				anotherParams := r.URL.Query()
				if len(anotherParams) != 0 {
					filteredParams := make(map[string]string)
					for param, value := range anotherParams {
						if len(value[0]) > 20 || len(param) > 20 {
							continue
						}
						filteredParams[param] = value[0]
					}

					db.DB.Where(filteredParams).Find(&modelList)
					respBody, _ := json.Marshal(prepareData(modelList))
					fmt.Fprint(w, string(respBody))
					return
				}
				WriteModelToJson(mainModel, modelList, w, 20)
				return
			}
		} else { // url with id

			id_int, _ := strconv.Atoi(id)

			if r.Method == http.MethodPost {
				r.ParseForm()
				ins := mainModel.Where("id = ?", id_int)
				for key, val := range r.PostForm {
					ins.Update(key, val[0])
				}
			} else if r.Method == http.MethodDelete {
				mainModel.Delete(&modelInstance, id_int)
				WriteDeletedToJson(w)
				return
			}

			mainModel.Find(&modelInstance, id_int)
			respBody, _ := json.Marshal(prepareData(modelInstance))
			fmt.Fprint(w, string(respBody))
		}
	}
}

func My404Handler(w http.ResponseWriter, r *http.Request) {
	Write404ToJson(w)
}

func LoggingAndJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Allow-Origin", "https://localhost:3000")

		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s%s\t%s", color.BlueString(req.Method), req.Host, req.URL.Path, color.CyanString(fmt.Sprintf("\t%fsec", time.Since(start).Seconds())))
	})
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, password := getUserAndPassword(*r)
	usrModel := db.DB.Model(db.User{})

	regToken := r.URL.Query()["token"]
	if len(regToken) != 0 {
		regRandomWord := db.RDB.RPop(ctx, fmt.Sprint(user.ID))
		randomWord, _ := regRandomWord.Result()
		hashedString := sha256.Sum256([]byte(fmt.Sprint(user.ID) + randomWord))

		if base64.URLEncoding.EncodeToString(hashedString[:]) == regToken[0] {
			usrModel.Find(&user, user.ID).Update("is_authorized", true)
		}
	} else {
		if !parseJWT(*r, *user) { // Если токена нет
			if checkPassword(*user, password) { // Если пароль совпадает
				writeJWT(w, *r, *user)
				w.WriteHeader(http.StatusOK)
				respBody, _ := json.Marshal(prepareData("ok"))
				fmt.Fprint(w, string(respBody))
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				respBody, _ := json.Marshal(prepareData("error"))
				fmt.Fprint(w, string(respBody))
				return
			}
		} else {
			w.WriteHeader(http.StatusOK)
			respBody, _ := json.Marshal(prepareData("Token!"))
			fmt.Fprint(w, string(respBody))
			return
		}
	}
}

// func MainHandler(w http.ResponseWriter, r *http.Request) {
// 	user, _ := getUserAndPassword(*r)
// 	if !parseJWT(*r, *user) { // Если токена нет
// 		w.WriteHeader(http.StatusUnauthorized)
// 		respBody, _ := json.Marshal(prepareData("Error!"))
// 		fmt.Fprint(w, string(respBody))
// 		return
// 	} else {
// 		w.WriteHeader(http.StatusOK)
// 		respBody, _ := json.Marshal(prepareData("Token!"))
// 		fmt.Fprint(w, string(respBody))
// 		return
// 	}
// }
