# clm
A simple load manager written in Go. CLM auto scales your stateless application based on the flags supplied.

** This does not support sessions as of now.!!!!

### CMD flags
	$ clm -h

	  -app string
    	Path to binary of web server.
	  -ip string
	    	IPv4 address of machine. (default "127.0.0.1")
	  -max int
	    	Maximum number of connections per instance. clm will scale as needed. (default 100)
	  -port string
	    	Port to listen on. (default "8080")
	  -wait int
	    	Time to wait for server to start. (default 10 seconds)


### Go project setup
You Go web application should retrieve the port number to listen on from env. variable $PORT.

Example
	
	...
	port := ":defaultport"
	if envport := os.ExpandEnv("$PORT"); envport != "" {
		port = fmt.Sprintf(":%s", envport)
	}
	...
	log.Fatal(http.ListenAndServe(port, nil) )

### OS ulimits

Make sure you're ulimits are high to prevent socket read errors. On linux :

	$ ulimit -n 10000
