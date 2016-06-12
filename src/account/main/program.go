package main

import (
	"account/db"
	"fmt"
	"lib/network"
	"os"
)

func main() {
	// By law, this header and the copyright must be included in derivations of
	// this work. You may change the name of the project, but not the license or
	// copyright. All derivations must include this copyright and license.
	fmt.Println(`GoConquer, Account Server`)
	fmt.Println("Copyright(C) 2016 Gareth Warry, Matt Moening")
	fmt.Println("Version 1.0, May 2016\n")
	fmt.Println("This work is licensed under the Creative Commons Attribution-");
	fmt.Println("NonCommercial-ShareAlike 4.0 International (CC-BY-NC) License.");
	fmt.Println("A copy of this license is available to you in the distribution");
	fmt.Println("of this software.\n");
	
	// Read in the user's configuration file for the server.
	fmt.Println("Initializing server...")
	err := db.Configuration.Decode("./configuration.json")
	if err != nil { fmt.Println(err.Error()); os.Exit(-1) }
	
	// Load flat-file database.
	if !db.LoadGameServers() { fmt.Printf("failed\n"); os.Exit(-1) }
	
	// Create the server instance and start listening.
	ch := make(chan bool)
	server := network.Server{} 
	server.OnConnect = OnConnect
	server.OnReceive = OnReceive
	go server.Listen(db.Configuration.Host, ch) 
	fmt.Println("Listening for new connections\n")
	
	// Terminate the program only when done listening for connections.
	serverstate := <-ch
	if !serverstate {
		fmt.Println("server terminated unexpectedly")
		os.Exit(-1)
	}
}