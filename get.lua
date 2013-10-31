local get = function (keys, min, max)
        local data = {}
        for i = 1, table.getn(keys), 1 do
                local h = {}
                h["target"] = keys[i]
                h["datapoints"] = redis.call("ZRANGEBYSCORE", keys[i], min, max)
                data[i] = h
        end
        return cjson.encode(data)
end
return get(redis.call("KEYS", KEYS[1]), ARGV[1], ARGV[2])