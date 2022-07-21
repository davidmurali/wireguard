package device

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// structure of database table rows
type allowedIPsRow struct {
	entries_id    string
	dest_ip       string
	dest_protocol string
	dest_port     string
	is_active     int
}

//data of db configuration xml
type DBDetails struct {
	DriverName string `xml:"DriverName"`
	UserName   string `xml:"UserName"`
	Password   string `xml:"Password"`
	IPAddress  string `xml:"IPAddress"`
	Port       string `xml:"Port"`
	DBName     string `xml:"DBName"`
	DBTable    string `xml:"DBTable"`
}

// from DB to nested map of AddressToProtocolToPort
type DestAddressMap map[string]DestProtocolMap
type DestProtocolMap map[string]DestPortMap
type DestPortMap map[string]string

var destAddressMap DestAddressMap = DestAddressMap{}

var db *sql.DB

//Adding rows to the nested map of AddressToProtocolToPort
func AddToAddressMap(address string, port string, protocol string) {
	protocol = strings.ToLower(protocol)
	protocolMap, exists := destAddressMap[address]
	if !exists {
		destAddressMap[address] = DestProtocolMap{}
		protocolMap = destAddressMap[address]
	}
	portMap, exists := protocolMap[protocol]
	if !exists {
		protocolMap[protocol] = DestPortMap{}
		portMap = protocolMap[protocol]
	}
	portMap[port] = port
}


// check and allow the incoming packets are within the bounded protocols
func CanByPassEncrypt(address string, port string, protocol string) bool {
	protocol = strings.ToLower(protocol)
	protocolMap, exists := destAddressMap[address]
	if !exists {
		protocolMap, exists = destAddressMap["*"]
	}
	if !exists {
		return false
	}
	portMap, exists := protocolMap[protocol]
	if !exists {
		portMap, exists = protocolMap["*"]
	}
	if !exists {
		return false
	}
	_, exists = portMap[port]
	if !exists {
		_, exists = portMap["*"]
	}
	return exists
}

//Connect DB and store data in nested map
func connectDBAndGetData() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dbFilePath := filepath.Dir(exe) + "//dbconf.xml"
	fileInfo, err := os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("please place the DB configuration file")
			fmt.Println(fileInfo)
			return
		}
	}
	data, _ := ioutil.ReadFile(dbFilePath)
	dbdetails := &DBDetails{}
	_ = xml.Unmarshal([]byte(data), &dbdetails)
	connectionString := dbdetails.UserName + ":" + dbdetails.Password + "@tcp(" + dbdetails.IPAddress + ":" + dbdetails.Port + ")/" + dbdetails.DBName
	db, err = sql.Open(dbdetails.DriverName, connectionString)
	if err != nil {
		fmt.Println("DBDetails are missing or wrong")
		return
	}
	StoreDBInLocalMap(dbdetails.DBTable)
}

func CloseDB() {
	db.Close()
}

//storing to local map
func StoreDBInLocalMap(dbtablename string) {
	defer CloseDB()
	queryString := "SELECT * FROM " + dbtablename + " Where is_active = 1"
	rows, err := db.Query(queryString)
	if err != nil {
		fmt.Println("DBDetails are missing or wrong in dbconf.xml")
		return
	}

	for rows.Next() {
		var row allowedIPsRow
		err := rows.Scan(&row.entries_id, &row.dest_ip, &row.dest_protocol, &row.dest_port, &row.is_active)
		if err != nil {
			fmt.Println("DB row is not retrieved")
		} else {
			AddToAddressMap(row.dest_ip, row.dest_protocol, row.dest_port)
		}
	}
}
