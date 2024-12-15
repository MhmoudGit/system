package main

func main() {
	server, err := NewServer()
	if err != nil {
		server.Logger.Error("Failed to initialize server: %s\n", "error", err)
	}

	server.Start()
}
