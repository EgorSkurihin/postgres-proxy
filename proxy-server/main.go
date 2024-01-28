package main

import (
	"flag"

	pgspy "github.com/egorskurihin/postgres-proxy/proxy-server/pkg"
)

func main() {
	proxyAddrPtr := flag.String("proxy-addr", "", "string")
	dbAddrPtr := flag.String("db-addr", "", "string")
	storagePathPtr := flag.String("storage-path", "", "string")

	flag.Parse()

	pgspy.Start(pgspy.StartParams{
		Addr:        *proxyAddrPtr,
		DBAddr:      *dbAddrPtr,
		StoragePath: *storagePathPtr,
	})
}
