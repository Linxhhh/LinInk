-- 限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber( ARGV[2])
-- 窗口的起始时间
local now = tonumber(ARGV[3])
local min = now - window

-- 移除 score < min 的元素
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

-- 获取当前窗口内的请求数量
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
-- local cnt = redis.call('ZCOUNT', key, min, '+inf')

if cnt >= threshold then
    -- 执行限流
    return "true"
else
    -- 向 zset 中添加新元素
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window)
    return "false"
end