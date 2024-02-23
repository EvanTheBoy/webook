local key = KEYS[1]
local cntKey = key..":cnt"
local val = ARGV[1]
local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
    --系统错误
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    --符合预期
    return 0
else
    --发送太频繁
    return -1
end