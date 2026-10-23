<h1 align="center">
  defango
  <br>
</h1>
<h4 align="center">URL / IP / Email defanging with Golang. Make IoC harmless.</h4>
<h6 align="center"> Coded with 💙 by edoardottt </h6>
<p align="center">
  <a href="https://github.com/edoardottt/defango/actions">
      <img src="https://github.com/edoardottt/defango/actions/workflows/go.yml/badge.svg" alt="go action">
  </a>
<br>
  <!--Tweet button-->
  <a href="https://twitter.com/intent/tweet?text=defango%20-%20URL%20%2F%20IP%20%2F%20Email%20defanging%20with%20Golang.%20Make%20IoC%20harmless.%20https%3A%2F%2Fgithub.com%2Fedoardottt%2Fdefango%20%23golang%20%23github%20%23linux%20%23infosec%20%23malware" target="_blank">Share on Twitter!
  </a>
</p>

<p align="center">
  <a href="#install-">Install</a> •
  <a href="#get-started-">Get Started</a> •
  <a href="#changelog-">Changelog</a> •
  <a href="#contributing-">Contributing</a> •
  <a href="#license-">License</a>
</p>

Install 📡
----------

```console
go get github.com/edoardottt/defango
```

CLI:

```console
go install github.com/edoardottt/defango/cmd/defango@latest
```

Get Started 🎉
----------

```go
package main

import (
    "fmt"
    "log"
    "github.com/edoardottt/defango"
)

func main() {
    fmt.Println(defango.IP("8.8.8.8:53")) // 8[.]8[.]8[.]8[:]53

    u, err := defango.URL("http://malicious.example.link.com:666/m4lw4r3.exe")
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(u) // hxxp://malicious[.]example[.]link[.]com[:]666/m4lw4r3[.]exe

    fmt.Println(defango.Email("mailto:edoardott@gmail.com")) // mailto[:]edoardott@gmail[.]com

    // Defang only the IoCs in a text.
    fmt.Println(defango.Text("C2 at evil.com, payload: https://1.2.3.4/x.exe"))
    // C2 at evil[.]com, payload: hxxps://1[.]2[.]3[.]4/x[.]exe

    fmt.Println(defango.Refang("hxxps://evil[.]com"))     // https://evil.com
    fmt.Println(defango.IsDefanged("hxxps://evil[.]com")) // true
}
```

`http`, `https`, `ftp` and `file` become `hxxp`, `hxxps`, `fxp` and `fxle`; any other scheme gets `[://]`.
Pass `defango.HostOnly()` to `URL` or `Text` to leave paths and queries as they are.

The CLI defangs its standard input, or refangs it with `-r`:

```console
$ echo "C2 at evil.com, payload: https://1.2.3.4/x.exe" | defango
C2 at evil[.]com, payload: hxxps://1[.]2[.]3[.]4/x[.]exe
$ echo "hxxps://evil[.]com" | defango -r
https://evil.com
```

Read the full [`package documentation here`](https://pkg.go.dev/github.com/edoardottt/defango).

Changelog 📌
-------

Detailed changes for each release are documented in the [release notes](https://github.com/edoardottt/defango/releases).

Contributing 🛠
-------

Just open an [issue](https://github.com/edoardottt/defango/issues) / [pull request](https://github.com/edoardottt/defango/pulls).

Before opening a pull request, download [golangci-lint](https://golangci-lint.run/usage/install/) and run

```console
golangci-lint run
```

If there aren't errors, go ahead :)

License 📝
-------

This repository is under [MIT License](https://github.com/edoardottt/defango/blob/main/LICENSE).  
[edoardottt.com](https://edoardottt.com) to contact me.
