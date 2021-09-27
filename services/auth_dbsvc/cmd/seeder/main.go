// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE. 
//
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
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

	//table name hardcoded on purpose, we don't want it to be in config files and environment variables...
	err = seeder.csvParse(*filePath, "seeding_tmp")
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
		level.Info(s.logger).Log("msg", txt[:6]+"...")

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

	//This table will be used only once for seeding in order to keep hashing responsility where it belongs
	createTable := `
	CREATE TABLE if not exists seeding_tmp(
		 role VARCHAR(255),
		 username VARCHAR(255) NOT NULL UNIQUE,
		 password TEXT NOT NULL
	)`

	_, err = s.connection.Exec(createTable)
	if err != nil {
		return fmt.Errorf("cannot create table seeding_tmp: %v", err)
	}

	qry := fmt.Sprintf(`INSERT INTO %s (role, username, password) VALUES (?, ?, ?)`, table)

	level.Info(s.logger).Log("msg", fmt.Sprintf("Executing Query (obfuscated): %v", qry))
	res, err := s.connection.Exec(qry,
		record[0],
		record[1],
		//password
		//we could hash it here but it would be a mistake for several reasons.
		//better to use a temporary table that we remove after the application starts for the first time
		record[2],
	)

	if err != nil {
		return fmt.Errorf("not inserted correctly, result is nil: %v", err)
	}

	rowInserted := res.RowsAffected()
	if rowInserted != 1 {
		err = errors.New("not inserted correctly, affected rows is not 1")
	}

	return
}
