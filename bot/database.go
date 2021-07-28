package heimdallr

import (
	"database/sql"
	"os"
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
	ID     string
	Reason string
	Time   time.Time
}

//IsolatedUser contains the basic information about an isolated User
type IsolatedUser struct {
	UserID    string
	StartTime time.Time
	EndTime   time.Time
	RoleIDs   []string
}

//Student contains the information about each student in the Hifz circles
type Student struct {
	ID        string
	Circle    string
	SheetLink string
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

	// dropTables := `DROP TABLE IF EXISTS archive cascade;`
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

CREATE TABLE IF NOT EXISTS isolatedUsers (
	id SERIAL PRIMARY KEY,
	roleIDs TEXT,
	start_time timestamp,
	end_time timestamp,
	user_id TEXT,
	FOREIGN KEY(user_id) REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS students (
	id SERIAL PRIMARY KEY,
	user_id TEXT,
	circle TEXT,
	sheetLink TEXT,
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
    IF (SELECT count(*) FROM archive) >= 1000 THEN 
	DELETE FROM archive WHERE id IN (SELECT id FROM archive ORDER BY time_ asc LIMIT 1); 
	END IF;
	RETURN NEW;
END;
$body$
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tr_check_number_of_row ON archive;

CREATE TRIGGER tr_check_number_of_row 
BEFORE INSERT ON archive
FOR EACH STATEMENT EXECUTE PROCEDURE check_number_of_row();

CREATE TABLE IF NOT EXISTS resources (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS new_channels (
	id SERIAL PRIMARY KEY,
	channel_ID TEXT NOT NULL,
	user_ID TEXT NOT NULL
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
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(19)

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
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT id, reason, time_ FROM infractions WHERE user_id=$1 ORDER BY time_",
		userID,
	)
	if err != nil {
		return infractions, errors.Wrap(err, "fetching infractions failed")
	}

	for rows.Next() {
		var infractionID string
		var infractionReason string
		var infractionTime time.Time
		err = rows.Scan(&infractionID, &infractionReason, &infractionTime)
		if err != nil {
			return nil, errors.Wrap(err, "parsing infraction row failed")
		}
		infractions = append(infractions, Infraction{infractionID, infractionReason, infractionTime})
	}

	if err = rows.Err(); err != nil {
		return infractions, errors.WithStack(err)
	}
	rows.Close()
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
func RemoveInfraction(ID string) error {
	_, err := db.Query(
		"DELETE FROM infractions WHERE id::text = $1::text",
		ID,
	)
	return errors.Wrap(err, "deleting infraction failed")
}

//AddStudent adds a new student to the database
func AddStudent(user discordgo.User, circle string, sheetLink string) error {
	err := AddUser(user)
	if err != nil {
		return errors.Wrap(err, "Adding user failed")
	}

	_, err = db.Exec("INSERT INTO students (user_id, circle, sheetLink) VALUES ($1, $2, $3)",
		user.ID, circle, sheetLink)
	return errors.Wrap(err, "adding student failed")
}

//GetStudent retrieve a student information from the database
func GetStudent(userID string) (Student, error) {
	var student Student
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT user_id, circle, sheetLink FROM students WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return student, errors.Wrap(err, "getting student failed")
	}
	for rows.Next() {
		var userID string
		var circle string
		var sheetLink string

		err := rows.Scan(&userID, &circle, &sheetLink)
		if err != nil {
			return student, err
		}
		student = Student{userID, circle, sheetLink}

	}
	rows.Close()
	return student, nil
}

//AddNewChannel adds a new channel to the database
func AddNewChannel(userID string, channelID string) error {
	_, err := db.Exec("INSERT INTO new_channels (user_ID,channel_ID) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, channelID)
	if err != nil {
		return errors.Wrap(err, "inserting new channel failed")
	}
	return nil
}

//RemoveInfraction removes an infraction for a user
func RemoveNewChannel(userID string) error {
	_, err := db.Query(
		"DELETE FROM new_channels WHERE user_ID::text = $1::text",
		userID,
	)
	return errors.Wrap(err, "deleting channel failed")
}

//GetStudents retrieve all the students information of a cirtain circle from the database
func GetStudents(circleName string) ([]Student, error) {
	var students []Student
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT user_id, circle, sheetLink FROM students WHERE circle=$1",
		circleName,
	)
	if err != nil {
		return students, errors.Wrap(err, "getting student failed")
	}
	for rows.Next() {
		var userID string
		var circle string
		var sheetLink string

		err := rows.Scan(&userID, &circle, &sheetLink)
		if err != nil {
			return students, err
		}
		students = append(students, Student{userID, circle, sheetLink})

	}
	rows.Close()
	return students, nil
}

//RemoveStudent removes a student from the database
func RemoveStudent(userID string) error {
	_, err := db.Query(
		"DELETE FROM students WHERE user_id = $1",
		userID,
	)
	return errors.Wrap(err, "deleting student failed")
}

//GetnewChannel gets a channel ID from the database using the user ID
func GetnewChannel(userID string) (string, error) {
	var ChannelID string = ""
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT channel_ID FROM new_channels WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return ChannelID, errors.Wrap(err, "getting channel ID failed")
	}
	for rows.Next() {
		err := rows.Scan(&ChannelID)
		if err != nil {
			return ChannelID, err
		}
	}
	rows.Close()
	return ChannelID, nil
}

//GetAllnewChannels gets the IDs of all the new channels
func GetAllnewChannels() ([]string, error) {
	var ChannelIDs []string
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT channel_ID FROM new_channels")
	if err != nil {
		return ChannelIDs, errors.Wrap(err, "getting channel IDs failed")
	}
	for rows.Next() {
		var ChannelID string
		err := rows.Scan(&ChannelID)
		if err != nil {
			return ChannelIDs, err
		}
		ChannelIDs = append(ChannelIDs, ChannelID)
	}
	rows.Close()
	return ChannelIDs, nil
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

	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT channelID, time_, content, user_id FROM archive WHERE messageID=$1 ORDER BY time_",
		messageID,
	)
	if err != nil {
		return message, errors.Wrap(err, "fetching message failed")
	}
	for rows.Next() {
		var channelID string
		var content string
		var messageTime time.Time
		var userID string
		err := rows.Scan(&channelID, &messageTime, &content, &userID)
		if err != nil {
			return message, err
		}
		message = Message{messageID, channelID, content, messageTime, userID}

	}
	rows.Close()
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

//AddIsolatedUser Adds an isolated user to the list of users
func AddIsolatedUser(user discordgo.User, start_time time.Time, end_time time.Time, roleIDs string) error {
	err := AddUser(user)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO isolatedUsers (roleIDs, start_time, end_time, user_id) VALUES ($1, $2, $3, $4)",
		roleIDs, start_time, end_time, user.ID)
	return errors.Wrap(err, "isolating user failed")
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
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
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
	rows.Close()
	return roles, nil
}

//GetIsolatedUserRoles retrieves the roles of an isolated member
func GetIsolatedUserRoles(userID string) ([]string, error) {
	var roles []string
	var roleIDs string
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT roleIDs, start_time, end_time FROM isolatedUsers WHERE user_id=$1 ORDER BY start_time",
		userID,
	)
	if err != nil {
		return roles, errors.Wrap(err, "fetching infractions failed")
	}

	for rows.Next() {
		var startTime time.Time
		var endTime time.Time

		err = rows.Scan(&roleIDs, &startTime, &endTime)
		if err != nil {
			return nil, errors.Wrap(err, "getting row failed")
		}

	}
	roles = strings.Split(roleIDs, ",")

	if err = rows.Err(); err != nil {
		return roles, errors.WithStack(err)
	}
	rows.Close()
	return roles, nil
}

//GetIsolatedEndTime retrieves the end time of an isolated user
func GetIsolatedEndTime(userID string) (time.Time, error) {
	var end_time time.Time
	var roleIDs string
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT roleIDs, start_time, end_time FROM isolatedUsers WHERE user_id=$1 ORDER BY start_time",
		userID,
	)
	if err != nil {
		return end_time, errors.Wrap(err, "fetching infractions failed")
	}

	for rows.Next() {
		var startTime time.Time

		err = rows.Scan(&roleIDs, &startTime, &end_time)
		if err != nil {
			return end_time, errors.Wrap(err, "getting row failed")
		}

	}

	if err = rows.Err(); err != nil {
		return end_time, errors.WithStack(err)
	}
	rows.Close()
	return end_time, nil
}

//GetAllIsolatedUsers retrieves a list of isolated users
func GetAllIsolatedUsers() ([]IsolatedUser, error) {

	var roleIDs string
	var userID string
	var startTime time.Time
	var endTime time.Time
	var isolatedUsers []IsolatedUser
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
	rows, err := db.Query(
		"SELECT user_id, roleIDs, start_time, end_time FROM isolatedUsers ORDER BY start_time")
	if err != nil {
		return isolatedUsers, errors.Wrap(err, "fetching isolated users failed")
	}

	for rows.Next() {

		err = rows.Scan(&userID, &roleIDs, &startTime, &endTime)
		if err != nil {
			return isolatedUsers, errors.Wrap(err, "getting row failed")
		}

		roles := strings.Split(roleIDs, ",")

		// isolatedUser := IsolatedUser{
		// 	userID:    userID,
		// 	roleIDs:   roles,
		// 	startTime: startTime,
		// 	endTime:   endTime,
		// }
		isolatedUsers = append(isolatedUsers, IsolatedUser{userID, startTime, endTime, roles})
	}

	if err = rows.Err(); err != nil {
		return isolatedUsers, errors.WithStack(err)
	}
	rows.Close()
	return isolatedUsers, nil
}

//RemoveMutedUser Removes a user from the database after being unmuted
func RemoveMutedUser(userID string) error {
	_, err := db.Query(
		"DELETE FROM mutedUsers WHERE user_id=$1",
		userID,
	)
	return errors.Wrap(err, "deleting user failed")
}

//RemoveIsolatedUser Removes a user from the database after being unisolated
func RemoveIsolatedUser(userID string) error {
	_, err := db.Query(
		"DELETE FROM isolatedUsers WHERE user_id=$1",
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
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
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
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
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
	rows.Close()
	return resources[0], nil
}

//SearchResources searches the database for resources matching the search terms
func SearchResources(searchTerms []string) ([]*Resource, error) {
	if db.Stats().OpenConnections >= db.Stats().MaxOpenConnections || db.Stats().InUse >= db.Stats().MaxOpenConnections { //closes the connection pool and opens a new one to clear out the connections
		db.Close()
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(db.Stats().MaxOpenConnections)
	}
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
	rows.Close()
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
