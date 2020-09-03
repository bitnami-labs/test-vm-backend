package main

import (
	"fmt"
	"net/http"
)

func main() {
	vms := VMServer{Cloud{vms: loadVMs()}, ":8080"}

	fmt.Println("API:")
	fmt.Printf("GET %v/vms             -> list all VMs\n", vms.address)
	fmt.Printf("PUT %v/vms/launch/{id} -> launch a VM by id\n", vms.address)
	fmt.Printf("PUT %v/vms/stop/{id}   -> stop a VM by id\n", vms.address)
	fmt.Printf("GET %v/vms/{id}        -> inspect a VM by id\n", vms.address)
	fmt.Printf("DELETE %v/vms/{id}     -> delete a VM by id\n", vms.address)
	http.HandleFunc("/vms", vms.list)
	http.HandleFunc("/vms/launch/", vms.launch)
	http.HandleFunc("/vms/stop/", vms.stop)
	http.HandleFunc("/vms/", vms.serveVM)

	http.ListenAndServe(vms.address, nil)
}
