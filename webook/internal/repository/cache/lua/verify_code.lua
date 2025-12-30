local key = KEYS[1]
-- 用户输入的 code
local expectedCode = ARGV[1]
local code = redis.call("get",key)
local cntKey = key..":cnt"
-- 如果验证码不存在
if code == false then
    return -3
end
-- 从 Redis 获取验证次数
local cnt = redis.call("get",cntKey)
-- 转成一个数字
if cnt == false or tonumber(cnt) <= 0 then
    -- 说明 用户一直输错或验证次数已用完
    return -1
elseif expectedCode == code then
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户输错了，减少验证次数
    redis.call("decr",cntKey)
    return -2
end
