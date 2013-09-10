local ttl = function (keys, min, max)
	local s = 0
	for i = 1, table.getn(keys), 1 do
		s = s + redis.call("ZREMRANGEBYSCORE", keys[i], min, max)
	end
	return s
end
return ttl(redis.call("KEYS", KEYS[1]), ARGV[1], ARGV[2])