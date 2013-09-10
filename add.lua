local add = function (key, value, timestamp)
	return redis.call("ZADD", key, timestamp, "[" .. tostring(value) .. "," .. tostring(timestamp) .. "]")
end
return add(KEYS[1], ARGV[1], ARGV[2])