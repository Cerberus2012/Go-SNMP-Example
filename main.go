package main

import (
	"log"
	"time"

	"github.com/gosnmp/gosnmp"
)

var sensorOids = []string{
	"1.3.6.1.4.1.20916.1.9.1.1.1.1",
	"1.3.6.1.4.1.20916.1.9.1.1.1.2",
	"1.3.6.1.4.1.20916.1.9.1.1.1.3",
}

var sensorOidsMap = map[int]string{
	0: "internalSensorC",
	1: "internalSensorF",
	2: "sensorLabel",
}

func main() {
	log.Println("Starting main")

	var snmp gosnmp.GoSNMP
	snmp.Target = "192.168.0.2"
	snmp.Version = gosnmp.Version1
	snmp.Port = 161
	snmp.Community = "public"
	snmp.Retries = 1
	snmp.Timeout = time.Duration(30) * time.Second
	snmp.MaxOids = 60

	log.Println("Max OIDS:", snmp.MaxOids)
	log.Printf("Target Port: %d", snmp.Port)
	log.Printf("Community: %s", snmp.Community)
	log.Printf("SNMP Version: %s", snmp.Version)
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc

	log.Println("Connecting to target")
	err := snmp.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
		return
	} else {
		log.Printf("Connection: %s", snmp.Conn)
		log.Println("Connection succesful")
	}

	for i := 0; i < 50; i++ {
		log.Println("Performing SNMP GET")
		// oids := []string{
		// 	"1.3.6.1.4.1.20916.1.9.1.1.1.1",
		// 	"1.3.6.1.4.1.20916.1.9.1.1.1.2",
		// 	"1.3.6.1.4.1.20916.1.9.1.1.1.3",
		// }
		result, err2 := snmp.GetNext(sensorOids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			log.Fatalf("Get() err: %v", err2)
		} else {
			log.Println("GET successful")
		}

		log.Println("GET Results:")
		// log.Println(result)
		for i, variable := range result.Variables {
			// fmt.Printf("%d: oid: %s \n", i, variable.Name)
			// log.Println(variable.Value)
			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch variable.Type {
			case gosnmp.OctetString:
				log.Printf("index: %d, Label: %s, Value: %s\n", i, sensorOidsMap[i], string(variable.Value.([]byte)))
			default:
				// ... or often you're just interested in numeric values.
				// ToBigInt() will return the Value as a BigInt, for plugging
				// into your calculations.
				log.Printf("Index: %d, Label: %s, Value: %d\n", i, sensorOidsMap[i], gosnmp.ToBigInt(variable.Value))
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	snmp.Conn.Close()
}
