package function

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type (
	//RequestBody JSON Request Struct
	RequestBody struct {
		Message string `json:"message"`
	}

	//JSONResponse Response Struct
	JSONResponse struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

var err error
var db *gorm.DB

var logger *log.Logger

//Initialization for logs
func init() {
	logger = log.New()
	logger.SetLevel(log.TraceLevel)
	logger.Formatter = &log.TextFormatter{}
	log.SetOutput(os.Stdout)
}

//InitializeDB Function for Database Connection
func InitializeDB() error {
	Logger().Info("IntializeDB Function")

	cfg := mysql.Config{
		User:                 "newroot",
		Passwd:               "newroot",
		Addr:                 "192.168.0.254:3306", //IP:PORT
		Net:                  "tcp",
		DBName:               "car_inventory",
		AllowNativePasswords: true,
	}
	psqlInfo := cfg.FormatDSN()

	db, err = gorm.Open("mysql", psqlInfo)
	if err != nil {
		return err
	}
	return nil
}

//Handle Handler Function for data processing and Response
func Handle(w http.ResponseWriter, r *http.Request) {
	Logger().Info("Handler Function")

	var requestBody RequestBody
	var response JSONResponse

	if r.Body == nil {
		Logger().Error("Body is nil")
		response.Status = "ERROR"
		response.Message = "Please send a request body"
		response.Data = "Request Body cannot be empty"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		Logger().Error(err)
		response.Status = "ERROR"
		response.Message = "Error in Request Body"
		response.Data = "Please send a correct request body"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = InitializeDB()
	defer db.Close()
	if err != nil {
		Logger().Error(err)
		response.Status = "ERROR"
		response.Message = "Error in Connecting to the Database"
		response.Data = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	db.Raw("SELECT CURRENT_TIMESTAMP FROM DUAL").Row().Scan(&response.Data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	Logger().Info("Success")
	json.NewEncoder(w).Encode(response)
}

//Logger with fields
func Logger() *log.Entry {
	var depth = 1
	function, file, line, _ := runtime.Caller(depth)
	functionObject := runtime.FuncForPC(function)
	entry := logger.WithFields(log.Fields{
		"file":     file,
		"function": functionObject.Name(),
		"line":     line,
	})

	logger.SetOutput(os.Stdout)
	return entry

}
