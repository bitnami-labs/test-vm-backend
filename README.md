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
Loading fake Cloud state from local file "vms.json"
Missing "vms.json", generating one...
Tip: You can tweak "vms.json"  adding VMs or changing states for next run.
Server listening at :8080
API:
GET     /vms[/]?                -> VMs JSON             # list All VMs
PUT     /vms/launch/\d+         -> Check status code    # launch VM by id
PUT     /vms/stop/\d+           -> Check status code    # stop a VM by id
GET     /vms/\d+                -> VM JSON              # inspect a VM by id
DELETE  /vms/\d+                -> VM JSON              # delete a VM by id

<- GET /vms
...
```
## Test drive with CURL

To test with curl, go to another terminal and write:

```bash
watch 'curl -s http://localhost:8080/vms |jq .'
```
Remove the `| jq .` bit tail if jq is not installed locally. It is optional but makes the JSON output more readable.

This shows how the server VMs change state as you interact with them from another terminal.

```
Every 2s: curl -s http://localhost:8080/vms |jq . 
{
  "0": {
    "vcpus": 1,
    "clock": 1500,
    "ram": 4096,
    "storage": 128,
    "network": 1000,
    "state": "Running"
  },
  "1": {
    "vcpus": 4,
    "clock": 3600,
    "ram": 32768,
    "storage": 512,
    "network": 10000,
    "state": "Stopped"
  },
  "2": {
    "vcpus": 2,
    "clock": 2200,
    "ram": 8192,
    "storage": 256,
    "network": 1000,
    "state": "Stopped"
  }
}
```

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

$ curl -s -X PUT http://localhost:8080/vms/launch/0 
$ 

$ curl -s -X PUT http://localhost:8080/vms/stop/0 
$ curl -s -X POST http://localhost:8080/vms/stop/0 
POST /vms/stop/0 not allowed
$ curl -s -X PUT http://localhost:8080/vms/stop/0 
Illegal transition from "Stopped" to "Stopping"
$ curl -s -X DELETE http://localhost:8080/vms/0 
$ curl -s http://localhost:8080/vms/0 |jq .
{}

```
# Customizing initial state

Notice the output lines in the example above:

```
Loading fake Cloud state from local file "vms.json"
Missing "vms.json", generating one...
Tip: You can tweak "vms.json"  adding VMs or changing states for next run.
...
```

If you run the server at leats once it will create a default `vms.json` file you can tweak to you liking. The initial contents of that file should look like the first call to the `/vms` endpoint:

```json
$ cat vms.json |jq .
[
  {
    "vcpus": 1,
    "clock": 1500,
    "ram": 4096,
    "storage": 128,
    "network": 1000,
    "state": "Stopped"
  },
  {
    "vcpus": 4,
    "clock": 3600,
    "ram": 32768,
    "storage": 512,
    "network": 10000,
    "state": "Stopped"
  },
  {
    "vcpus": 2,
    "clock": 2200,
    "ram": 8192,
    "storage": 256,
    "network": 1000,
    "state": "Stopped"
  }
]
```

From that you can add/remove or tweak VM entries and re-run to start from a new initial state.
