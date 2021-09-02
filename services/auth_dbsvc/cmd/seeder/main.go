package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service/db"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var filePath = flag.String("file", "./scripts/seeding/users.csv", "seeding file for the users table")

type seeder struct {
	connection *pg.DB
	logger     *logrus.Logger
}

func NewSeeder(connection *pg.DB, logger *logrus.Logger) *seeder {
	return &seeder{
		connection: connection,
		logger:     logger,
	}
}

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	// prepare pg connection
	pgClient := db.NewPostgresClientFromConfig(cfg)
	connection := pgClient.GetConnection()
	defer connection.Close()

	seeder := NewSeeder(connection, cfg.Logger)

	err = seeder.csvParse(*filePath, "users")
	if err != nil {
		cfg.Logger.Fatal(err, "Could not seed table!")
	}

	cfg.Logger.Info("table seeded successfully!")
}

func (s *seeder) csvParse(filePath string, tableName string) (err error) {

	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file error: %v", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	idx := 0
	fieldNames := []string{"role", "username", "password"}

	for sc.Scan() {
		txt := sc.Text()
		s.logger.Info(txt)

		if idx == 0 {
			csvReader := csv.NewReader(strings.NewReader(txt))
			var record []string
			record, err = csvReader.Read()
			if err != nil {
				return
			}
			err = s.sqlInsert(fieldNames, record, tableName)
			if err != nil {
				s.logger.Fatal(err)
			}
		}
		idx++
	}
	if err = sc.Err(); err != nil {
		return
	}
	return
}

func (s *seeder) sqlInsert(cols []string, record []string, table string) (err error) {

	colsStr := "(role, username, password)"

	valStr := "( "
	for idx, val := range record {

		if val != "null" {
			val = fmt.Sprintf("'%v'", val)
		}

		if idx < len(cols)-1 {
			valStr += fmt.Sprintf("%v, ", val)
		} else {
			//password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost+1)
			if err != nil {
				return err
			}
			valStr += fmt.Sprintf("'%v' )", string(hashedPassword))
		}
	}
	qry := fmt.Sprintf(
		`INSERT INTO "%s" %s VALUES %s`,
		table, colsStr, valStr)

	s.logger.Info("Executing Query: ", qry)
	res, err := s.connection.Exec(qry)

	if err != nil {
		return fmt.Errorf("not inserted correctly, result is nil: %v", err)
	}

	rowInserted := res.RowsAffected()
	if rowInserted != 1 {
		err = errors.New("not inserted correctly, affected rows is not 1")
	}

	return
}
