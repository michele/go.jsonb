package jsonb_test

import (
	jsonb "github.com/michele/go.jsonb"

	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const (
	DB_USER     = "jsonb_user"
	DB_PASSWORD = "jsonb_pass"
	DB_NAME     = "jsonb_test"
)

var theTestMap = map[string]interface{}{
	"key":         "value",
	"another_key": "another_value",
	"nested_key": map[string]interface{}{
		"child_key": "child_value",
	},
}

var theTestString = `{"field":{"another_key":"another_value","key":"value","nested_key":{"child_key":"child_value"}}}`

var theTestJSONB = jsonb.JSONB{
	"key":         "value",
	"another_key": "another_value",
	"nested_key": map[string]interface{}{
		"child_key": "child_value",
	},
}

type TestStruct struct {
	Field jsonb.JSONB `json:"field"`
}

func TestJSONBMarshalling(t *testing.T) {
	theTest := TestStruct{
		Field: theTestMap,
	}

	b, err := json.Marshal(theTest)

	assert.Nil(t, err)

	assert.Equal(t, theTestString, string(b))
}

func TestJSONBUnmarshalling(t *testing.T) {
	theTest := TestStruct{}

	err := json.Unmarshal([]byte(theTestString), &theTest)

	assert.Nil(t, err)

	assert.Equal(t, jsonb.JSONB(theTestMap), theTest.Field)
}

func TestDBInteraction(t *testing.T) {
	// psql -c 'CREATE USER jsonb_user WITH PASSWORD 'jsonb_pass';' -U your_user -hlocalhost
	// psql -c 'GGRANT ALL ON SEQUENCE public.test_id_seq TO jsonb_user; GRANT ALL ON TABLE public.test TO jsonb_user;' -U your_user -hlocalhost -djsonb_test
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	var lastInsertId int
	err = db.QueryRow("INSERT INTO test(field) VALUES($1) returning id;", theTestJSONB).Scan(&lastInsertId)
	assert.Nil(t, err)

	checkErr(err)

	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM test WHERE id=$1", lastInsertId)

	assert.Nil(t, err)
	checkErr(err)

	for rows.Next() {
		var id int
		var field jsonb.JSONB
		err = rows.Scan(&id, &field)
		assert.Nil(t, err)
		checkErr(err)
		assert.Equal(t, theTestJSONB, field)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
