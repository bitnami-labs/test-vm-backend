# test-vmbackend
Test back-end to be used by front-end projects.

Hassle free for candidates, just use the binaries provided and no need to setup any dev env or dependencies, so they can just focus on their front-end stuff.

# End user experience

The code in this repo compiles to a binaries for Linux, Mac or Windows. You get a ZIP file with all those binaries in a root folder (`./test-vmbackend-{version}/`), and the source code in another (`./test-vmbackend/`).

Just pick the binary of the right architecture and run it. Unless you want to build from sources.

# Building from sources

You need to have [Go installed locally](https://golang.org/doc/install) beforehand. Be sure it is Go 1.14 or 1.15 as those are the versions the source code has been tested with.

Once Go is installed locally just do:

```
$ cd test-vmbackend # if you were not already in the source folder.
$ go build
```

# Run in Docker

From within the `test-vmbackend` folder with the `Dockerfile` just do:

```
./run_in_docker.sh
```

Parameters are accepted just as if that was a regular binary invocation.

# Sample usage

Launch the server on a terminal:

```bash
$ ./test-vmbackend # or ./run_in_docker.sh
2020/09/15 19:53:56 Test-VMBackend version Development
2020/09/15 19:53:56 Loading fake Cloud state from local file "vms.json"
2020/09/15 19:53:56 Missing "vms.json", generating one...
2020/09/15 19:53:56 Tip: You can tweak "vms.json" adding VMs or changing states for next run.
2020/09/15 19:53:56 Server listening at :8080
API:
GET     /vms                    -> VMs JSON             # list All VMs
PUT     /vms/{vm_id}/launch     -> Check status code    # launch VM by id
PUT     /vms/{vm_id}/stop       -> Check status code    # stop VM by id
GET     /vms/{vm_id}            -> VM JSON              # inspect a VM by id
DELETE  /vms/{vm_id}            -> Check status code    # delete a VM by id

<- GET /vms
...
```

### Run on another port or address

Use the `-address` flag:

```
$ ./test-vmbackend -address "0.0.0.0:6060"
2020/09/15 19:43:54 Test-VMBackend version Development
2020/09/15 19:43:54 Loading fake Cloud state from local file "vms.json"
2020/09/15 19:43:54 Server listening at 0.0.0.0:6060
API:
GET	/vms                	-> VMs JSON            	# list All VMs
PUT	/vms/{vm_id}/launch 	-> Check status code   	# launch VM by id
PUT	/vms/{vm_id}/stop   	-> Check status code   	# stop VM by id
GET	/vms/{vm_id}        	-> VM JSON             	# inspect a VM by id
DELETE	/vms/{vm_id}        	-> Check status code   	# delete a VM by id
```

Same works for the docker invocation:
```
$ ./run_in_docker.sh --address=:6060
...
2020/09/17 21:23:04 Server listening at :6060
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

$ curl -s -X PUT http://localhost:8080/vms/0/launch
$ 

$ curl -s -X PUT http://localhost:8080/vms/0/stop
$ curl -s -X PUT http://localhost:8080/vms/0/stop
illegal transition from "Stopped" to "Stopping"
$ 

$ curl -s -X DELETE http://localhost:8080/vms/0
$ curl -s -X PUT http://localhost:8080/vms/0/stop
not found VM with id 0
$ curl -s http://localhost:8080/vms/0
{}
```

### Demotest

You can run `demotest.sh` for a quick happy path only test drive:

```
$ ./demotest.sh 
Expects test-vmbackend running on default port: 8080
GET http://localhost:8080/vms
{"0":{"vcpus":1,"clock":1500,"ram":4096,"storage":128,"network":1000,"state":"Stopped"},"1":{"vcpus":4,"clock":3600,"ram":32768,"storage":512,"network":10000,"state":"Stopped"},"2":{"vcpus":2,"clock":2200,"ram":8192,"storage":256,"network":1000,"state":"Stopped"}}
GET http://localhost:8080/vms/0
{"vcpus":1,"clock":1500,"ram":4096,"storage":128,"network":1000,"state":"Stopped"}
PUT http://localhost:8080/vms/0/launch

GET http://localhost:8080/vms/0
{"vcpus":1,"clock":1500,"ram":4096,"storage":128,"network":1000,"state":"Starting"}
Wait for started...
GET http://localhost:8080/vms/0
{"vcpus":1,"clock":1500,"ram":4096,"storage":128,"network":1000,"state":"Running"}
PUT http://localhost:8080/vms/0/stop

GET http://localhost:8080/vms/0
{"vcpus":1,"clock":1500,"ram":4096,"storage":128,"network":1000,"state":"Stopping"}
Wait for stopped
GET http://localhost:8080/vms/0
{"vcpus":1,"clock":1500,"ram":4096,"storage":128,"network":1000,"state":"Stopped"}
DELETE http://localhost:8080/vms/0

GET http://localhost:8080/vms
{"1":{"vcpus":4,"clock":3600,"ram":32768,"storage":512,"network":10000,"state":"Stopped"},"2":{"vcpus":2,"clock":2200,"ram":8192,"storage":256,"network":1000,"state":"Stopped"}}
Demotest: OK/PASS
```

Demotest expects the backend just launched to work, from initial state.

# Customizing initial state

Notice the output lines in the example above:

```
Loading fake Cloud state from local file "vms.json"
Missing "vms.json", generating one...
Tip: You can tweak "vms.json"  adding VMs or changing states for next run.
...
```

If you run the server at least once it will create a default `vms.json` file you can tweak to you liking. The initial contents of that file should look like the first call to the `/vms` endpoint:

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
