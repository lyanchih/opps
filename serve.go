package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.adlinktech.com/lyan.hung/opps/engine"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	address string
	port    uint32
)

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "OPPS serve",
		RunE:  runServe,
	}

	f := cmd.Flags()
	f.StringVarP(&address, "address", "a", "0.0.0.0", "Serve Listen Address")
	f.Uint32VarP(&port, "port", "p", 7070, "Serve Listen Port")
	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	runOpps(cmd, args)

	http.HandleFunc("/hook/", handleHook)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}

func handleHook(w http.ResponseWriter, req *http.Request) {
	log.Printf("Hook %s %s\n", req.Method, req.URL)
	if req.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	engineName := strings.TrimPrefix(req.URL.Path, "/hook/")
	e, err := engine.TranslateEngine(engineName)
	if err == engine.ErrEngineNotSupport {
		w.WriteHeader(404)
		return
	}

	b := req.Body
	defer b.Close()
	data, err := ioutil.ReadAll(b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	err = e.HandleHook(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
	}
}
