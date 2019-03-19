package main
 
import (
    "fmt"
    "net/http"
	"database/sql"
    "log"
	config "./config" 
    mydb "./mydb"
    ms "./email"
	_ "github.com/lib/pq"
    "os"
    "strconv"
    "strings"
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

	registration.ID = "email-service"
	registration.Name = "email-service"
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
	p := os.Getenv("EMAIL_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8083"
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
    http.HandleFunc("/CreateICSfromDBbyID",CreateICSfromDBbyID)
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

//Create an ics file from the event information in database based on ID
func CreateICSfromDBbyID(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email := r.FormValue("email") // Data from the form        
    if event, err1 := mydb.GetEventByID(email); err1 == nil {
        log.Printf("%v\n", event)
        return
    } else {
        log.Printf("error was: %v\n",err1)
    }        
}

//Sending an ICS as email
func SendEmail(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email := r.FormValue("email")
    name := r.FormValue("name")
    subject := r.FormValue("subject")
    message := r.FormValue("message")
    ms.SendEmail(name,email,subject,message)
}