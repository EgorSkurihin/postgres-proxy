package main

import (
	"flag"

	web "github.com/egorskurihin/postgres-proxy/web-server/pkg"
)

func main() {
	addrPtr := flag.String("addr", "", "string")
	usersFilePath := flag.String("users-file-path", "", "string")
	dataFilePath := flag.String("storage-path", "", "string")

	flag.Parse()

	web.Start(web.StartParams{
		Addr:             *addrPtr,
		DataFilePath:     *dataFilePath,
		AuthUserFilePath: *usersFilePath,
	})
}

/* [
    {
        "login": "admin",
        "password": "$2a$14$Wr1Kuv0LDW5dl33rAJ7Vle5UGmmgLEiBsKuc5khCfqaREIyw.azAe"
    }
]
*/
