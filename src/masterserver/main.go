package main

import (
    "share/logger"
    "masterserver/def"
    "masterserver/rpc"
    "github.com/jmoiron/sqlx"
    _"github.com/go-sql-driver/mysql"
)

var log = logger.Instance()

var g_ServerConfig  = def.ServerConfig
var g_RPCHandler    = def.RPCHandler
var g_LoginDatabase = def.LoginDatabase

func main() {
    log.Info("MasterServer init")

    // read config
    g_ServerConfig.Read()

    // register events
    RegisterEvents()

    // init RPC handler
    g_RPCHandler.Init()
    g_RPCHandler.Port = g_ServerConfig.Port

    // register RPC packets
    rpc.RegisterPackets()

    // connect to login database
    log.Info("Attempting to connect to the Login database...")
    if db, err := sqlx.Connect("mysql", g_ServerConfig.LoginDB()); err != nil {
        log.Fatalf("[DATABASE] %s", err.Error())
    } else {
        log.Info("Successfully connected to the Login database!")
        *g_LoginDatabase = *db

        var version []string
        db.Select(&version, "SELECT VERSION()")
        log.Debugf("[DATABASE] Version: %s", version[0])
    }

    // start RPC Server
    g_RPCHandler.Run()
}