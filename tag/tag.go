package tag

import (
	"os"
	"log"
	"fmt"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DBNAME string
	DBUSERNAME string
	DBPWD string
	TABLENAME string
)

func init() {
	DBNAME = os.Getenv("DBNAME")
	DBUSERNAME = os.Getenv("DBUSERNAME")
	DBPWD = os.Getenv("DBPWD")
	TABLENAME = "tags"
}

type Tag struct {
	tag_name string
}


// get tag using tag name
// return search status dan tag_url
func Get(name, channel, guild string, scope int) (bool, string) {

	db := dbConn()
	defer db.Close()

	query := fmt.Sprintf(`
		SELECT tag_url FROM %s
		WHERE
		tag_name = ?
		`, TABLENAME)


	var url string
	var e error

	if scope == 2 {
		query = fmt.Sprintf(`
			%s
			tag_scope = ?;`, query)

		e = db.QueryRow(query, name, scope).Scan(&url)

	} else if scope == 1 {
		query = fmt.Sprintf(`
			%s
			tag_guild = ?
			tag_scope = ?;`, query)

		e = db.QueryRow(query, name, guild, scope).Scan(&url)
	} else {
		query = fmt.Sprintf(`
			%s
			tag_channel = ?
			tag_scope = ?;`, query)

		e = db.QueryRow(query, name, channel, scope).Scan(&url)
	}

	if e != nil {
		log.Println(e)
		return false, "nothing found"
	}

	return true, url
}

// add tag
// return add status
func Add(name, url, owner, channel, guild string, scope int) bool {

	db := dbConn()
	defer db.Close()

	query := fmt.Sprintf(`
		INSERT INTO %s(tag_name, tag_url, tag_owner, tag_channel, tag_guild, tag_scope)
		VALUES(?, ?, ?, ?, ?, ?);`, TABLENAME)

	insForm, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return false
	}

	insForm.Exec(name, url, owner, channel, guild, scope)

	return true
}

// delete tag
// return delete status
func Delete(name, url, owner, channel, guild string, scope int) bool {
	return false
}

// get a list of all available tags
func List(channel, guild string) string {

	db := dbConn()
	defer db.Close()

	query := fmt.Sprintf(`
		SELECT tag_name, tag_scope FROM %s
		WHERE
		(tag_scope = 0 and tag_channel = ?) OR 
		(tag_scope = 1 AND tag_guild = ?) OR
		tag_scope = 2;`,
		 TABLENAME)

	// execute query
	rows, err := db.Query(query, channel, guild)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var (
		channelTag []string
		guildTag []string
		publicTag []string
	)
	for rows.Next() {

		var (
			tag_name string
			tag_scope int
		)
		if err := rows.Scan(&tag_name, &tag_scope); err != nil {
			log.Println(err)
		}
		if tag_scope == 0 {
			channelTag = append(channelTag, tag_name)
		} else if tag_scope == 1 {
			guildTag = append(guildTag, tag_name)
		} else {
			publicTag = append(publicTag, tag_name)
		}

	}
	result := fmt.Sprintf("public: %s\nguild: %s\nchannel: %s",
						strings.Join(publicTag, ","),
						strings.Join(guildTag, ","),
						strings.Join(channelTag, ","))

	return result
}

// get metadata of a tag using name
// return search status and string of metadata
func Info(name, channel, guild string) (bool, string) {
	return true, "owner: yourmom"
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", DBUSERNAME, DBPWD, DBNAME))

	if err != nil {
		log.Println(err)
	}

	return db
}