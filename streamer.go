package main

func RunStreamer(history []string, live <-chan []byte, output chan<- []byte) {
	for _, m := range history {
		output <- []byte(m)
	}
	for m := range live {
		output <- m
	}
	close(output)
}
