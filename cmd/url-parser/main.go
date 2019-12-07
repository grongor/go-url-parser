package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/grongor/go-url-parser"
)

func usage() {
	app := "url-parser"
	fmt.Printf("%s\n    Parses given URLs and prints the desired parts of them.\n\n", app)
	fmt.Println("Usage:")
	fmt.Printf("    %s <option> [url]...\n\n", app)
	fmt.Println("Options:")
	fmt.Printf("    -c, --component=COMPONENT\n\n")
	fmt.Printf("        Prints single URL component from each. URL.\n\n")
	fmt.Println("        Valid values: scheme, authority, auth, user, password,")
	fmt.Println("                      hostport, host, tld, port, path, query,")
	fmt.Printf("                      fragment, basePath, file, ext, relativeUrl\n\n")
	fmt.Printf("        Additionally, some components may be further formatted:\n")
	fmt.Println("            host:x, host:x:y            print only desired (sub)domains")
	fmt.Println("            hostport:x, hostport:x:y    same as host:*, but also includes port")
	fmt.Println("            path:x, path:x:y            print only desired parts of the path")
	fmt.Println("            query:NAME                  print only query parameter named NAME")
	fmt.Println("            fragment:NAME               print only fragment part named NAME")
	fmt.Println("")
	fmt.Println("                x    starting position; use -x to start from end")
	fmt.Println("                y    count; how many parts of the component you want to print")
	fmt.Println("")
	fmt.Printf("    -f, --format=FORMAT\n\n")
	fmt.Printf("        Prints URLs formatted according to FORMAT.\n\n")
	fmt.Println("        FORMAT is an arbitrary string containing component placeholders")
	fmt.Println("        enclosed in curly brackets. Each {COMPONENT} will be replaced")
	fmt.Println("        as if url-parser was called with option -c COMPONENT")
	fmt.Println("")
	fmt.Printf("Examples:\n\n")
	fmt.Printf("reference URL: https://www.example.com/lorem/ipsum/dolor/sit.html?metus=lectus#at=nostra&unde=omnis\n\n")
	fmt.Println("url-parser -c host:1 <url>               example.com")
	fmt.Println("url-parser -c path:-0:2 <url>            /dolor/sit.html")
	fmt.Println("url-parser -c path:-1 <url>              /lorem/ipsum/dolor")
	fmt.Println("url-parser -c path:0:2 <url>             /lorem/ipsum")
	fmt.Println("url-parser -c fragment:unde <url>        unde=omnis")
	fmt.Println("url-parser -c relativeUrl <url>          /lorem/ipsum/dolor/sit.html?metus=lectus#at=nostra&unde=omnis")
	fmt.Println("url-parser -f {scheme}://{host} <url>    https://www.example.com")
}

var (
	component string
	format    string
)

func main() {
	flag.StringVar(&component, "component", "", "")
	flag.StringVar(&component, "c", "", "")

	flag.StringVar(&format, "format", "", "")
	flag.StringVar(&format, "f", "", "")

	flag.Usage = usage
	flag.Parse()

	if component == "" && format == "" || component != "" && format != "" {
		fmt.Fprintln(os.Stderr, "You must specify either --component or --format option.")
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			process(scanner.Text())
		}
	} else {
		for _, rawUrl := range flag.Args() {
			process(rawUrl)
		}
	}
}

func process(rawUrl string) {
	var result string
	var err error
	if component == "" {
		result, err = urlparser.Format(rawUrl, format)
	} else {
		result, err = urlparser.Component(rawUrl, component)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	fmt.Println(result)
}
