package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Token string
}

type ReqCreateUserJSON struct {
	Name string `json:"name" validate:"required"`
}

type ResCreateUserJSON struct {
	Token string `json:"token"`
}

type ResGetUserJSON struct {
	Name string `json:"name"`
}

type ReqUpdateUserJSON struct {
	Name string `json:"name" validate:"required"`
}

func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "ca_mission"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	count := 0
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				panic(err)
			}
			db, err = gorm.Open(DBMS, CONNECT)
		}
	}
	return db
}

// UserTokenをunixtimeから生成して返す関数
func GenerateUserToken() string {
	//生成したuuidが被っていないかチェックするようにした方が良いかも
	unixtime := strconv.FormatInt(time.Now().Unix(), 10)
	return unixtime
}

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world.")
	// header := r.Header
	// fmt.Println(header)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// JSONに変換
	var userCreate ReqCreateUserJSON
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body[:length], &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// MySQLへデータを格納
	user := User{
		Name:  jsonBody["name"].(string),
		Token: GenerateUserToken(),
	}
	db := sqlConnect()
	db.AutoMigrate(&user)
	defer db.Close()
	db.NewRecord(&user)
	db.Create(&user)

	// 返り値の設定
	userCreate.Name = jsonBody["name"].(string)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(
		&ResCreateUserJSON{
			Token: user.Token,
		})
	if err != nil {
		log.Fatal(err)
	}
}

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.Header.Get("x-token"))

// 	if err := userUsecase.Get(m); err != nil {
// 		log.Fatal(err)
// 		writeError(w, http.StatusInternalServerError, apierror.ErrInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json;charset=utf-8")
// 	w.WriteHeader(http.StatusOK)

// 	if err := json.NewEncoder(w).Encode(&ResGetUserJSON{
// 		Name: m.Name,
// 	}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	//jsonからgoの構造体にデコードする
// 	var user ReqUpdateUserJSON
// 	//http通信などのストリームデータをデコード
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	//バリデーション
// 	validate := validator.New()
// 	err = validate.Struct(&user)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	m := &model.User{
// 		Name:  user.Name,
// 		Token: r.Header.Get("x-token"),
// 	}

// 	err = userUsecase.Update(m)
// 	if err != nil {
// 		log.Fatal(err)
// 		writeError(w, http.StatusInternalServerError, apierror.ErrInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json;charset=utf-8")
// 	w.WriteHeader(http.StatusNoContent)
// }

func main() {
	fmt.Println("Starting Server at http://localhost:8080")
	// db := sqlConnect()
	// db.AutoMigrate(&User{})
	// defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/test", HandlerFunc)

	subrouterUser := router.PathPrefix("/user").Subrouter()
	subrouterUser.HandleFunc("/create", CreateUser).Methods("POST")
	// subrouterUser.HandleFunc("/get", GetUser).Methods("GET")
	// subrouterUser.HandleFunc("/update", UpdateUser).Methods("PUT")
	router.Handle("/", router)
	http.ListenAndServe(":8080", router)
}
