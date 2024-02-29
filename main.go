package main

func main() {
	server := InitWebServer()
	if err := server.Run(":8081"); err != nil {
		return
	}
}
