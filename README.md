<div align="center">
<img src="./free-go.svg" width="250"/>

<i><a href="https://en.wikipedia.org/wiki/Freebox" target="_blank">Freebox</a> client in Golang</i>
</div>


## Usage

First get the library locally:

```shell
go get github.com/nikolalohinski/free-go
```

And use it as follows:

```go
package main

import (
    "github.com/nikolalohinski/free-go/client"
    "github.com/nikolalohinski/free-go/types"
)

func main() {
    var (
        endpoint     = "mafreebox.freebox.fr"
        version      = "v10" 
        appID        = "..." // See https://dev.freebox.fr/sdk/os/login/ to understand 
        privateToken = "..." // how to define an app and generate a private token
    )

    freebox, err := client.New(endpoint, version).
        WithAppID(appID).
        WithPrivateToken(privateToken)
    if err != nil {
        panic(err)
    }

    permissions, err := freebox.Login()
    if err != nil {
        panic(err)
    }

    fmt.Println(permissions)
}
```

## Development

### Requirements

* Install `go` (`>= v1.20`) following the [official instructions](https://go.dev/doc/install) ;
* Install `mage` using the [online documentation](https://magefile.org/²) ;
* Run the following to fetch all the required tools:
  ```shell
  mage install
  ```
* Verify the previous steps by running:
  ```shell
  mage
  ```

### Tests

To run the unit tests:

```shell
mage go:test
```

To generate and open a coverage report:

```shell
mage go:cover
```

To run the integration tests, you will first need the following environment variables defined:
* `FREEBOX_ENDPOINT`: IP Address or DNS name to reach out to your Freebox. Usually `mafreebox.freebox.fr` works ;
* `FREEBOX_VERSION`: API version of the freebox you want to run against. For example `v10` ;
* `FREEBOX_APP_ID`: The ID of the application you created to authenticate to the Freebox (see [the login documentation](https://dev.freebox.fr/sdk/os/login/)) ;
* `FREEBOX_TOKEN`: The private token to authenticate to the Freebox (see [the login documentation](https://dev.freebox.fr/sdk/os/login/)) ;

Then, you should be able to run:

```shell
mage go:integration
```