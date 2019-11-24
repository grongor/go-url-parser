url-parser
==========

```
url-parser
    Parses given URLs and prints the desired parts of them.

Usage:
    url-parser <option> [url]...

Options:
    -c, --component=COMPONENT

        Prints single URL component from each. URL.

        Valid values: scheme, authority, auth, user, password,
                      hostport, host, tld, port, path, query,
                      fragment, basePath, file, ext, relativeUrl

        Additionally, some components may be further formatted:
            host:x, host:x:y            print only desired (sub)domains
            hostport:x, hostport:x:y    same as host:*, but also includes port
            path:x, path:x:y            print only desired parts of the path
            query:NAME                  print only query parameter named NAME
            fragment:NAME               print only fragment part named NAME

                x    starting position; use -x to start from end
                y    count; how many parts of the component you want to print

    -f, --format=FORMAT

        Prints URLs formatted according to FORMAT.

        FORMAT is an arbitrary string containing component placeholders
        enclosed in curly brackets. Each {COMPONENT} will be replaced
        as if url-parser was called with option -c COMPONENT

Examples:

reference URL: https://www.example.com/lorem/ipsum/dolor/sit.html?metus=lectus#at=nostra&unde=omnis

url-parser -c host:1 <url>               example.com
url-parser -c path:-0:2 <url>            /dolor/sit.html
url-parser -c path:-1 <url>              /lorem/ipsum/dolor
url-parser -c path:0:2 <url>             /lorem/ipsum
url-parser -c fragment:unde <url>        unde=omnis
url-parser -c relativeUrl <url>          /lorem/ipsum/dolor/sit.html?metus=lectus#at=nostra&unde=omnis
url-parser -f {scheme}://{host} <url>    https://www.example.com
```
