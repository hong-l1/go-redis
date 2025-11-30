package dict

import "math/rand"

func (d *MyDict) Get(key string) (val interface{}, exists bool) {
	shard := d.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	hash := getHash(key)
	if val, found := shard.ht[0].lookup(key, hash); found {
		return val, true
	}
	if shard.rehashIdx != -1 {
		if val, found := shard.ht[1].lookup(key, hash); found {
			return val, true
		}
	}
	return nil, false
}

func (d *MyDict) Len() int {
	var total uint64
	for _, shard := range d.shards {
		shard.RLock()
		if shard.rehashIdx == -1 {
			total += shard.ht[0].used
		} else {
			total += shard.ht[0].used + shard.ht[1].used
		}
		shard.RUnlock()
	}
	return int(total)
}

func (d *MyDict) Put(key string, val interface{}) (result int) {
	shard := d.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	if shard.rehashIdx != -1 {
		shard.rehashStep()
	}
	hash := getHash(key)
	if shard.ht[0].update(key, val, hash) {
		return 0
	}
	if shard.rehashIdx != -1 {
		if shard.ht[1].update(key, val, hash) {
			return 0
		}
	}
	targetHt := shard.ht[0]
	if shard.rehashIdx != -1 {
		targetHt = shard.ht[1]
	}
	idx := hash & targetHt.mask
	newEntry := &Entry{
		Key:  key,
		Val:  val,
		Next: targetHt.buckets[idx],
	}
	targetHt.buckets[idx] = newEntry
	targetHt.used++
	if shard.rehashIdx == -1 {
		shard.expand()
	}
	return 1
}

func (d *MyDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, exist := d.Get(key)
	if !exist {
		return d.Put(key, val)
	}
	return 0
}

func (d *MyDict) PutIfExists(key string, val interface{}) (result int) {
	_, exist := d.Get(key)
	if exist {
		return d.Put(key, val)
	}
	return 0
}

func (d *MyDict) Remove(key string) (interface{}, int) {
	shard := d.getShard(key)
	shard.RUnlock()
	defer shard.RUnlock()
	if shard.rehashIdx != -1 {
		shard.rehashStep()
	}
	hash := getHash(key)
	if val, ok := shard.ht[0].del(key, hash); ok {
		// 即使在 ht[0] 删除了，也要注意负载检查，Redis 有缩容逻辑，这里暂略
		return val, 1
	}
	if shard.rehashIdx != -1 {
		if val, ok := shard.ht[1].del(key, hash); ok {
			return val, 1
		}
	}
	return nil, 0
}

func (d *MyDict) ForEach(consumer Consumer) {
	for _, shard := range d.shards {
		shard.RLock()
		ht0 := shard.ht[0]
		for _, head := range ht0.buckets {
			for e := head; e != nil; e = e.Next {
				consumer(e.Key, e.Val)
			}
		}
		if shard.rehashIdx != -1 {
			ht1 := shard.ht[1]
			for _, head := range ht1.buckets {
				for e := head; e != nil; e = e.Next {
					consumer(e.Key, e.Val)
				}
			}
		}
		shard.RUnlock()
	}
}

func (d *MyDict) Keys() []string {
	keys := d.Len()
	arr := make([]string, 0, keys)
	for _, shard := range d.shards {
		shard.RLock()
		ht0 := shard.ht[0]
		for _, head := range ht0.buckets {
			for e := head; e != nil; e = e.Next {
				arr = append(arr, e.Key)
			}
		}
		if shard.rehashIdx != -1 {
			ht1 := shard.ht[1]
			for _, head := range ht1.buckets {
				for e := head; e != nil; e = e.Next {
					arr = append(arr, e.Key)
				}
			}
		}
		shard.RUnlock()
	}
	return arr
}

func (d *MyDict) RandomKeys(limit int) []string {
	result := make([]string, 0, limit)
	maxAttempts := limit * 2
	attempts := 0
	for len(result) < limit && attempts < maxAttempts {
		attempts++
		s := d.pickWeightedShard()
		if s == nil {
			break
		}
		s.RLock()
		key, found := s.getRandomKeyFromShard()
		s.RUnlock()
		if found {
			result = append(result, key)
		}
	}
	return result
}

func (d *MyDict) RandomDistinctKeys(limit int) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, limit)
	maxAttempts := limit * 2
	attempts := 0
	for len(result) < limit && attempts < maxAttempts {
		attempts++
		s := d.pickWeightedShard()
		s.RLock()
		key, found := s.getRandomKeyFromShard()
		s.RUnlock()
		if found {
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				result = append(result, key)
			}
		}
	}
	return result
}

func (d *MyDict) Clear() {
	for k, v := range d.shards {
		v.Lock()
		d.shards[k] = NewShard()
		v.Unlock()
	}
}
func (d *MyDict) pickWeightedShard() *Shard {
	counts := make([]int64, len(d.shards))
	var total int64 = 0
	for i, s := range d.shards {
		n := int64(s.ht[0].used + s.ht[1].used)
		counts[i] = n
		total += n
	}
	if total == 0 {
		return nil
	}
	r := rand.Int63n(total)
	var accum int64 = 0
	for i, count := range counts {
		accum += count
		if r < accum {
			return d.shards[i]
		}
	}
	for i, count := range counts {
		if count > 0 {
			return d.shards[i]
		}
	}
	return nil
}
