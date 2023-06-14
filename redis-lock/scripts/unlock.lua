local val = redis.call('get', KEYS[1])
if val == ARGV[1] then
    return redis.call('del', KEYS[1])
else
    return 0
end