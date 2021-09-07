package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var filePath = flag.String("file", "./scripts/seeding/users.csv", "seeding file for the users table")

type seeder struct {
	connection *pg.DB
	logger     log.Logger
}

func NewSeeder(connection *pg.DB, logger log.Logger) *seeder {
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
		level.Error(cfg.Logger).Log("err", errors.Wrap(err, "Could not seed table!"))
		os.Exit(1)
	}

	level.Info(cfg.Logger).Log("msg", "table seeded successfully!")
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
		level.Info(s.logger).Log("msg", txt)

		if idx == 0 {
			csvReader := csv.NewReader(strings.NewReader(txt))
			var record []string
			record, err = csvReader.Read()
			if err != nil {
				return
			}
			err = s.sqlInsert(fieldNames, record, tableName)
			if err != nil {
				level.Error(s.logger).Log("err", err)
				os.Exit(1)
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
			//TODO refactor
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(val+auth.PWDSALT), bcrypt.DefaultCost+1)
			if err != nil {
				return err
			}
			valStr += fmt.Sprintf("'%v' )", string(hashedPassword))
		}
	}
	qry := fmt.Sprintf(
		`INSERT INTO "%s" %s VALUES %s`,
		table, colsStr, valStr)

	level.Info(s.logger).Log("msg", fmt.Sprintf("Executing Query: %v", qry))
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
