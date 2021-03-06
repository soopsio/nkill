# nkill

Kills all processes listening on the given TCP ports.

#### Install

You need go installed and GOBIN in your PATH. Once that is done, run the command:

```bash
   $ go get -u github.com/soopsio/nkill
```

#### Usage

To kill any process listening to the port 8080:

```go
func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Kills all processes listening on the given TCP ports.\nusage: nkill port")
	}

	// if os.Getpid() != 0 {
	// 	log.Println("WARNING: You are not running this script as superuser.")
	// }

	for _, port := range os.Args[1:] {
		p, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			log.Printf("%s is not a valid port number\n", port)
			continue
		}
		nkill.KillPort(p)
	}

}

```

##### Inspiration

http://voorloopnul.com/blog/a-python-netstat-in-less-than-100-lines-of-code/
