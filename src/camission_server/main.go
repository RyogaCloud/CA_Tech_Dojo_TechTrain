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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// usersテーブルにおける各データの形
type User struct {
	gorm.Model
	Name  string
	Token string
}

// charactersテーブルにおける各データの形
type Character struct {
	gorm.Model
	Name        string
	Probability int
}

// ユーザ情報作成APIで返す形
type ReqCreateUserJSON struct {
	Name string `json:"name"`
}

// ユーザ情報作成APIで返す形
type ResCreateUserJSON struct {
	Token string `json:"token"`
}

// ユーザ情報取得APIで返す形
type ResGetUserJSON struct {
	Name string `json:"name"`
}

// ガチャ実行APIで返す形
type ResGachaJSON struct {
	Results []ResultGachaJSON `json:"results"`
}

// ガチャ実行APIでガチャの結果を格納する形
type ResultGachaJSON struct {
	CharacterID int    `json:"characterID"`
	Name        string `json:"name"`
}

// MySQLと接続する関数
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

// uuidを生成して返す関数
func GenerateUserToken() (string, error) {
	//生成したuuidが被っていないかチェックするようにした方が良いかも
	uuid, err := uuid.NewRandom()
	return uuid.String(), err
}

// Hello worldを出力する関数
func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world.")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// POSTされたデータをJSONに変換
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
	var reqCreateUserJSON map[string]interface{}
	err = json.Unmarshal(body[:length], &reqCreateUserJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// usersテーブルへデータを格納
	token, err := GenerateUserToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ca_missionのMySQLに接続
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	// usersテーブルにデータを格納
	user := User{
		Name:  reqCreateUserJSON["name"].(string),
		Token: token,
	}
	db.Table("users").NewRecord(&user)
	db.Table("users").Create(&user)

	// 返り値の設定
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	// ca_missionのMySQLに接続
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	// usersテーブルからtokenが一致するデータを取得
	var user User
	db.Table("users").Where("token = ?", r.Header.Get("x-token")).First(&user)

	// 返り値の設定
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(
		&ResGetUserJSON{
			Name: user.Name,
		})
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// PUTされたデータをJSONに変換
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
	var reqUpdateUserJSON map[string]interface{}
	err = json.Unmarshal(body[:length], &reqUpdateUserJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ca_missionのMySQLに接続
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	// usersテーブルを書き換える
	db.Table("users").Model(&User{}).Where("token = ?", r.Header.Get("x-token")).Update("name", reqUpdateUserJSON["name"].(string))

	// 返り値の設定
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
}

func DrawGacha(w http.ResponseWriter, r *http.Request) {
	// POSTされたデータをJSONに変換
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
	var reqDrawGachaJSON map[string]interface{}
	err = json.Unmarshal(body[:length], &reqDrawGachaJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ca_missionのMySQLに接続
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	// 取得キャラクターを配列に格納
	var resultsGacha []ResultGachaJSON
	// ガチャの実行
}

func main() {
	// 起動したサーバURLの表示
	fmt.Println("Starting Server at http://localhost:8080")

	// ルーティングの作成
	router := mux.NewRouter()
	router.HandleFunc("/test", HandlerFunc)
	subrouterUser := router.PathPrefix("/user").Subrouter()
	subrouterUser.HandleFunc("/create", CreateUser).Methods("POST")
	subrouterUser.HandleFunc("/get", GetUser).Methods("GET")
	subrouterUser.HandleFunc("/update", UpdateUser).Methods("PUT")
	subrouterGacha := router.PathPrefix("/gacha").Subrouter()
	subrouterGacha.HandleFunc("/draw", DrawGacha).Methods("POST")
	router.Handle("/", router)
	http.ListenAndServe(":8080", router)
}
