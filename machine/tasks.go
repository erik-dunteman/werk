package machine

import (
	"fmt"
	"time"
)

/*
* Define your task handlers here.
* These functions are the functions ran by the workers.
 */

// Sleep
func sleepHandler(t int) (bool, error) {
	time.Sleep(2 * time.Second)
	return true, nil
}

// Say Hello
func helloHandler(name string) (string, error) {
	out := fmt.Sprintf("Hello %s", name)
	return out, nil
}

// // Do xyz stuff - uncomment this block to add a new function
// func xyzHandler(t int) (bool, error) {
// 	// add logic here
// 	return true, nil
// }

// Tasks ...
var Tasks = map[string]interface{}{
	"sleep": sleepHandler,
	"hello": helloHandler,
	// "xyz": xyzHandler, // uncomment this line to add a new function
}
