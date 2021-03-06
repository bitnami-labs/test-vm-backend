# Test VM Backend
Test back-end to be used by front-end projects.

Hassle free for candidates, just use the binaries provided and no need to setup any dev env or dependencies, so they can just focus on their front-end stuff.

## End user experience

The code in this repo compiles to a binaries for Linux, Mac or Windows. You get a ZIP file with all those binaries in a root folder (`./test-vm-backend-{version}/`), and the source code in another (`./test-vm-backend/`).

Just pick the binary of the right architecture and run it. Unless you want to build from sources.

## Building from sources

You need to have [Go installed locally](https://golang.org/doc/install) beforehand. Be sure it is Go 1.14 or 1.15 as those are the versions the source code has been tested with.

Once Go is installed locally just do:

~~~bash
$ cd test-vm-backend # if you were not already in the source folder.
$ go build
~~~

## Run in Docker

From within the `test-vm-backend` folder with the `Dockerfile` just do:

~~~bash
$ ./run_in_docker.sh
~~~

Parameters are accepted just as if that was a regular binary invocation.

## Sample usage

Launch the server on a terminal:

~~~bash
$ ./test-vm-backend # or...

$ ./run_in_docker.sh
Sending build context to Docker daemon  9.282MB
...
Successfully tagged test-vm-backend:latest
2020/11/10 09:32:07 Test VM Backend version Development
2020/11/10 09:32:07 Loading fake Cloud state from local file "vms.json"
API:
GET	    /vms                	-> VMs JSON            	# list All VMs
PUT	    /vms/{vm_id}/launch 	-> Check status code   	# launch VM by id
PUT	    /vms/{vm_id}/stop   	-> Check status code   	# stop VM by id
GET	    /vms/{vm_id}        	-> VM JSON             	# inspect a VM by id
DELETE	/vms/{vm_id}        	-> Check status code   	# delete a VM by id
2020/11/10 09:32:07 No UI folder given. Not serving any static files.
2020/11/10 09:32:07 Unlike a real production service this API accepts:
- Any Origin on CORS requests.
- Preflight OPTIONS request with any headers.
2020/11/10 09:32:07 Server listening at :8080
<- GET /vms
...
~~~

## Run on another port or address

Use the `--address` flag:

~~~bash
$ ./test-vm-backend --address "0.0.0.0:6060"
2020/11/10 09:36:13 Test VM Backend version Development
2020/11/10 09:36:13 Loading fake Cloud state from local file "vms.json"
API:
GET	    /vms                -> VMs JSON            	# list All VMs
PUT	    /vms/{vm_id}/launch -> Check status code   	# launch VM by id
PUT	    /vms/{vm_id}/stop   -> Check status code   	# stop VM by id
GET	    /vms/{vm_id}        -> VM JSON             	# inspect a VM by id
DELETE	/vms/{vm_id}        -> Check status code   	# delete a VM by id
2020/11/10 09:36:13 No UI folder given. Not serving any static files.
2020/11/10 09:36:13 Unlike a real production service this API accepts:
- Any Origin on CORS requests.
- Preflight OPTIONS request with any headers.
2020/11/10 09:36:13 Server listening at :6060

~~~

Same works for the docker invocation:
~~~bash
$ ./run_in_docker.sh --address=:6060
...
2020/09/17 21:23:04 Server listening at :6060
...
~~~

#### VM data & units

Each VM JSON looks like this:

~~~json
{
  "vcpus": 1,
  "clock": 1500,
  "ram": 4096,
  "storage": 128,
  "network": 1000,
  "state": "Running"
}
~~~

Units:
- `vcpus` is a integer number of virtual CPU cores.
- `clock` is measured in  `Mhz`.
- `ram` is `MiB`.
- `storage` is `GiB`.
- `network` is `Mbps`.
- `state` is one of `"Stopped"`, `"Starting"`, `"Running"`, `"Stopping"`.

## Testing

### Test drive with CURL

To test with curl, go to another terminal and write:

~~~bash
watch 'curl -s http://localhost:8080/vms |jq .'
~~~
Remove the `| jq .` bit tail if jq is not installed locally. It is optional but makes the JSON output more readable.

This shows how the server VMs change state as you interact with them from another terminal.

~~~json
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
~~~

Then issue requests on another terminal:

~~~bash
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
~~~

### Demotest

You can run `demotest.sh` for a quick happy path only test drive.

The demo test script expects the backend in default address `:8080` and the `uiFolder` to be set, possible to `./ui/`. More on the `uiFolder` argument in the next section.

~~~bash
$ ./demotest.sh 
Expects test-vm-backend running on default port: 8080
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
~~~

Demotest expects the backend just launched to work, from initial state.

**NOTE: The start and stop delays have random durations around 10 and 5 seconds each. DO NOT rely on as durations and poll for completion like demotest.sh does.**

## CORS bypass

This backend is intended as a quick tool to help develop a frontend quickly.

Developers may chose to have the code running in the browser call the API directly. If the code was loaded form a different process, the origin domain will differ and [the browser will send CORS headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS).

Production APIs usually do not need to handle CORS requests because are rarely accessed directly from a browser:
- Most of the time a load-balancer or a reverse proxy will reroute traffic from to a web ui or API as needed, all under the same domain.
- Or the WebApp backend will resend the API request itself on behalf o the browser code.

To make things simpler for the frontend deployment, this server bypasses CORS issues by allowing any requesting `Origin` and supporting pre-flight `OPTIONS` requests that allow any requested headers.

Needless to say this is not a safe setup for production, but **this is not a production-ready server**.

## Serving a simple UI from a local folder

If the Frontend consists mostly on code running on the browser directly, it might be handy to serve the frontend files as static files from this same backend.

That way the browser will not see the frontend being loaded from a domain, and then its code using an API from another domain. Otherwise the setup would cause Cross Origin Resource Sharing errors.

~~~
[Code in Browser] <=> [test-vm-backend API + ui frontend file serving]
~~~

To that end the `test-vm-backend` accepts `string` parameter `--uiFolder` so that the given folder files will be served to the browser from the same backend program:

~~~bash
$ ./test-vm-backend --uiFolder=./ui/
2020/11/01 10:51:27 Test VM Backend version Development
...
2020/11/01 10:51:27 Serving static files for the UI at "./ui/"
~~~

A Sample static file served from the this program:

~~~bash
$ curl http://localhost:8080/ui/vms.html
<!doctype html>
<html>
<head>
<!-- 
  Sample static file draft for starting a Single Page UI App

  => Feel free to remove, replace or extend as needed
-->
<title>Sample UI</title>
<meta charset="UTF-8">
<meta name="description" content="HTML5 sample web ui">
</head>
<body>
<h1>VMs manager Web UI</h1>
<div id="vms"></div>
<script>
  document.getElementById("vms").innerHTML = "Load VMs here...";
</script> will be
</body>
</html>
~~~

**This feature is totally optional**. In fact, by default, when uiFolder is not set or set to empty no static files are served.

## Customizing initial state

Notice the output lines in the example above:

~~~
Loading fake Cloud state from local file "vms.json"
Missing "vms.json", generating one...
Tip: You can tweak "vms.json"  adding VMs or changing states for next run.
...
~~~

If you run the server at least once it will create a default `vms.json` file you can tweak to your liking. The initial contents of that file should look like the first call to the `/vms` endpoint:

~~~json
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
~~~

From that you can add/remove or tweak VM entries and re-run to start from a new initial state.

If you are running from the container, note that by default the `vms.json` file used is the one from within the container, not your host filesystem.
