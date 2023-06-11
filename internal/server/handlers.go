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
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

// UserHandler godoc
// @Summary Work with users
// @Description CRUD and auth users
// @Tags users
// @Produce  json
// @Param all query bool false "Get all users"
// @Param limit query int false "Get n users"
// @Param user_id path int false "Get user with id"
// @Success 200 {object} db.User
// @Success 200 {array} db.User
// @Failure 404
// @Param Authorization header string false "Insert your access token" default(Bearer <Add access token here>)
// @Router /users/{user_id}/ [get]
// @Router /users/{user_id}/ [post]
// @Router /users/{user_id}/ [delete]
// @Router /users/ [get]
// @Router /users/ [post]
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

			user = db.User{
				ID:           int32(id),
				Name:         r.PostFormValue("name"),
				Email:        r.PostFormValue("email"),
				Password:     base64.URLEncoding.EncodeToString(hashedPassword[:]),
				IsAuthorized: false,
			}

			if user.Name == "" || user.Email == "" {
				respBody, _ := json.Marshal(prepareData("Bad request"))
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, string(respBody))
				return
			}

			if err := usrModel.Create(&user).Error; err != nil {
				respBody, _ := json.Marshal(prepareData(err.Error()))
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, string(respBody))
				return
			}

			result := db.RDB.RPush(ctx, fmt.Sprint(id), randomWord)
			if _, err := result.Result(); err != nil {
				log.Fatal(err)
			}

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

// ModelHandler godoc
// @Summary Work with all existing models
// @Description CRUD for models
// @Produce  json
// @Param all query bool false "Get all models"
// @Param limit query int false "Get n models"
// @Success 200 {object} db.Privilege
// @Success 200 {array} db.Privilege
// @Success 200 {object} db.Card
// @Success 200 {array} db.Card
// @Success 200 {object} db.Attachment
// @Success 200 {array} db.Attachment
// @Success 200 {object} db.Label
// @Success 200 {array} db.Label
// @Success 200 {object} db.Column
// @Success 200 {array} db.Column
// @Success 200 {object} db.Desk
// @Success 200 {array} db.Desk
// @Success 200 {object} db.Workspace
// @Success 200 {array} db.Workspace
// @Success 200 {object} db.UserPrivilege
// @Success 200 {array} db.UserPrivilege
// @Success 200 {object} db.CardsLabel
// @Success 200 {array} db.CardsLabel
// @Failure 404
// @Failure 401
// @Param Authorization header string false "Insert your access token" default(Bearer <Add access token here>)
// @Param privilege_id path int false "Get privilege with id"
// @Router /privilege/{privilege_id}/ [get]
// @Router /privilege/{privilege_id}/ [post]
// @Router /privilege/{privilege_id}/ [delete]
// @Router /privilege/ [get]
// @Router /privilege/ [post]
// @Param cards_id path int false "Get card with id"
// @Router /cards/{cards_id}/ [get]
// @Router /cards/{cards_id}/ [post]
// @Router /cards/{cards_id}/ [delete]
// @Router /cards/ [get]
// @Router /cards/ [post]
// @Param attachment_id path int false "Get attachment with id"
// @Router /attachment/{attachment_id}/ [get]
// @Router /attachment/{attachment_id}/ [post]
// @Router /attachment/{attachment_id}/ [delete]
// @Router /attachment/ [get]
// @Router /attachment/ [post]
// @Param label_id path int false "Get label with id"
// @Router /label/{label_id}/ [get]
// @Router /label/{label_id}/ [post]
// @Router /label/{label_id}/ [delete]
// @Router /label/ [get]
// @Router /label/ [post]
// @Param column_id path int false "Get column with id"
// @Router /column/{column_id}/ [get]
// @Router /column/{column_id}/ [post]
// @Router /column/{column_id}/ [delete]
// @Router /column/ [get]
// @Router /column/ [post]
// @Param desk_id path int false "Get desk with id"
// @Router /desk/{desk_id}/ [get]
// @Router /desk/{desk_id}/ [post]
// @Router /desk/{desk_id}/ [delete]
// @Router /desk/ [get]
// @Router /desk/ [post]
// @Param workspace_id path int false "Get workspace with id"
// @Router /workspace/{workspace_id}/ [get]
// @Router /workspace/{workspace_id}/ [post]
// @Router /workspace/{workspace_id}/ [delete]
// @Router /workspace/ [get]
// @Router /workspace/ [post]
// @Param user_privilege_id path int false "Get user_privilege with id"
// @Router /user_privilege/{user_privilege_id}/ [get]
// @Router /user_privilege/{user_privilege_id}/ [post]
// @Router /user_privilege/{user_privilege_id}/ [delete]
// @Router /user_privilege/ [get]
// @Router /user_privilege/ [post]
// @Param cards_label_id path int false "Get cards_label with id"
// @Router /cards_label/{cards_label_id}/ [get]
// @Router /cards_label/{cards_label_id}/ [post]
// @Router /cards_label/{cards_label_id}/ [delete]
// @Router /cards_label/ [get]
// @Router /cards_label/ [post]
func ModelHandler(model interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := reflect.TypeOf(model)
		modelList := reflect.New(reflect.SliceOf(t)).Elem().Interface()
		modelInstance := reflect.New(t).Interface()
		mainModel := db.DB.Model(model)

		ids := r.URL.Query()["id"]
		is_all := r.URL.Query()["all"]
		limit := r.URL.Query()["limit"]
		id := ""
		if len(ids) > 0 {
			id = ids[0]
		}

		r.ParseForm()

		if id == "" { // url without id
			if r.Method == http.MethodPost {

				id, _ := strconv.Atoi(r.PostFormValue("id"))
				data := prepareDataParams(r.PostForm)

				if len(data) == 0 {
					respBody, _ := json.Marshal(prepareData("Bad request"))
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprint(w, string(respBody))
					return
				}

				if err := mainModel.Create(data).Error; err != nil {
					respBody, _ := json.Marshal(prepareData(err.Error()))
					w.WriteHeader(http.StatusConflict)
					fmt.Fprint(w, string(respBody))
					return
				}

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
		if !parseJWT(*r) { // Если токена нет
			if checkPassword(*user, password) { // Если пароль совпадает
				token := writeJWT(*user)
				w.WriteHeader(http.StatusOK)
				respBody, _ := json.Marshal(prepareData(struct {
					Token string `json:"token"`
				}{Token: token}))
				fmt.Fprint(w, string(respBody))
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				respBody, _ := json.Marshal(prepareData("Unauthorized"))
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

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "swagger") || strings.Contains(req.URL.Path, "auth") || strings.HasSuffix(req.URL.Path, "users/") && req.Method == http.MethodPost {
			next.ServeHTTP(w, req.Clone(ctx))
		} else {
			if !parseJWT(*req) {
				w.WriteHeader(http.StatusUnauthorized)
				respBody, _ := json.Marshal(prepareData("Unauthorized"))
				fmt.Fprint(w, string(respBody))
				return
			} else {
				next.ServeHTTP(w, req.Clone(ctx))
			}
		}
	})
}
