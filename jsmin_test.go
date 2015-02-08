package jsmin_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/omeid/jsmin"
)

var before = strings.NewReader(`
// (c) 2001 Douglas Crockford
// 2001 June 3
// is
// The -is- object is used to identify the browser. Every browser edition
// identifies itself, but there is no standard way of doing it, and some of
// the identification is deceptive. This is because the authors of web
// browsers are liars. For example, Microsoft's IE browsers claim to be
// Mozilla 4. Netscape 6 claims to be version 5.
var is = {
ie: navigator.appName == 'Microsoft Internet Explorer',
java: navigator.javaEnabled(),
ns: navigator.appName == 'Netscape',
ua: navigator.userAgent.toLowerCase(),
version: parseFloat(navigator.appVersion.substr(21)) ||
parseFloat(navigator.appVersion),
win: navigator.platform == 'Win32'
}
is.mac = is.ua.indexOf('mac') &gt;= 0;
if (is.ua.indexOf('opera') &gt;= 0) {
is.ie = is.ns = false;
is.opera = true;
}
if (is.ua.indexOf('gecko') &gt;= 0) {
is.ie = is.ns = false;
is.gecko = true;
}`)

const after = `
var is={ie:navigator.appName=='Microsoft Internet Explorer',java:navigator.javaEnabled(),ns:navigator.appName=='Netscape',ua:navigator.userAgent.toLowerCase(),version:parseFloat(navigator.appVersion.substr(21))||parseFloat(navigator.appVersion),win:navigator.platform=='Win32'}
is.mac=is.ua.indexOf('mac')&gt;=0;if(is.ua.indexOf('opera')&gt;=0){is.ie=is.ns=false;is.opera=true;}
if(is.ua.indexOf('gecko')&gt;=0){is.ie=is.ns=false;is.gecko=true;}`

func TestMinify(t *testing.T) {

	reader, err := jsmin.Minify(before)
	if err != nil {
		t.Fatal(err)
	}

	out, _ := ioutil.ReadAll(reader)
	if string(out) != after {
		t.Fatalf("Whops. We failed.\nOUTPUT:\n%s\n", out)
	}
}
