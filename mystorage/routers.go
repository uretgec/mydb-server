package main

import "github.com/tidwall/redcon"

func setupRoutes() *redcon.ServeMux {

	mux := redcon.NewServeMux()
	mux.HandleFunc("detach", handlerDetach)
	mux.HandleFunc("quit", handlerQuit)
	mux.HandleFunc("ping", handlerPing)

	//mux.HandleFunc("pset", handlerProcessSet)
	mux.HandleFunc("set", handlerSet)
	mux.HandleFunc("get", handlerGet)
	mux.HandleFunc("mget", handlerMGet)
	mux.HandleFunc("list", handlerList)
	mux.HandleFunc("prevlist", handlerPrevList)
	mux.HandleFunc("del", handlerDel)
	mux.HandleFunc("exists", handlerExists)
	mux.HandleFunc("vexists", handlerValueExists)
	mux.HandleFunc("bstats", handlerBStats)
	mux.HandleFunc("backup", handlerBackup)
	mux.HandleFunc("restore", handlerRestore)

	mux.HandleFunc("publish", handlerPublish)
	mux.HandleFunc("subscribe", handlerSubscribe)
	mux.HandleFunc("psubscribe", handlerPsubscribe)

	return mux
}
