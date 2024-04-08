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
    "fmt"
    "context"

    "github.com/nikolalohinski/free-go/client"
)

func main() {
    var (
        endpoint     = "mafreebox.freebox.fr"
        version      = "v10" 
        appID        = "..." // See https://dev.freebox.fr/sdk/os/login/ to learn
        privateToken = "..." // how to define an app and generate a private token
    )

    ctx := context.Background()

    freebox, err := client.New(endpoint, version).
        WithAppID(appID).
        WithPrivateToken(privateToken)
    if err != nil {
        panic(err)
    }

    permissions, err := freebox.Login(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println(permissions)

    vms, err := freebox.ListVirtualMachines(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println(vms)
}
```

For details on how to use this client, please refer to the `Client` interface in [`client/client.go`](./client/client.go).

## Supported and planned endpoints

- [x] [Authentication](https://dev.freebox.fr/sdk/os/login/) : `/login/*`
  - [x] Request authorization
  - [x] Track authorization progress (as part of the `Request authorization` process)
  - [x] Getting the challenge value 
  - [x] Opening a session
  - [x] Closing the current session
- [x] [Discovery over HTTP](https://dev.freebox.fr/sdk/os/) : `/api_version`
- [ ] [Lan](https://dev.freebox.fr/sdk/os/lan/#lan) : `/lan/*`
  - [x] Getting the list of browsable LAN interfaces
  - [x] Getting the list of hosts on a given interface
  - [x] Getting a host information
  - [ ] Updating a host information
  - [ ] Wake on LAN
  - [ ] Get the current Lan configuration
  - [ ] Update the current Lan configuration
- [x] [Port forwarding](https://dev.freebox.fr/sdk/os/nat/#port-forwarding): `/fw/redir/*`
  - [x] Getting the list of port forwarding
  - [x] Getting a specific port forwarding
  - [x] Updating a port forwarding
  - [x] Add a port forwarding
  - [x] Delete a port forwarding
- [ ] [Virtual machines](http://mafreebox.freebox.fr/#Fbx.os.app.help.app) (UNSTABLE) : `/vm/*`
  - [x] Get VM System Info
  - [x] Get Installable VM distributions
  - [x] Get the list of all VMs
  - [x] Get a VM
  - [x] Add a VM
  - [x] Delete a VM
  - [x] Update a VM
  - [ ] Start a VM
  - [ ] Send a powerbutton signal to a VM
  - [ ] Stop a VM
  - [ ] Reset a VM
  - [ ] Watch for VM status changes
  - [ ] VM virtual console
  - [ ] VM virtual screen
  - [ ] Get information on a virtual disk
  - [ ] Create a virtual disk
  - [ ] Resize a virtual disk
  - [ ] Get a virtual disk task
  - [ ] Delete a virtual disk task
- [ ] [Websocket API](https://dev.freebox.fr/sdk/os/) : `/ws/*`
  - [x] WebSocket event API
  - [ ] WebSocket file Upload API
- [ ] [Download API](https://dev.freebox.fr/sdk/os/download/) : `/downloads/*`
  - [ ] Get a download task
  - [ ] List download tasks
  - [ ] Delete a download task
  - [ ] Update a download task
  - [ ] Get a download log
  - [ ] Add a new download task
- [ ] [Filesystem API](https://dev.freebox.fr/sdk/os/fs/) : `/fs/*`
  - [x] Get file information
  - [ ] Download a file
  - [x] Remove files
  - [ ] List files
  - [ ] Move files
  - [ ] Copy files
  - [ ] Concatenate files
  - [ ] Create an archive
  - [ ] Extract a file
  - [ ] Repair a file
  - [ ] Hash a file
  - [ ] Get a hash value
  - [ ] Create a directory
  - [ ] Rename a file/folder

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

## About

This project aims to provide the base Go components to create a `terraform` provider for the Freebox Delta to be able to leverage its VMs scheduling capabilities via infrastructure as code. It also aims to reach feature parity with [`juju2013/go-freebox`](https://github.com/juju2013/go-freebox), [`moul/go-freebox`](https://github.com/moul/go-freebox) and eventually [`hacf-fr/freebox-api`](https://github.com/hacf-fr/freebox-api) but with actual unit tests, maximal code coverage and integration testing.

Contributions are welcomed !