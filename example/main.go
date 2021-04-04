package main

import (
	"errors"
	"fmt"
	"time"

	"bitbucket.org/armorblox/armorblox/pkg/dag"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize worker functions
	a := func() error {
		fmt.Println("Function A starts")
		time.Sleep(1 * time.Second)
		fmt.Println("Function A Finish")
		return nil
	}
	b := func() error {
		fmt.Println("Function B starts")
		time.Sleep(2 * time.Second)
		fmt.Println("Function B Finish")
		return nil
	}
	c := func() error {
		fmt.Println("Function C starts")
		time.Sleep(2 * time.Second)
		// err := errorGenerator()
		// if err != nil {
		// 	return err
		// }
		fmt.Println("Function C Finish")
		return nil
	}
	d := func() error {
		fmt.Println("Function D starts")
		time.Sleep(3 * time.Second)
		err := errorGenerator()
		if err != nil {
			return err
		}
		fmt.Println("Function D Finish")
		return nil
	}
	e := func() error {
		fmt.Println("Function E starts")
		time.Sleep(1 * time.Second)
		fmt.Println("Function E Finish")
		return nil
	}
	f := func() error {
		fmt.Println("Function F starts")
		time.Sleep(2 * time.Second)
		fmt.Println("Function F Finish")
		return nil
	}
	g := func() error {
		fmt.Println("Function G starts")
		time.Sleep(3 * time.Second)
		fmt.Println("Function G Finish")
		return nil
	}
	h := func() error {
		fmt.Println("Function H starts")
		time.Sleep(2 * time.Second)
		fmt.Println("Function H Finish")
		return nil
	}

	// Define the dag
	var conditions = []dag.Dependencies{
		{a, "a", nil, []string{"a", "b", "c", "d", "e", "f"}, true},
		{b, "b", []string{"a"}, []string{"b"}, true},
		{c, "c", []string{"a"}, []string{"c"}, true},
		{d, "d", []string{"a"}, []string{"d"}, false},
		{e, "e", []string{"c", "d"}, []string{"e"}, false},
		{f, "f", []string{"b"}, []string{"f"}, false},
		{g, "g", []string{"f", "e"}, []string{"g"}, false},
		{h, "h", []string{"f", "e"}, []string{"h"}, false},
	}

	// Initialize the execution status
	dag.Init(conditions)

	// Call the installer to start the execution
	dag.Installer(conditions)

	// Verify execution status
	if dag.VerifyStatus() {
		fmt.Println("Execution succeeded!!")
	} else {
		status, err := dag.FetchError()
		log.Infof("--- Execution Errors ---")
		log.Errorf("Status: %v", status)
		log.Infof("0: Initialized, 1: In progress, 2: Succeeded, 3: Retrying, 4: Failed")
		for message, error := range err {
			log.Errorf("Error: %v, %v", message, error)
		}
	}
}

// Test function that demostrate error generation
func errorGenerator() error {
	return errors.New("Error found")
}
