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
        version      = "latest" // Or pin a specific version such as v10.
        appID        = "..."    // See https://dev.freebox.fr/sdk/os/login/ and/or
        privateToken = "..."    // head to the next section of the documentation
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

## Generating credentials

At the time of this writing, generating credentials can only be done via the Freebox API. Please see [the documentation of this `terraform` provider](https://nikolalohinski.github.io/terraform-provider-freebox/provider.html#generating-credentials) which leverages `free-go` to provide a simple CLI to interact with the API and generate tokens.

## Supported and planned endpoints

- [x] [Authentication](https://dev.freebox.fr/sdk/os/login/) : `/login/*`
  - [x] Request authorization
  - [x] Track authorization progress (as part of the `Request authorization` process)
  - [x] Getting the challenge value
  - [x] Opening a session
  - [x] Closing the current session
- [ ] [Connection](https://dev.freebox.fr/sdk/os/connection/) : `/connection/*`
  - [ ] Get the current Connection status
  - [ ] Get the current Connection configuration
  - [ ] Update the Connection configuration
  - [ ] Get the current IPv6 Connection configuration
  - [ ] Update the IPv6 Connection configuration
  - [ ] Get the status of a DynDNS service
  - [ ] Set the config of a DynDNS service
- [x] [Discovery over HTTP](https://dev.freebox.fr/sdk/os/) : `/api_version`
- [ ] [Lan](https://dev.freebox.fr/sdk/os/lan/#lan) : `/lan/*`
  - [x] Getting the list of browsable LAN interfaces
  - [x] Getting the list of hosts on a given interface
  - [x] Getting a host information
  - [ ] Updating a host information
  - [ ] Wake on LAN
  - [ ] Get the current Lan configuration
  - [ ] Update the current Lan configuration
- [ ] [DHCP](https://dev.freebox.fr/sdk/os/dhcp/#dhcp) : `/dhcp/*`
  - [ ] Get the current DHCP configuration
  - [ ] Update the current DHCP configuration
  - [x] List the DHCP static leases
  - [x] Get a given DHCP static lease
  - [x] Update DHCP static lease
  - [x] Delete a DHCP static lease
  - [x] Add a DHCP static lease
  - [ ] Get the list of DHCP dynamic leases
- [x] [Port forwarding](https://dev.freebox.fr/sdk/os/nat/#port-forwarding): `/fw/redir/*`
  - [x] Getting the list of port forwarding
  - [x] Getting a specific port forwarding
  - [x] Updating a port forwarding
  - [x] Add a port forwarding
  - [x] Delete a port forwarding
- [ ] [Incoming port configuration](https://dev.freebox.fr/sdk/os/nat/#incoming-port-configuration) : `/fw/incoming/*`
  - [ ] Getting the list of incoming ports
  - [ ] Getting a specific incoming port
  - [ ] Updating an incoming port
- [ ] [Virtual machines](http://mafreebox.freebox.fr/#Fbx.os.app.help.app) (UNSTABLE) : `/vm/*`
  - [x] Get VM System Info
  - [x] Get Installable VM distributions
  - [x] Get the list of all VMs
  - [x] Get a VM
  - [x] Add a VM
  - [x] Delete a VM
  - [x] Update a VM
  - [x] Start a VM
  - [x] Send a powerbutton signal to a VM
  - [x] Stop a VM
  - [ ] Reset a VM
  - [ ] VM virtual console
  - [ ] VM virtual screen
  - [x] Get information on a virtual disk
  - [x] Create a virtual disk
  - [x] Resize a virtual disk
  - [x] Get a virtual disk task
  - [x] Delete a virtual disk task
- [x] [Websocket API](https://dev.freebox.fr/sdk/os/) : `/ws/*`
  - [x] WebSocket event API
  - [x] WebSocket file Upload API
- [ ] [Download API](https://dev.freebox.fr/sdk/os/download/) : `/downloads/*`
  - [x] Get a download task
  - [x] List download tasks
  - [x] Delete a download task
  - [x] Update a download task
  - [ ] Get a download log
  - [x] Add a new download task
- [x] [Upload API](https://dev.freebox.fr/sdk/os/upload/) : `/upload/*`
  - [x] Get an upload task
  - [x] List upload tasks
  - [x] Delete an upload task
  - [x] Cancel an upload task
  - [x] Cleanup upload tasks
  - [x] Start a new upload
- [ ] [Filesystem API](https://dev.freebox.fr/sdk/os/fs/) : `/fs/*`
  - [x] Get file information
  - [x] Download a file
  - [x] Remove files
  - [ ] List files
  - [x] Move files
  - [x] Copy files
  - [ ] Concatenate files
  - [ ] Create an archive
  - [x] Extract a file
  - [ ] Repair a file
  - [x] Hash a file
  - [x] Get a hash value
  - [x] Create a directory
  - [ ] Rename a file/folder
  - [x] List every task
  - [x] Get a task
  - [x] Delete a task
  - [x] Update a task
- [ ] [File Sharing Link](https://dev.freebox.fr/sdk/os/share/) : `/share_link/*`
  - [ ] List File Sharing links
  - [ ] Create a File Sharing link
  - [ ] Retrieve a File Sharing link
  - [ ] Delete a File Sharing link
- [ ] [Wi-Fi](https://dev.freebox.fr/sdk/os/wifi/) : `/wifi/*`
  - [ ] Get the current Wi-Fi global configuration
  - [ ] Update the Wi-Fi global configuration
  - [ ] List the Wi-Fi Access Points
  - [ ] Get a specific Access Point
  - [ ] Update an Access Point configuration
  - [ ] Get the Wi-Fi allowed combinations for the given Access Point
  - [ ] List the Wi-Fi Stations (connected devices)
  - [ ] List the Basic Service Sets
  - [ ] Get a specific Basic Service Set
  - [ ] Update a Basic Service Set
  - [ ] List the neighbors for the given Access Point
  - [ ] List the Wi-Fi channels usages for the given Access Point
  - [ ] Refresh the radar informations
  - [ ] Get the Wi-Fi Planning configuration
  - [ ] Update the Wi-Fi Planning configuration
  - [ ] List the MAC Filter entries
  - [ ] Get a specific MAC Filter entry
  - [ ] Update a MAC Filter entry
  - [ ] Delete a MAC Filter entry
  - [ ] Create a MAC Filter entry
  - [ ] Reset the Wi-Fi configuration
- [ ] [System](https://dev.freebox.fr/sdk/os/system/) : `/system/*`
  - [ ] Reboot the Freebox
- [ ] [AirMedia](https://dev.freebox.fr/sdk/os/airmedia/) : `/airmedia/*`
  - [ ] Get the AirMedia configuration
  - [ ] Update the AirMedia configuration
  - [ ] Get the list of AirMedia receivers
  - [ ] Sending a new request to an AirMedia receiver
- [ ] [Call](https://dev.freebox.fr/sdk/os/call/) : `/call/*`
  - [ ] List the calls
  - [ ] Delete all calls
  - [ ] Mark all calls as read
  - [ ] Get a call
  - [ ] Delete a call
  - [ ] Update a call entry
- [ ] [Contact](https://dev.freebox.fr/sdk/os/contacts/) : `/contact/*`
  - [ ] List the contacts
  - [ ] Get a contact
  - [ ] Create a contact
  - [ ] Delete a contact
  - [ ] Update a contact
  - [ ] List the contact numbers
  - [ ] Get a contact number
  - [ ] Create a contact number
  - [ ] Delete a contact number
  - [ ] Update a contact number
- [ ] [FreePlugs](https://dev.freebox.fr/sdk/os/freeplug/) : `/freeplug/*`
  - [ ] List the Freeplugs networks and its members
  - [ ] Get a specific Freeplug
  - [ ] Reset a Freeplug
- [ ] [Parental](https://dev.freebox.fr/sdk/os/parental/) : `/parental/*`
  - [ ] Get parental filter configuration
  - [ ] Update parental filter configuration
  - [ ] List the parental filter rules
  - [ ] Get a parental filter rule
  - [ ] Delete a parental filter rule
  - [ ] Update a parental filter rule
  - [ ] Create a parental filter rule
  - [ ] Get the planning for a parental filter rule
  - [ ] Update the planning for a parental filter rule
- [ ] [LCD](https://dev.freebox.fr/sdk/os/lcd/) : `/lcd/*`
  - [ ] Get the current LCD configuration
  - [ ] Update the LCD configuration
- [ ] [Switch](https://dev.freebox.fr/sdk/os/switch/) : `/switch/*`
  - [ ] Get the switch status and the list of ports
  - [ ] Get a specific port configuration
  - [ ] Update a port configuration
- [ ] [Universal Plug and Play Audio Video](https://dev.freebox.fr/sdk/os/upnpav/) : `/upnpav/*`
  - [ ] Get the UPnP AV configuration
  - [ ] Update UPnP AV configuration

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
