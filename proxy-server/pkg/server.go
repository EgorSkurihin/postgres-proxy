package pgspy

type StartParams struct {
	Addr        string
	DBAddr      string
	StoragePath string
}

func Start(startParams StartParams) {
	pgHost := startParams.DBAddr
	proxyHost := startParams.Addr

	proxy := NewProxy(pgHost, proxyHost, startParams.StoragePath)
	proxy.Start()
}
