package main

import (
	"fmt"
	"strconv"

	"github.com/uretgec/mydb-server/mystorage/loggers"
	"github.com/uretgec/mydb-server/mystorage/services"

	"github.com/tidwall/redcon"
)

func checkCmdRequired(args [][]byte, required int) bool {
	// args 0 dan başladığı için her zaman required verisine bir ekliyoruz.
	required += 1
	return len(args) >= required
}

func connAccept(conn redcon.Conn) bool {
	loggers.Sugar.With("Client", conn.RemoteAddr()).Info("New client connection accepted")
	return true
}

func connClose(conn redcon.Conn, err error) {
	loggers.Sugar.With("clientip", conn.RemoteAddr()).Info("client connection closed")
	if err != nil {
		loggers.Sugar.With("error", err).Error("redcon conn closed")
	}
}

// CMD: detach
func handlerDetach(conn redcon.Conn, cmd redcon.Command) {
	hconn := conn.Detach()
	loggers.Sugar.Info("connection has been detached")

	go func() {
		defer hconn.Close()
		hconn.WriteString("OK")
		hconn.Flush()
	}()
}

// CMD: quit
func handlerQuit(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("OK")
	conn.Close()
}

// CMD: ping
func handlerPing(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("pong")
}

// CMD: set bucketName id "json_stringfy_data"
func handlerSet(conn redcon.Conn, cmd redcon.Command) {
	if *dbReadOnly {
		conn.WriteError("storage read only mode active")
		return
	}

	if !checkCmdRequired(cmd.Args, 3) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	id := cmd.Args[2]
	val := cmd.Args[3]

	uid, err := services.Store.Set(bucketName, id, val)
	if err != nil {
		loggers.Sugar.With("error", err).Error("set command error")
		conn.WriteError(fmt.Sprintf("Set command error: %s", err.Error()))
		return
	}

	// Always return uid or id (id has to be sometimes nil)
	conn.WriteAny(uid)
}

// CMD: get bucketName id
func handlerGet(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	id := cmd.Args[2]

	result, err := services.Store.Get(bucketName, id)
	if err != nil {
		loggers.Sugar.With("error", err).Error("get command error")
		conn.WriteError(fmt.Sprintf("Get command error: %s", err.Error()))
		return
	}

	conn.WriteAny(result)
}

// CMD: mget bucketName id1 id2 id3 ....
func handlerMGet(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	var ids [][]byte
	for i := 2; i < len(cmd.Args); i++ {
		ids = append(ids, cmd.Args[i])
	}

	result, err := services.Store.MGet(bucketName, ids...)
	if err != nil {
		loggers.Sugar.With("error", err).Error("mget command error")
		conn.WriteError(fmt.Sprintf("MGet command error: %s", err.Error()))
		return
	}

	conn.WriteAny(result)
}

// CMD: list bucketName cursor perpage
func handlerList(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 3) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	cursor := cmd.Args[2]
	perpage, _ := strconv.Atoi(string(cmd.Args[3]))

	result, err := services.Store.List(bucketName, cursor, perpage)
	if err != nil {
		loggers.Sugar.With("error", err).Error("list command error")
		conn.WriteError(fmt.Sprintf("Get command error: %s", err.Error()))
		return
	}

	if result == nil {
		conn.WriteAny([]string{})
		return
	}

	conn.WriteAny(result)
}

// CMD: prevlist bucketName cursor perpage
func handlerPrevList(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 3) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	cursor := cmd.Args[2]
	perpage, _ := strconv.Atoi(string(cmd.Args[3]))

	result, err := services.Store.PrevList(bucketName, cursor, perpage)
	if err != nil {
		loggers.Sugar.With("error", err).Error("list command error")
		conn.WriteError(fmt.Sprintf("Get command error: %s", err.Error()))
		return
	}

	if result == nil {
		conn.WriteAny([]string{})
		return
	}

	conn.WriteAny(result)
}

// CMD: exists bucketName key
func handlerExists(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	id := cmd.Args[2]

	found, err := services.Store.Exist(bucketName, id)
	if err != nil {
		loggers.Sugar.With("error", err).Error("exists command error")
		conn.WriteError(fmt.Sprintf("exists command error: %s", err.Error()))
		return
	}

	conn.WriteAny(found)
}

// CMD: vexists bucketName value
func handlerValueExists(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	id := cmd.Args[2]

	found, err := services.Store.ValueExist(bucketName, id)
	if err != nil {
		loggers.Sugar.With("error", err).Error("vexists command error")
		conn.WriteError(fmt.Sprintf("vexists command error: %s", err.Error()))
		return
	}

	conn.WriteAny(found)
}

// CMD: del bucketName id
func handlerDel(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]
	id := cmd.Args[2]

	err := services.Store.Del(bucketName, id)
	if err != nil {
		loggers.Sugar.With("error", err).Error("del command error")
		conn.WriteError(fmt.Sprintf("Del command error: %s", err.Error()))
		return
	}

	conn.WriteInt64(1)
}

// CMD: bstats bucketName
func handlerBStats(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 1) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	bucketName := cmd.Args[1]

	conn.WriteInt(services.Store.BStats(bucketName))
}

// CMD: backup path/ filename
func handlerBackup(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	path := string(cmd.Args[1])
	filename := string(cmd.Args[2])

	err := services.Store.Backup(path, filename)
	if err != nil {
		loggers.Sugar.With("error", err).Error("backup command error")
		conn.WriteError(fmt.Sprintf("Backup command error: %s", err.Error()))
		return
	}

	conn.WriteString("OK")
}

// CMD: restore path/ filename
func handlerRestore(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	path := string(cmd.Args[1])
	filename := string(cmd.Args[2])

	err := services.Store.Restore(path, filename)
	if err != nil {
		loggers.Sugar.With("error", err).Error("restore command error")
		conn.WriteError(fmt.Sprintf("Restore command error: %s", err.Error()))
		return
	}

	conn.WriteString("OK")
}

// CMD: publish channel message
func handlerPublish(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 2) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	channelName := string(cmd.Args[1])
	msg := string(cmd.Args[2])

	ok := ps.Publish(channelName, msg)

	conn.WriteInt(ok)
}

// CMD: subscribe channel
func handlerSubscribe(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 1) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	channelName := string(cmd.Args[1])

	ps.Subscribe(conn, channelName)
}

// CMD: psubscribe channel
func handlerPsubscribe(conn redcon.Conn, cmd redcon.Command) {
	if !checkCmdRequired(cmd.Args, 1) {
		loggers.Sugar.With("args", string(cmd.Raw)).Error("not enough required count")
		conn.WriteError("not enough required count")
		return
	}

	channelName := string(cmd.Args[1])

	ps.Psubscribe(conn, channelName)
}
