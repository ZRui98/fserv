package main

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"time"
	"flag"
	"crypto/rand"
	"encoding/binary"

	"github.com/alexedwards/argon2id"
	"github.com/gabriel-vasile/mimetype"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/golang/glog"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length = uint32(len(alphabet))
)

var JWT_KEY = []byte(os.Getenv("SECRET_KEY"))
var ROOT_DIR = os.Getenv("FILE_DIR")
var REGISTRATION_KEY = os.Getenv("REGISTRATION_KEY")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encode(n uint32) string {
	if n == 0 {
		return "0"
	}

	b := make([]byte, 0, 512)
	for n > 0 {
		r := n % length
		n /= length
		b = append([]byte{alphabet[r]}, b...)
	}
	return string(b)
}
var image = []string{"image/jpeg", "image/png", "image/gif"}
var video = []string{"video/x-flv", "video/mp4", "video/mov", "video/x-msvideo", "video/x-ms-wmv"}
var userFolders = []string{"txt", "videos", "images"}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		jwtRetCode, username := validateJWT(r)
		if jwtRetCode != http.StatusOK {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		glog.Info("Uploading File")
		r.ParseMultipartForm(32 << 20)
		uploadedFile, handler, err := r.FormFile("uploadFile")
		if err != nil {
			glog.Errorf("Error parsing form file:: %v\n", err)
			return
		}
		defer uploadedFile.Close()
		var seed []byte
		seed, err = generateRandomBytes(4)
		if err != nil {
			seed, err = generateRandomBytes(4)
			glog.Errorf("Error generating seed:: %v\n", err)
		}

		mime, err := mimetype.DetectReader(uploadedFile)
		if err != nil {
			glog.Errorf("Error unrecognized mimetype:: %v\n", err)
			return
		}
		folder := "txt"
		if mimetype.EqualsAny(mime.String(), image...) {
			folder = "images"
		} else if mimetype.EqualsAny(mime.String(), video...) {
			folder = "videos"
		}
		token := encode(binary.BigEndian.Uint32(seed))
		file_path := path.Join(ROOT_DIR, "files", username, folder, token + path.Ext(handler.Filename))
		glog.Info("Saving to: %s File: %s", file_path, handler.Filename)
		queryString := `INSERT INTO "files" ("url_id", "file_path", "file_name", "owner") VALUES ($1, $2, $3, $4);`
		_, err = pool.Exec(context.Background(), queryString, token, file_path, handler.Filename, username)
		if err != nil {
			glog.Errorf("Error inserting data into DB: %v\n", err)
			return
		}

		f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			glog.Errorf("Error creating file pointer: %v\n", err)
			return
		}
		defer f.Close()
		if _, err := uploadedFile.Seek(0, 0); err != nil {
			glog.Errorf("Error seeking uploaded file: %v\n", err)
			return
		}
		io.Copy(f, uploadedFile)
		t, _ := template.ParseFiles("templates/upload.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("templates/upload.html")
		t.Execute(w, nil)
	}
}

func validateJWT(r *http.Request) (int, string) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, ""
		}
		return http.StatusBadRequest, ""
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, ""
		}
		return http.StatusBadRequest, ""
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, ""
	}
	var lastLogin time.Time
	err = pool.QueryRow(context.Background(), `SELECT "last_login" FROM "users" WHERE "username"=$1`, claims.Username).Scan(&lastLogin)
	if err != nil {
		glog.Errorf("Getting last login time failed: %v\n", err)
		return http.StatusInternalServerError, ""
	}
	if time.Unix(claims.IssuedAt, 0).Before(lastLogin) {
		glog.Error("User already logged in, token expired")
		return http.StatusRequestTimeout, ""
	}
	return http.StatusOK, claims.Username
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			glog.Errorf("Parsing form failed:: %v\n", err)
			return
		}
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		var hashedPassword string
		err = pool.QueryRow(context.Background(), `SELECT "password" FROM "users" WHERE "username"=$1`, username).Scan(&hashedPassword)
		if err != nil {
			glog.Errorf("Querying DB failed:: %v\n", err)
			return
		}
		match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
		if err != nil {
			glog.Errorf("Error validating password:: %v\n", err)
			return
		}

		if !match {
			w.WriteHeader(http.StatusUnauthorized)
			glog.Error("Password did not match")
			return
		}

		loginTimestamp := time.Now()
		_, err = pool.Exec(context.Background(), `UPDATE "users" SET last_login=$1 WHERE username=$2`, loginTimestamp, username)
		expirationTime := time.Now().Add(3 * time.Hour)
		claims := &Claims {
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				IssuedAt: loginTimestamp.Unix() + 1,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JWT_KEY)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value: tokenString,
			Expires: expirationTime,
			Secure: true,
			SameSite: http.SameSiteNoneMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
	}
}

func register( w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			glog.Errorf("Parsing form failed:: %v\n", err)
			return
		}
		glog.Info("Creating User")
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		registration_key := r.PostFormValue("key")
		if (registration_key != REGISTRATION_KEY) {
			glog.Error("Wrong Registration Key")
			return
		}
		hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			glog.Errorf("Failed to verify password:: %v\n", err)
			return
		}
		loginTimestamp := time.Now()
		queryString := `INSERT INTO "users" (username, password, last_login) VALUES ($1, $2, $3);`
		_, err = pool.Exec(context.Background(), queryString, username, hashedPassword, loginTimestamp)
		if err != nil {
			glog.Errorf("Querying DB failed:: %v\n", err)
			return
		}
		userDir := "files/" + username
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			os.Mkdir(userDir, os.ModePerm)
			for _, s := range userFolders {
				os.Mkdir(userDir + "/" + s, os.ModePerm)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/register.html")
		t.Execute(w, nil)
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters, present := query["id"]
	if !present || len(filters) == 0 {
		glog.Error("ID not present/valid")
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/album.html")
		t.Execute(w, nil)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func setupRoutes() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/s", view)
}

func usage() {
    flag.PrintDefaults()
    os.Exit(2)
}

func init() {
    flag.Usage = usage
    flag.Set("logtostderr", "true")
    flag.Set("v", "2")
    flag.Parse()
}
var pool *pgxpool.Pool
func main() {
	glog.Info("Starting fserv")
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		glog.Fatalf("Error in DB Configuration:: %v\n", err)
		os.Exit(1)
	}
	pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		glog.Fatalf("Unable to connect to database:: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	setupRoutes()
	err = http.ListenAndServeTLS(":2446", "cert.pem", "key.pem", nil)
	glog.Flush()
}
