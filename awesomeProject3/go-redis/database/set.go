package database

import (
	"awesomeProject3/go-redis/datastruct/set"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
	"strconv"
)

func (d *DB) getAsSet(key string) (set.Set, resp.Reply) {
	entity, exists := d.GetEntity(key)
	if !exists {
		return nil, nil
	}
	s, ok := entity.Data.(set.Set)
	if !ok {
		return nil, reply.NewErrReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return s, nil
}
func execSAdd(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	members := args[1:]
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		s = set.NewMySet()
		db.PutEntity(key, &database.DataEntity{Data: s})
	}
	result := 0
	for _, v := range members {
		result += s.Add(string(v))
	}
	return reply.NewIntReply(result)
}
func execSIsMember(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	member := string(args[1])
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		return reply.NewIntReply(0)
	}
	if s.Has(member) {
		return reply.NewIntReply(1)
	}
	return reply.NewIntReply(0)
}
func execSMembers(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		return reply.NewEmptyMultiBulkReply()
	}
	members := s.Members()
	result := make([][]byte, len(members))
	for i, v := range members {
		result[i] = []byte(v)
	}
	return reply.NewMultiBulkReply(result)
}
func execSRem(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	members := args[1:]
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		return reply.NewIntReply(0)
	}
	result := 0
	for _, v := range members {
		result += s.Remove(string(v))
	}
	if s.Len() == 0 {
		db.Remove(key)
	}
	return reply.NewIntReply(result)
}
func execSCard(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		return reply.NewIntReply(0)
	}
	return reply.NewIntReply(s.Len())
}
func execSPop(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	count := 1
	if len(args) > 1 {
		var err error
		count, err = strconv.Atoi(string(args[1]))
		if err != nil {
			return reply.NewErrReply("ERR value is not an integer or out of range")
		}
	}
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		if len(args) > 1 {
			return reply.NewEmptyMultiBulkReply()
		}
		return reply.NewNilBulkReply()
	}
	members := s.RandomDistinctMembers(count)
	result := make([][]byte, len(members))
	for i, v := range members {
		s.Remove(v)
		result[i] = []byte(v)
	}
	if s.Len() == 0 {
		db.Remove(key)
	}
	if len(args) > 1 {
		return reply.NewMultiBulkReply(result)
	}
	if len(members) > 0 {
		return reply.NewBulkReply(result[0])
	}
	return reply.NewNilBulkReply()
}
func execSRandMember(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	count := 1
	if len(args) > 1 {
		var err error
		count, err = strconv.Atoi(string(args[1]))
		if err != nil {
			return reply.NewErrReply("ERR value is not an integer or out of range")
		}
	}
	s, errReply := db.getAsSet(key)
	if errReply != nil {
		return errReply
	}
	if s == nil {
		return reply.NewNilBulkReply()
	}
	var members []string
	if count > 0 {
		members = s.RandomDistinctMembers(count)
	} else if count < 0 {
		members = s.RandomMembers(-count)
	} else {
		return reply.NewEmptyMultiBulkReply()
	}
	result := make([][]byte, len(members))
	for i, v := range members {
		result[i] = []byte(v)
	}
	if len(args) > 1 {
		return reply.NewMultiBulkReply(result)
	}
	if len(members) > 0 {
		return reply.NewBulkReply(result[0])
	}
	return reply.NewNilBulkReply()
}

func init() {
	RegisterCommand("sadd", execSAdd, -3)
	RegisterCommand("sismember", execSIsMember, 3)
	RegisterCommand("smembers", execSMembers, 2)
	RegisterCommand("srem", execSRem, -3)
	RegisterCommand("scard", execSCard, 2)
	RegisterCommand("spop", execSPop, -2)
	RegisterCommand("srandmember", execSRandMember, -2)
}
