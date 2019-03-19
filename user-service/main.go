package main
 
import (
    "fmt"
    "net/http"
	"database/sql"
    "log"
	config "./config"
	mydb "./mydb"
	helper "./helpers"
	"crypto/sha512"
	_ "github.com/lib/pq"
	"encoding/json"
    "os"
    "strconv"
    "strings"
	"github.com/dgrijalva/jwt-go"
    consulapi "github.com/hashicorp/consul/api"
)

type JwtToken struct {
	Token string `json:"token"`
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "user-service"
	registration.Name = "user-service"
	address := hostname()
	registration.Address = address
	port, err := strconv.Atoi(port()[1:len(port())]) 
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	consul.Agent().ServiceRegister(registration)
}

func port() string {
	p := os.Getenv("USER_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8081"
	}
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}

func main() {
    registerServiceWithConsul()
    db := connectToDatabase()
    http.HandleFunc("/Signup",Signup)
    http.HandleFunc("/ChangePassword",ChangePassword)
    http.HandleFunc("/login",login)
    http.HandleFunc("/CreateTable",CreateTable)
	fmt.Printf("user service is up on port: %s", port())
	http.ListenAndServe(port(), nil)
    defer db.Close()
}
//Database connection
func connectToDatabase() *sql.DB {
    dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    config.HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.PORT)
    db, err := sql.Open("postgres", dbinfo)
    if err != nil {
        fmt.Println(err)
    }
    log.Printf("Postgres started at %s PORT", config.PORT)
    mydb.SetDatabase(db)		
    return db
}

func Signup(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()
    uName := r.FormValue("username")     // Data from the form
    email := r.FormValue("email")        // Data from the form
    pwd := r.FormValue("password")       // Data from the form
    pwdConfirm := r.FormValue("confirm") // Data from the form

    // Empty data checking
    uNameCheck := helper.IsEmpty(uName)
    emailCheck := helper.IsEmpty(email)
    pwdCheck := helper.IsEmpty(pwd)
    pwdConfirmCheck := helper.IsEmpty(pwdConfirm)

    if uNameCheck || emailCheck || pwdCheck || pwdConfirmCheck {
        fmt.Fprintf(w, "ErrorCode is -10 : There is empty data.")
        return
    }
    if pwd == pwdConfirm {
        flag := mydb.Signup(uName,email,pwd)
        if flag == 1{
            fmt.Fprintln(w, "Account Created")
        }
    } else {
        fmt.Fprintln(w, "Password information must be the same.")
    }
}

func ChangePassword(w http.ResponseWriter,r *http.Request){
    r.ParseForm()
    email := r.FormValue("email")     // Data from the form
    oldPassword := r.FormValue("OldPassword")  // Data from the form
    newPassword := r.FormValue("NewPassword")  // Data from the form
    confirmPassword := r.FormValue("ConfirmPassword") // Data from the form
    if confirmPassword==newPassword{
        flag:=mydb.ChangePassword(email,oldPassword,newPassword)
        if flag == 1{
            fmt.Fprintln(w, "Password Changed successfully")
        }
    } else {
        fmt.Fprintln(w, "Error")
    }
}

func login(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    email := r.FormValue("email")  // Data from the form
    pwd := r.FormValue("password") // Data from the form

    // Empty data checking
    emailCheck := helper.IsEmpty(email)
    pwdCheck := helper.IsEmpty(pwd)

    if emailCheck || pwdCheck {
        fmt.Fprintf(w, "ErrorCode is -10 : There is empty data.")
        return
    }
    //Getting JWT
    if user, err := mydb.Login(email, pwd); err == nil {
        hasher := sha512.New()
        hasher.Write([]byte(pwd))
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "username": email,
        })
        tokenString, error := token.SignedString([]byte("secret"))
        if error != nil {
            fmt.Println(error)
        }
        json.NewEncoder(w).Encode(JwtToken{Token: tokenString})

        var file, err = os.Create(`creds.txt`)
        if err != nil {
            
        }  
        fmt.Fprintf(file,tokenString) 
        fmt.Fprintln(w,"Login Successful")
        defer file.Close()
        
        log.Printf("User has logged in: %v\n", user)
        return
    } else {
        log.Printf("Failed to log user in with email: %v %v, error was: %v\n", email,pwd, err)
    }
}

func CreateTable(w http.ResponseWriter, r *http.Request) {
    mydb.CreateTable()
}