package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"prushka/internal/db"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Data struct {
	Data any `json:"data"`
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func prepareData(rawStruct any) Data {
	rt := reflect.TypeOf(rawStruct)
	if rt.Kind() == reflect.Struct {
		return Data{[]any{rawStruct}}
	}

	return Data{rawStruct}
}

func WriteModelToJson(model *gorm.DB, container any, w http.ResponseWriter, limit int) {
	if limit == 0 {
		model.Find(&container)
	} else {
		model.Limit(limit).Find(&container)
	}
	respBody, _ := json.Marshal(prepareData(container))
	fmt.Fprint(w, string(respBody))
}

func Write404ToJson(w http.ResponseWriter) {
	respBody, _ := json.Marshal(struct {
		Text string `json:"message"`
	}{Text: "404 Not found"})
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, string(respBody))
}

func WriteDeletedToJson(w http.ResponseWriter) {
	respBody, _ := json.Marshal(struct {
		Text string `json:"message"`
	}{Text: "Successfuly deleted"})
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, string(respBody))
}

func seedGenerator(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func prepareDataParams(data map[string][]string) map[string]interface{} {
	dataParams := make(map[string]interface{})
	for key, value := range data {
		if key == "id" {
			dataParams["ID"] = value[0]
			continue
		}
		if strings.Contains(key, "_id") {
			dataParams[key] = value[0]
			continue
		}
		r := []rune(key)
		r[0] = unicode.ToUpper(r[0])
		s := string(r)
		dataParams[s] = value[0]
	}
	return dataParams
}

func sendRegisterEmail(email string, id int32) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	randomWord := seedGenerator(10)
	token := sha256.Sum256([]byte(fmt.Sprint(id) + randomWord))
	link := os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/auth/?token=" + base64.URLEncoding.EncodeToString(token[:])
	db.RDB.RPush(ctx, fmt.Sprint(id), randomWord)

	// user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	subject := "Prushka. Подтверждение адреса электронной почты"
	body := fmt.Sprintf(
		"Для подтверждения адреса вам необходимо перейти по следующей ссылке: <a href=%s>%s</a>",
		link, link,
	)

	request := Mail{
		Sender:  os.Getenv("SMTP_SENDER"),
		To:      []string{email},
		Subject: subject,
		Body:    body,
	}

	addr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")
	msg := BuildMessage(request)

	auth := smtp.PlainAuth("", os.Getenv("SMTP_SENDER"), password, os.Getenv("SMTP_HOST"))
	e := smtp.SendMail(addr, auth, os.Getenv("SMTP_SENDER"), []string{email}, []byte(msg))

	if e != nil {
		log.Fatal(e)
	}
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func getUserAndPassword(r http.Request) (*db.User, string) {
	var user db.User
	usrModel := db.DB.Model(db.User{})
	r.ParseForm()
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	usrModel.Where("email = ?", email).First(&user)

	return &user, password
}

func checkPassword(u db.User, password string) bool {
	passwordRandomWord := db.RDB.LPop(ctx, fmt.Sprint(u.ID))
	randomWord, _ := passwordRandomWord.Result()
	result := db.RDB.LPush(ctx, fmt.Sprint(u.ID), randomWord)
	if _, err := result.Result(); err != nil {
		log.Fatal(err)
	}

	hashedPassword := sha256.Sum256([]byte(password + randomWord))

	return base64.URLEncoding.EncodeToString(hashedPassword[:]) == u.Password
}

func writeJWT(u db.User) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	expire := time.Now()
	if u.IsAuthorized {
		expire.Add(72 * time.Hour)
	} else {
		expire.Add(24 * time.Hour)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"UserId": u.ID, "ExpiresAt": jwt.NewNumericDate(expire)})
	ss, _ := token.SignedString([]byte(os.Getenv("JWT_SIGN_STRING")))

	return ss
}

func parseJWT(r http.Request) bool {
	tokenString := r.Header.Get("Authorization")
	tokens := strings.Fields(tokenString)

	if len(tokens) != 2 {
		return false
	}

	if tokens[0] != "Bearer" {
		return false
	}

	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokens[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGN_STRING")), nil
	}, jwt.WithLeeway(5*time.Second))
	var userID int32
	for _, el := range claims {
		val, ok := el.(int32)
		if ok {
			userID = val
		}
	}
	var usr db.User
	db.DB.Model(db.User{}).Where("id = ?", int(userID)).First(&usr)
	return fmt.Sprint(userID) == fmt.Sprint(usr.ID) && token.Valid
}
