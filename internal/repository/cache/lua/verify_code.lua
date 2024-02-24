local key = KEYS[1]

local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 出现异常
    return -1
elseif expectedCode == code then
    -- 输入对了
    redis.call("set", cntKey, -1)
    return 0
else
    -- 输错了
    redis.call("decr", cntKey, -1)
    return -2
end