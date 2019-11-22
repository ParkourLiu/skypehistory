package main

import (
	"mtcomm/db/mysql"
)

//插入school
func iu_name(name string) error {
	sql := "INSERT INTO `s_user`(`u_name`) VALUES(?) ON DUPLICATE KEY UPDATE ut=NOW();"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{name}})
}

func iu_friends(name, friend_id, display_name, Type string) error {
	sql := "INSERT INTO `s_friends`(`u_name`,`friend_id`,`display_name`,`type`) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE display_name=?, `ut`=NOW();"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{name, friend_id, display_name, Type, display_name}})
}

func iu_chat(name, friend_id string, msg message) error {
	sql := "INSERT INTO `s_chat`(`u_name`,`friend_id`,`version`,`messagetype`,`content`,`type`,`from`)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `ut`=NOW();"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{name, friend_id, msg.Version, msg.Messagetype, msg.Content, msg.Type, msg.From}})
}

func s_name() ([]map[string]string, error) {
	sql := "SELECT * FROM `s_user`"
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

func s_friends(reqData ReqData) ([]map[string]string, error) {
	sql := "SELECT * FROM `s_friends` WHERE `u_name`=? "
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{reqData.Name}})
}

func s_chat(reqData ReqData) ([]map[string]string, error) {
	sql := "SELECT * FROM `s_chat` WHERE `u_name`=? AND friend_id=? AND `version`<?  ORDER BY `version` DESC LIMIT 20"
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{reqData.Name, reqData.Friend_id, reqData.Version}})
}
