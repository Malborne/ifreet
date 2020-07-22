package heimdallr

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
	// Register SQL driver
	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
)

//Infraction contains the reason and time for a user infraction
type Infraction struct {
	Reason string
	Time   time.Time
}

//Message contains the basic information about the message
type Message struct {
	messageID string
	channelID string
	content   string
	Time      time.Time
	userID    string
}

//Resource represents a learning resource
type Resource struct {
	ID      int
	Name    string
	Content string
	Tags    []string
}

var db *sql.DB

//OpenDb opens a connection to the database and creates the tables if they don't exist
func OpenDb(file string) error {
	var err error
	// db, err = sql.Open("sqlite3", file)
	db, err = sql.Open("postgres", file)

	if err != nil {
		return errors.Wrap(err, "opening database failed")
	}

	// CREATE TABLE IF NOT EXISTS whitelistedUsers (
	// 	id SERIAL PRIMARY KEY,
	// 	time_ timestamp,
	// 	user_id TEXT,
	// 	FOREIGN KEY(user_id) REFERENCES users(id)
	// );

	// dropTables := `DROP TABLE IF EXISTS users cascade;
	// 			  DROP TABLE IF EXISTS infractions cascade;
	// 			  DROP TABLE IF EXISTS mutedUsers cascade;
	// 			  DROP TABLE IF EXISTS resources cascade;
	// 			  DROP TABLE IF EXISTS resource_tags cascade;
	// 			  DROP TABLE IF EXISTS resource_tags_resources cascade;
	// 			  DROP TABLE IF EXISTS invites cascade;`

	createTableStatement := `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	username TEXT
);

CREATE TABLE IF NOT EXISTS infractions (
	id SERIAL PRIMARY KEY,
	reason TEXT,
	time_ timestamp,
	user_id TEXT,
	FOREIGN KEY(user_id) REFERENCES users(id)
);
  
CREATE TABLE IF NOT EXISTS mutedUsers (
	id SERIAL PRIMARY KEY,
	roleIDs TEXT,
	time_ timestamp,
	user_id TEXT,
	FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS archive (
	id SERIAL PRIMARY KEY,
	messageID TEXT,
	channelID TEXT,
	content TEXT,
	time_ timestamp,
	user_id TEXT,
	FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE OR REPLACE FUNCTION check_number_of_row()
RETURNS TRIGGER AS
$body$
BEGIN
    IF (SELECT count(*) FROM archive) > 1000 THEN 
	DELETE FROM archive WHERE id IN (SELECT id FROM archive ORDER BY time_ asc LIMIT 1); 
	END IF;
	RETURN NEW;
END;
$body$
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tr_check_number_of_row ON archive;

CREATE TRIGGER tr_check_number_of_row 
BEFORE INSERT ON archive
FOR EACH ROW EXECUTE PROCEDURE check_number_of_row();

CREATE TABLE IF NOT EXISTS resources (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS resource_tags (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS resource_tags_resources (
	resource_id INTEGER,
	resource_tag_id INTEGER,
	PRIMARY KEY(resource_id, resource_tag_id),
	FOREIGN KEY(resource_id) REFERENCES resources(id),
	FOREIGN KEY(resource_tag_id) REFERENCES resource_tags(id)
);
`
	// _, err = db.Exec(dropTables)
	// if err != nil {
	// 	return errors.Wrap(err, "deleting database tables failed")
	// }
	_, err = db.Exec(createTableStatement)
	return errors.Wrap(err, "creating database tables failed")
}

//CloseDb closes the database connection
func CloseDb() error {
	return errors.Wrap(db.Close(), "closing database failed")
}

//GetInfractions gets the list of infractions for a user
func GetInfractions(userID string) ([]Infraction, error) {
	var infractions []Infraction
	rows, err := db.Query(
		"SELECT reason, time_ FROM infractions WHERE user_id=$1 ORDER BY time_",
		userID,
	)
	if err != nil {
		return infractions, errors.Wrap(err, "fetching infractions failed")
	}

	for rows.Next() {
		var infractionReason string
		var infractionTime time.Time
		err = rows.Scan(&infractionReason, &infractionTime)
		if err != nil {
			return nil, errors.Wrap(err, "parsing infraction row failed")
		}
		infractions = append(infractions, Infraction{infractionReason, infractionTime})
	}

	if err = rows.Err(); err != nil {
		return infractions, errors.WithStack(err)
	}

	return infractions, nil
}

//AddInfraction adds an infraction for a user
func AddInfraction(user discordgo.User, infraction Infraction) error {
	err := AddUser(user)
	if err != nil {
		return errors.Wrap(err, "Adding user failed")
	}

	_, err = db.Exec("INSERT INTO infractions (reason, time_, user_id) VALUES ($1, $2, $3)",
		infraction.Reason, infraction.Time, user.ID)
	return errors.Wrap(err, "inserting infraction failed")
}

//RemoveInfraction removes an infraction for a user
func RemoveInfraction(timestamp time.Time) error {
	_, err := db.Query(
		"DELETE FROM infractions WHERE time_=$1",
		timestamp,
	)
	return errors.Wrap(err, "deleting infraction failed")
}

//AddtoArchive adds a message to the archive table
func AddtoArchive(user discordgo.User, m *discordgo.MessageCreate) error {
	err := AddUser(user)
	if err != nil {
		return errors.Wrap(err, "Adding user failed")
	}

	_, err = db.Exec("INSERT INTO archive (messageID, channelID, content, time_, user_id) VALUES ($1, $2, $3, $4, $5)",
		m.ID, m.ChannelID, m.Content, time.Now(), user.ID)
	return errors.Wrap(err, "inserting infraction failed")
}

//GetFromArchive gets a message from the archive table
func GetFromArchive(messageID string) (Message, error) {
	var message Message
	rows, err := db.Query(
		"SELECT channelID, time_, content, user_id FROM archive WHERE messageID=$1 ORDER BY time_",
		messageID,
	)
	if err != nil {
		return message, errors.Wrap(err, "fetching message from archive failed")
	}

	for rows.Next() {
		var channelID string
		var content string
		var messageTime time.Time
		var userID string
		err = rows.Scan(&channelID, &messageTime, &content, &userID)
		if err != nil {
			return message, errors.Wrap(err, "parsing infraction row failed")
		}
		message = Message{messageID, channelID, content, Time, userID}
	}

	if err = rows.Err(); err != nil {
		return message, errors.WithStack(err)
	}

	return message, nil
}

//RemovefromArchive Removes a message from the archive table
func RemovefromArchive(messageID string) error {
	_, err := db.Query(
		"DELETE FROM archive WHERE messageID=$1",
		messageID,
	)
	return errors.Wrap(err, "deleting message failed")
}

//AddMutedUser Adds a muted user to the list of users
func AddMutedUser(user discordgo.User, time time.Time, roleIDs string) error {
	err := AddUser(user)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO mutedUsers (roleIDs, time_, user_id) VALUES ($1, $2, $3)",
		roleIDs, time, user.ID)
	return errors.Wrap(err, "muting user failed")
}

//AddWhitelistedUser Adds a user to the whitelist to be able to post links
// func AddWhitelistedUser(user discordgo.User, time time.Time) error {
// 	err := AddUser(user)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("INSERT INTO mutedUsers (time_, user_id) VALUES ($1, $2)",
// 		time, user.ID)
// 	return errors.Wrap(err, "muting user failed")
// }

//IsuserWhitelisted checks if a user is whitelisted
// func IsuserWhitelisted(userID string) (bool, error) {
// 	_, err := db.Exec("SELECT EXISTS (SELECT 1 from whitelistedUsers where user_id=$1)", userID)
// 	if err != nil {
// 		return false, errors.Wrap(err, "fetching infractions failed")
// 	}
// }

//GetMutedUserRoles retrieves the muted roles of a muted member
func GetMutedUserRoles(userID string) ([]string, error) {
	var roles []string
	var roleIDs string
	rows, err := db.Query(
		"SELECT roleIDs, time_ FROM mutedUsers WHERE user_id=$1 ORDER BY time_",
		userID,
	)
	if err != nil {
		return roles, errors.Wrap(err, "fetching infractions failed")
	}

	for rows.Next() {
		var mutedTime time.Time
		err = rows.Scan(&roleIDs, &mutedTime)
		if err != nil {
			return nil, errors.Wrap(err, "parsing infraction row failed")
		}

		// roleIDs = append(infractions, Infraction{infractionReason, infractionTime})
	}
	roles = strings.Split(roleIDs, ",")

	if err = rows.Err(); err != nil {
		return roles, errors.WithStack(err)
	}

	return roles, nil
}

//RemoveMutedUser Removes a user from the database after being unmuted
func RemoveMutedUser(userID string) error {
	_, err := db.Query(
		"DELETE FROM mutedUsers WHERE user_id=$1",
		userID,
	)
	return errors.Wrap(err, "deleting user failed")
}

//AddUser adds a user or updates the username if it already exists
func AddUser(user discordgo.User) error {
	// ON CONFLICT (id) DO UPDATE SET username=$2
	// ON CONFLICT (id) IGNORE
	_, err := db.Exec("INSERT INTO users (id,username) VALUES ($1, $2) ON CONFLICT DO NOTHING", user.ID, user.Username)
	if err != nil {
		return errors.Wrap(err, "inserting user failed")
	}
	_, err = db.Exec("UPDATE users SET username=$1 WHERE id=$2", user.Username, user.ID)
	return errors.Wrap(err, "updating user failed")
}

//GetResourceByID gets a resource from the database by ID
func GetResourceByID(id int) (*Resource, error) {
	rows, err := db.Query(
		"SELECT resources.id, resources.name, content, resource_tags.name AS tag"+
			" FROM resources"+
			" 	LEFT JOIN resource_tags_resources ON resource_id = resources.id"+
			"	LEFT JOIN resource_tags ON resource_tag_id = resource_tags.id"+
			" WHERE resources.id = $1",
		id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "getting resource by id failed")
	}
	resources, err := getResources(rows)
	if err != nil {
		return nil, err
	}
	return resources[0], nil
}

//GetResourceByName gets a resource frm the database by name
func GetResourceByName(name string) (*Resource, error) {
	rows, err := db.Query(
		"SELECT resources.id, resources.name, content, resource_tags.name AS tag"+
			" FROM resources"+
			" 	LEFT JOIN resource_tags_resources ON resource_id = resources.id"+
			"	LEFT JOIN resource_tags ON resource_tag_id = resource_tags.id"+
			" WHERE LOWER(resources.name) LIKE LOWER('%' || $1 || '%')",
		name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "getting resource by name failed")
	}
	resources, err := getResources(rows)
	if err != nil {
		return nil, err
	}
	return resources[0], nil
}

//SearchResources searches the database for resources matching the search terms
func SearchResources(searchTerms []string) ([]*Resource, error) {
	query := "SELECT resources.id, resources.name, content, resource_tags.name AS tag" +
		" FROM resources" +
		" 	LEFT JOIN resource_tags_resources ON resource_id = resources.id" +
		"	LEFT JOIN resource_tags ON resource_tag_id = resource_tags.id" +
		" WHERE"
	var queryTerms []interface{}
	for i, searchTerm := range searchTerms {
		query += " LOWER(resources.name) LIKE LOWER('%' || ? || '%')" +
			"		OR LOWER(resources.content) LIKE LOWER('%' || ? || '%')" +
			"		OR LOWER(resource_tags.name) LIKE LOWER('%' || ? || '%')"
		if i < len(searchTerms)-1 {
			query += " OR"
		}
		for i := 0; i < 3; i++ {
			queryTerms = append(queryTerms, searchTerm)
		}
	}
	rows, err := db.Query(query, queryTerms...)
	if err != nil {
		return nil, errors.Wrap(err, "searching resources failed")
	}
	resources, err := getResources(rows)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func getResources(rows *sql.Rows) ([]*Resource, error) {
	var resources []*Resource
	resourcesByID := make(map[int]Resource)
	for rows.Next() {
		var id int
		var name, content string
		var tag sql.NullString
		err := rows.Scan(&id, &name, &content, &tag)
		if err != nil {
			return nil, errors.Wrap(err, "parsing resource row failed")
		}
		resource, ok := resourcesByID[id]
		if !ok {
			resource = Resource{
				ID:      id,
				Name:    name,
				Content: content,
			}
			resourcesByID[id] = resource
			resources = append(resources, &resource)
		}
		if tag.Valid {
			resource.Tags = append(resource.Tags, tag.String)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return resources, nil
}

//AddResource adds a resource to the database
func AddResource(resource Resource) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, errors.Wrap(err, "starting transaction failed")
	}

	result, err := tx.Exec("INSERT INTO resources (name, content) VALUES ($1, $2)",
		resource.Name, resource.Content)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Wrap(rollbackErr, "rollback failed")
		} else {
			err = errors.Wrap(err, "inserting resource failed")
		}
		return -1, err
	}

	resourceID, err := result.LastInsertId()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Wrap(rollbackErr, "rollback failed")
		} else {
			err = errors.Wrap(err, "getting id of inserted resource failed")
		}
		return -1, err
	}

	for _, tag := range resource.Tags {
		_, err = tx.Exec("INSERT OR REPLACE INTO resource_tags (name) VALUES ($1)", tag)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "rollback failed")
			} else {
				err = errors.Wrap(err, "inserting resource tag failed")
			}
			return -1, err
		}

		result := tx.QueryRow("SELECT id FROM resource_tags WHERE name = $1", tag)
		var tagID int64
		if err = result.Scan(&tagID); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "rollback failed")
			} else {
				err = errors.Wrap(err, "getting id of resource tag failed")
			}
			return -1, err
		}

		_, err = tx.Exec("INSERT INTO resource_tags_resources (resource_id, resource_tag_id) VALUES ($1, $2)", resourceID, tagID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "rollback failed")
			} else {
				err = errors.Wrap(err, "inserting relationship between resource and resource tag failed")
			}
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Wrap(rollbackErr, "rollback failed")
		} else {
			err = errors.Wrap(err, "committing transaction failed")
		}
		return -1, err
	}

	return resourceID, nil
}
