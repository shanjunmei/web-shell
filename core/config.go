package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config Represents Server ,Database  and port details
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Users    []User
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Database        string
	Schema          string
	User            string
	Password        string
	MigrateRequired bool
}

type ServerConfig struct {
	Port string
	Host string
}
type User struct {
	User     string
	Password string
}

func (config Config) ContainUser(user, password string) bool {
	return Contain(config.Users, User{User: user, Password: password}, func(s, t User) bool {
		return (s.User == t.User && s.Password == t.Password)
	})
}

func Contain[T any](s []T, d T, f func(e1, e2 T) bool) bool {
	for _, k := range s {
		if f(k, d) {
			return true
		}
	}
	return false
}

func LoadConfig(path string) Config {
	var config Config
	LoadFile(path, &config)
	return config
}

func LoadFile(path string, entity interface{}) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(buf, &entity)
}

var config = Config{}

// Init Represents to Read and parse the configuration file
func Init(filepath string, logFile *os.File) (*gorm.DB, error) {
	config = LoadConfig(filepath)
	//Create Log Folder Path in the application root directory

	db, err := GetPostgresDB()
	if err != nil {
		log.Println(err)
	}
	return db, err
}

// Get Represents to get instance of  Config
func Get() Config {
	if &config == nil {
		panic("config have init")
	}
	return config
}

// GetPostgresDB - Create SQL Database
func GetPostgresDB() (*gorm.DB, error) {
	host := Get().Database.Host
	user := Get().Database.User
	port := Get().Database.Port
	password := Get().Database.Password
	databaseName := Get().Database.Database
	schema := Get().Database.Schema
	desc := fmt.Sprintf("host=%s  port=%d  user=%s  password=%s dbname=%s sslmode=disable search_path=%s binary_parameters=yes", host, port, user, password, databaseName, schema)
	log.Println("DB Connection : " + desc)
	db, err := createConnection(desc)
	return db, err
}

func createConnection(desc string) (*gorm.DB, error) {
	_db, err := sql.Open("postgres", desc)
	if err != nil {
		return nil, err
	}
	config := postgres.Config{Conn: _db}
	return gorm.Open(postgres.New(config))
}
