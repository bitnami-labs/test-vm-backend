# test-vmbackend
Test back-end to be used by front-end projects.

Hassle free for candidates, just use the binaries provided and no need to setup any dev env or dependencies, so they can just focus on their front-end stuff.

# End user experience

The code in this repo compiles to a binary for Linux, Mac or Windows (or all of them a the same time).

The user just downloads the binary of the right architecture and runs it.

# Sample usage

Launch the server on a terminal:

```bash
$ go build && ./test-vmbackend 
Server listening at :8080
API
GET /vms             -> VMs JSON          # list all VMs
PUT /vms/launch/{id} -> Check status code # launch VM by id
PUT /vms/stop/{id}   -> Check status code # a VM by id
GET /vms/{id}        -> VM JSON           # inspect a VM by id
DELETE /vms/{id}     -> Check status code # delete a VM by id
...
```
## Test drive with CURL

To test with curl, go to another terminal and write:

```bash
watch 'curl -s http://localhost:8080/vms |jq .'
```
Remove the `| jq .` bit tail if jq is not installed locally. It is optional but makes the JSON output more readable.

This shows how the server VMs change state as you interact with them from another terminal.

Then issue requests on another terminal:

```bash
$ curl -s http://localhost:8080/vms/0 |jq .
{
  "vcpus": 1,
  "clock": 1500,
  "ram": 4096,
  "storage": 128,
  "network": 1000,
  "state": "Stopped"
}

```
