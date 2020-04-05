// Script that has all the SQL Instructions
package Uni
/*
	Startup: The code executed when the database connects
	
*/

var unisql map[string]map[string]string = map[string]map[string]string{
	"sqlite3": map[string]string{
		"Startup": `
CREATE TABLE IF NOT EXISTS ServerData(id text not null, adminrole text default null);
CREATE TABLE IF NOT EXISTS NSFWList(userid text not null);
CREATE TABLE IF NOT EXISTS Modules(gID text not null, cID text not null, modules bigint not null default 0);
CREATE TABLE IF NOT EXISTS DerpiFilters(gID text not null, cID text not null, filterID bigint not null);
CREATE TABLE IF NOT EXISTS UniBucks(userID text not null, value real not null);
CREATE TABLE IF NOT EXISTS DailyUniBucks(userID text not null, nanoseconds bigint not null);
CREATE TABLE IF NOT EXISTS UserStocks(userid text not null, name text not null, quantity real not null);
CREATE TABLE StockValues(name text not null, value real not null);
PRAGMA journal_mode = WAL;
PRAGMA temp_store = 2;
PRAGMA synchronous = EXTRA;
`,
	"CheckGuild": "SELECT id FROM ServerData WHERE id IS '%s';",
	"CreateGuild": "INSERT INTO ServerData VALUES ('%s', null);",
	"GetGuildAdminRole": "SELECT adminrole FROM ServerData WHERE id IS '%s';",
	"CheckChannelModules": "SELECT cID FROM Modules WHERE cID IS '%s';",
	"CreateChannelModules": "INSERT INTO Modules VALUES ('%s', '%s', 0);",
	"GetChannelModules": "SELECT modules FROM Modules WHERE cID IS '%s';",
	"UpdateChannelModules": "UPDATE Modules SET modules = %d WHERE cID IS '%s'",
	"CheckNSFW": "SELECT userid FROM NSFWList WHERE userid IS '%s'",
	"GiveNSFW": "INSERT INTO NSFWList VALUES ('%s');",
	"RevokeNSFW": "DELETE FROM NSFWList WHERE userid IS '%s';",
	"GetDerpiFilter": "SELECT filterID FROM DerpiFilters WHERE cID IS '%s';",
	"InsertDerpiFilter": "INSERT INTO DerpiFilters VALUES ('%s', '%s', %v);",
	"UpdateDerpiFilter": "UPDATE DerpiFilters SET filterID = %v WHERE cID IS '%s';",
	},
	"postgres": map[string]string{
		"Startup": "", // TODO
	},
}

var availabledrivers []string = []string{"sqlite3"}