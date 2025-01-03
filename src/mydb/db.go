package mydb
import (
 "database/sql"
 "fmt"
 "log"
 "os"
 "github.com/joho/godotenv"
 _ "github.com/lib/pq"
)

func setenv() string {
 err := godotenv.Load(".env")
 if err != nil {
  log.Fatalf("Error loading .env file: %v", err)
 }

 dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
  os.Getenv("POSTGRES_USER"),
  os.Getenv("POSTGRES_PASSWORD"),
  os.Getenv("POSTGRES_DATABASE"),
  os.Getenv("POSTGRES_HOST"),
  os.Getenv("POSTGRES_PORT"),
 )
 log.Println("Connecting to database with connection string:", dsn) 
 return dsn
}

func Database() *sql.DB {
 Db, err := sql.Open("postgres", setenv())
 if err != nil {
  log.Printf("connect db error, check .env => %v", err)
  return nil // Работает, даже если подключение не удалось
 }
 err = Db.Ping()
 if err != nil {
  log.Printf("db connected but can't ping: %v", err)
  return nil // Работает, даже если пинг не удался
 }

 return Db
}