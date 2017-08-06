package logs

import (
	"log"
	"database/sql"

    _ "github.com/mattn/go-sqlite3"

    //"GoHole/config"
)

var instance *sql.DB = nil

type QueryLog struct {
    Id int
    ClientIp string
    Domain string
    Cached int
    Timestamp int64
}

type ClientLog struct {
    ClientIp string
    Queries int
}

func GetInstance() *sql.DB {
    if instance == nil {
    	var err error = nil
    	instance, err = sql.Open("sqlite3", "./gohole.db")
    	if err != nil{
    		log.Fatal(err)
    	}
    }

    return instance
}

func SetupDB(){
	sqlCmd := `
	create table if not exists queries (id integer not null primary key autoincrement, clientip text, domain text, cached integer, timestamp integer);
	`

	_, err := GetInstance().Exec(sqlCmd)
	if err != nil{
		log.Fatal(err)
	}
}

func AddQuery(clientip, domain string, cached int, timestamp int64) (error){
	tx, err := GetInstance().Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO queries(clientip, domain, cached, timestamp) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(clientip, domain, cached, timestamp)
	if err != nil {
		return err
	}
	err = tx.Commit()

	return err
}

func GetQueriesByClientIp(clientIp string) ([]QueryLog, error){
	result := []QueryLog{}
	rows, err := GetInstance().Query("select id, clientip, domain, cached, timestamp from queries where clientip='"+ clientIp +"'")
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, cached int
		var timestamp int64
		var clientip, domain string
		err = rows.Scan(&id, &clientip, &domain, &cached, &timestamp)
		if err != nil {
			return result, err
		}
		
		l := QueryLog{
			Id: id,
			ClientIp: clientip,
			Domain: domain,
			Cached: cached,
			Timestamp: timestamp,
		}
		result = append(result, l)
	}
	err = rows.Err()
	return result, err
}

func GetQueriesByDomain(domain string) ([]QueryLog, error){
	result := []QueryLog{}
	rows, err := GetInstance().Query("select id, clientip, domain, cached, timestamp from queries where domain LIKE '%"+ domain +"%'")
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, cached int
		var timestamp int64
		var clientip, domain string
		err = rows.Scan(&id, &clientip, &domain, &cached, &timestamp)
		if err != nil {
			return result, err
		}
		
		l := QueryLog{
			Id: id,
			ClientIp: clientip,
			Domain: domain,
			Cached: cached,
			Timestamp: timestamp,
		}
		result = append(result, l)
	}
	err = rows.Err()
	return result, err
}

func GetClients() ([]ClientLog, error){
	result := []ClientLog{}
	rows, err := GetInstance().Query("select clientip, count(*) from queries GROUP BY clientip")
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var clientip string
		var queries int
		err = rows.Scan(&clientip, &queries)
		if err != nil {
			return result, err
		}

		c := ClientLog{
			ClientIp: clientip,
			Queries: queries,
		}
		result = append(result, c)
	}
	err = rows.Err()
	return result, err
}

func Flush() (error){
	sqlCmd := `
	delete from queries;
	`

	_, err := GetInstance().Exec(sqlCmd)

	return err
}


