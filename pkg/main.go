package pkg

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// Standard way of handling db drivers
	_ "github.com/mattn/go-sqlite3"
)

var inAlert bool
var cfg Config

// Run the main loop of the application
func Run(logFile string, config Config) {
	cfg = config

	// setup the db
	os.Remove(cfg.DbPath)
	db, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrateDB(db)

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Force file to the EOF
	s := bufio.NewScanner(file)
	for s.Scan() {
	}

	for {
		duration, err := time.ParseDuration(cfg.IntervalDuration)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(duration)

		now := time.Now()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			section := sectionFromLine(line)
			insertSection(section, now.Unix(), db)
		}

		avgWindow, _ := time.ParseDuration(cfg.RollingAvgDuration)
		handleTrafficStats(now, avgWindow, db)

	}

}

func sectionFromLine(line string) string {
	logParts := strings.Split(line, "\"")
	requestParts := strings.Split(logParts[1], " ")
	siteParts := strings.Split(requestParts[1], "/")
	return "/" + siteParts[1]
}

func handleTrafficStats(now time.Time, duration time.Duration, db *sql.DB) {
	log.Println("===== Last", cfg.IntervalDuration, "=====")
	printTopN(now.Unix(), db)
	// Get the results of average traffic from the last 2 mins

	log.Println("===== General Stats =====")
	a := getAvgTraffic(now, duration, db)
	log.Println("Average requests per "+duration.String()+":", a)

	s := getSumTraffic(now, duration, db)
	log.Println("Sum of requests in last "+duration.String()+":", s)

	t := getTotalTraffic(db)
	log.Println("Total requests since starting:", t)

	newState := overThreshold(a, cfg.TrafficThreshold)

	if sendAlert(newState, inAlert) {
		inAlert = true
		log.Println("High traffic generated an alert - hits =", a, "triggered at", now.String())
	}

	if recoverAlert(newState, inAlert) {
		inAlert = false
		log.Println("High trafic alert recovered")
	}

}

func getAvgTraffic(now time.Time, duration time.Duration, db *sql.DB) float64 {
	queryString := "SELECT AVG(c) AS a " +
		"FROM (SELECT ts, COUNT(*) AS c " +
		"  FROM logs " +
		"  GROUP BY ts) " +
		"WHERE ts >= " + strconv.FormatInt(now.Add(-duration).Unix(), 10) + ";"

	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var a float64
		err = rows.Scan(&a)
		if err != nil {
			// TODO: Better handle the errors when there hasnt been any traffic in the last
			// duration.
			return 0.0
		}
		return a
	}
	return 0.0
}

func getSumTraffic(now time.Time, duration time.Duration, db *sql.DB) int {
	// We use a count because we have a row per logline
	queryString := "SELECT COUNT(*) AS s " +
		"FROM logs " +
		"WHERE ts >= " + strconv.FormatInt(now.Add(-duration).Unix(), 10) + ";"

	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var s int
		err = rows.Scan(&s)
		if err != nil {
			// TODO: Better handle the errors when there hasnt been any traffic in the last
			// duration.
			return 0
		}
		return s
	}
	return 0
}

func getTotalTraffic(db *sql.DB) int {
	queryString := "SELECT COUNT(*) AS c FROM logs;"

	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c int
		err = rows.Scan(&c)
		if err != nil {
			// TODO: Better handle the errors when there hasnt been any traffic in the last
			// duration.
			return 0
		}
		return c
	}
	return 0
}

func overThreshold(value float64, threshold float64) bool {
	if value >= threshold {
		return true
	}
	return false
}

func sendAlert(newState bool, currentState bool) bool {
	if newState && !currentState {
		return true
	}
	return false
}

func recoverAlert(newState bool, currentState bool) bool {
	if !newState && currentState {
		return true
	}
	return false
}

func migrateDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE logs (ts INTEGER, section TEXT)")
	if err != nil {
		log.Fatal(err)
	}
}

func insertSection(section string, ts int64, db *sql.DB) {
	sqlStatement := "INSERT INTO logs VALUES (" + strconv.FormatInt(ts, 10) + ", '" + section + "')"
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
}

func printTopN(ts int64, db *sql.DB) {

	queryString := "SELECT section, count(*) AS c FROM logs " +
		"WHERE ts = " + strconv.FormatInt(ts, 10) + " " +
		"GROUP BY section " +
		"ORDER BY c DESC LIMIT " + strconv.FormatInt(cfg.TopCount, 10) + ";"
	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c int
		var section string
		err = rows.Scan(&section, &c)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(c, section)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// Config stores the global config.
type Config struct {
	IntervalDuration   string
	TrafficThreshold   float64
	RollingAvgDuration string
	DbPath             string
	TopCount           int64
}
