#!/usr/bin/env tarantool

local log_level = 5

if os.getenv("LOG_LEVEL") ~= nil then
    log_level = tonumber(os.getenv("LOG_LEVEL"));
end

box.cfg {
    username = nil;
    work_dir = nil;
    wal_dir = "/data";
    snap_dir = "/data";
    vinyl_dir = "/data";
    listen = 3301;
    pid_file = "/run/tarantool.pid";
    background = false;
    slab_alloc_arena = tonumber(os.getenv("MEMORY_LIMIT"));
    slab_alloc_minimal = 64;
    slab_alloc_maximal = 10485760;
    slab_alloc_factor = 1.06;
    snapshot_period = 0;
    snapshot_count = 6;
    panic_on_snap_error = true;
    panic_on_wal_error = true;
    rows_per_wal = 500000;
    snap_io_rate_limit = nil;
    wal_mode = "none";
    wal_dir_rescan_delay = 2.0;
    io_collect_interval = nil;
    readahead = 16320;
    log_level = log_level;
    logger_nonblock = true;
    too_long_threshold = 0.5;
}


local function bootstrap()
    box.schema.user.grant('guest', 'read,write,execute', 'universe')
end

local function dd(...)
    local x = debug.getinfo(2)
    local dbg = string.format('[%s:%d][%s]', x.source, x.currentline, x.name)
    require('log').info('\n ++DEBUG++ %s \n  %s', dbg, require('yaml').encode({ { ... } }))
end

box.once('once1', bootstrap)

package.cpath = '/usr/lib/tarantool/func.so;;'

box.schema.func.create('gex', { language = 'C' })
box.schema.user.grant('guest', 'execute', 'function', 'gex')

local name = 'cache'
local cache_space = name
local cache_data = 'c_' .. name

box.schema.space.create(cache_space, { if_not_exists = true })
box.space[cache_space]:create_index('primary', { type = 'hash', parts = { 1, 'string' }, if_not_exists = true })
box.space[cache_space]:create_index('cid', { type = 'TREE', unique = false, parts = { 2, 'unsigned' }, if_not_exists = true })

box.schema.space.create(cache_data, { if_not_exists = true })
box.space[cache_data]:create_index('primary', { type = 'hash', parts = { 1, 'unsigned' }, if_not_exists = true })

box.space[cache_space]:insert({'key1', 5})
box.space[cache_data]:insert({5, 'data1'})

local net_box = require('net.box')
local capi_connection = net_box:new(3301)
local res = capi_connection:call('gex', 'cache', 'key1')

dd(res)