package go_test_redis

import (
	"reflect"
	"testing"
)

func TestParseInfoResponse(t *testing.T) {
	in := `
# Persistence
loading:1
rdb_changes_since_last_save:0
rdb_bgsave_in_progress:0
rdb_last_save_time:1603112806
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:-1
rdb_current_bgsave_time_sec:-1
rdb_last_cow_size:0
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok
aof_last_cow_size:0
module_fork_in_progress:0
module_fork_last_cow_size:0
loading_start_time:1603112806
loading_total_bytes:40133580
loading_loaded_bytes:2094393
loading_loaded_perc:5.22
loading_eta_seconds:1
`
	resp := parseInfoResponse(in)
	if !reflect.DeepEqual(resp, map[string]string{
		"loading":                      "1",
		"rdb_changes_since_last_save":  "0",
		"rdb_bgsave_in_progress":       "0",
		"rdb_last_save_time":           "1603112806",
		"rdb_last_bgsave_status":       "ok",
		"rdb_last_bgsave_time_sec":     "-1",
		"rdb_current_bgsave_time_sec":  "-1",
		"rdb_last_cow_size":            "0",
		"aof_enabled":                  "0",
		"aof_rewrite_in_progress":      "0",
		"aof_rewrite_scheduled":        "0",
		"aof_last_rewrite_time_sec":    "-1",
		"aof_current_rewrite_time_sec": "-1",
		"aof_last_bgrewrite_status":    "ok",
		"aof_last_write_status":        "ok",
		"aof_last_cow_size":            "0",
		"module_fork_in_progress":      "0",
		"module_fork_last_cow_size":    "0",
		"loading_start_time":           "1603112806",
		"loading_total_bytes":          "40133580",
		"loading_loaded_bytes":         "2094393",
		"loading_loaded_perc":          "5.22",
		"loading_eta_seconds":          "1",
	}) {
		t.Fatal(resp)
	}

	in = `
# Server
redis_version:6.0.8
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:25b38681eed52ae
redis_mode:standalone
os:Darwin 19.6.0 x86_64
arch_bits:64
multiplexing_api:kqueue
atomicvar_api:atomic-builtin
gcc_version:4.2.1
process_id:31040
run_id:4d0ef4a90616edf4331047ed988fd4bd9fe73bee
tcp_port:6379
uptime_in_seconds:62459
uptime_in_days:0
hz:10
configured_hz:10
lru_clock:9339745
executable:/usr/local/opt/redis/bin/redis-server
config_file:/usr/local/etc/redis.conf
io_threads_active:0

# Clients
connected_clients:1
client_recent_max_input_buffer:2
client_recent_max_output_buffer:0
blocked_clients:0
tracking_clients:0
clients_in_timeout_table:0

# Memory
used_memory:61450272
used_memory_human:58.60M
used_memory_rss:2383872
used_memory_rss_human:2.27M
used_memory_peak:61450272
used_memory_peak_human:58.60M
used_memory_peak_perc:100.00%
used_memory_overhead:1018874
used_memory_startup:1001472
used_memory_dataset:60431398
used_memory_dataset_perc:99.97%
allocator_allocated:61403664
allocator_active:2345984
allocator_resident:2345984
total_system_memory:17179869184
total_system_memory_human:16.00G
used_memory_lua:37888
used_memory_lua_human:37.00K
used_memory_scripts:0
used_memory_scripts_human:0B
number_of_cached_scripts:0
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:0.04
allocator_frag_bytes:18446744073650493936
allocator_rss_ratio:1.00
allocator_rss_bytes:0
rss_overhead_ratio:1.02
rss_overhead_bytes:37888
mem_fragmentation_ratio:0.04
mem_fragmentation_bytes:-59019792
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:16986
mem_aof_buffer:0
mem_allocator:libc
active_defrag_running:0
lazyfree_pending_objects:0

# Persistence
loading:0
rdb_changes_since_last_save:0
rdb_bgsave_in_progress:0
rdb_last_save_time:1603112806
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:-1
rdb_current_bgsave_time_sec:-1
rdb_last_cow_size:0
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok
aof_last_cow_size:0
module_fork_in_progress:0
module_fork_last_cow_size:0

# Stats
total_connections_received:5
total_commands_processed:4
instantaneous_ops_per_sec:0
total_net_input_bytes:509
total_net_output_bytes:19516
instantaneous_input_kbps:0.01
instantaneous_output_kbps:6.56
rejected_connections:0
sync_full:0
sync_partial_ok:0
sync_partial_err:0
expired_keys:0
expired_stale_perc:0.00
expired_time_cap_reached_count:0
expire_cycle_cpu_milliseconds:756
evicted_keys:0
keyspace_hits:0
keyspace_misses:0
pubsub_channels:0
pubsub_patterns:0
latest_fork_usec:0
migrate_cached_sockets:0
slave_expires_tracked_keys:0
active_defrag_hits:0
active_defrag_misses:0
active_defrag_key_hits:0
active_defrag_key_misses:0
tracking_total_keys:0
tracking_total_items:0
tracking_total_prefixes:0
unexpected_error_replies:0
total_reads_processed:13
total_writes_processed:8
io_threaded_reads_processed:0
io_threaded_writes_processed:0

# Replication
role:master
connected_slaves:0
master_replid:605b7aacc1c67bcf3620c3fc8f018b35ce4a31f8
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

# CPU
used_cpu_sys:21.235879
used_cpu_user:14.505765
used_cpu_sys_children:0.000000
used_cpu_user_children:0.000000

# Modules

# Cluster
cluster_enabled:0

# Keyspace
db0:keys=8,expires=0,avg_ttl=0
`
	resp = parseInfoResponse(in)
	if !reflect.DeepEqual(resp, map[string]string{
		"redis_version":                   "6.0.8",
		"redis_git_sha1":                  "00000000",
		"redis_git_dirty":                 "0",
		"redis_build_id":                  "25b38681eed52ae",
		"redis_mode":                      "standalone",
		"os":                              "Darwin 19.6.0 x86_64",
		"arch_bits":                       "64",
		"multiplexing_api":                "kqueue",
		"atomicvar_api":                   "atomic-builtin",
		"gcc_version":                     "4.2.1",
		"process_id":                      "31040",
		"run_id":                          "4d0ef4a90616edf4331047ed988fd4bd9fe73bee",
		"tcp_port":                        "6379",
		"uptime_in_seconds":               "62459",
		"uptime_in_days":                  "0",
		"hz":                              "10",
		"configured_hz":                   "10",
		"lru_clock":                       "9339745",
		"executable":                      "/usr/local/opt/redis/bin/redis-server",
		"config_file":                     "/usr/local/etc/redis.conf",
		"io_threads_active":               "0",
		"connected_clients":               "1",
		"client_recent_max_input_buffer":  "2",
		"client_recent_max_output_buffer": "0",
		"blocked_clients":                 "0",
		"tracking_clients":                "0",
		"clients_in_timeout_table":        "0",
		"used_memory":                     "61450272",
		"used_memory_human":               "58.60M",
		"used_memory_rss":                 "2383872",
		"used_memory_rss_human":           "2.27M",
		"used_memory_peak":                "61450272",
		"used_memory_peak_human":          "58.60M",
		"used_memory_peak_perc":           "100.00%",
		"used_memory_overhead":            "1018874",
		"used_memory_startup":             "1001472",
		"used_memory_dataset":             "60431398",
		"used_memory_dataset_perc":        "99.97%",
		"allocator_allocated":             "61403664",
		"allocator_active":                "2345984",
		"allocator_resident":              "2345984",
		"total_system_memory":             "17179869184",
		"total_system_memory_human":       "16.00G",
		"used_memory_lua":                 "37888",
		"used_memory_lua_human":           "37.00K",
		"used_memory_scripts":             "0",
		"used_memory_scripts_human":       "0B",
		"number_of_cached_scripts":        "0",
		"maxmemory":                       "0",
		"maxmemory_human":                 "0B",
		"maxmemory_policy":                "noeviction",
		"allocator_frag_ratio":            "0.04",
		"allocator_frag_bytes":            "18446744073650493936",
		"allocator_rss_ratio":             "1.00",
		"allocator_rss_bytes":             "0",
		"rss_overhead_ratio":              "1.02",
		"rss_overhead_bytes":              "37888",
		"mem_fragmentation_ratio":         "0.04",
		"mem_fragmentation_bytes":         "-59019792",
		"mem_not_counted_for_evict":       "0",
		"mem_replication_backlog":         "0",
		"mem_clients_slaves":              "0",
		"mem_clients_normal":              "16986",
		"mem_aof_buffer":                  "0",
		"mem_allocator":                   "libc",
		"active_defrag_running":           "0",
		"lazyfree_pending_objects":        "0",
		"loading":                         "0",
		"rdb_changes_since_last_save":     "0",
		"rdb_bgsave_in_progress":          "0",
		"rdb_last_save_time":              "1603112806",
		"rdb_last_bgsave_status":          "ok",
		"rdb_last_bgsave_time_sec":        "-1",
		"rdb_current_bgsave_time_sec":     "-1",
		"rdb_last_cow_size":               "0",
		"aof_enabled":                     "0",
		"aof_rewrite_in_progress":         "0",
		"aof_rewrite_scheduled":           "0",
		"aof_last_rewrite_time_sec":       "-1",
		"aof_current_rewrite_time_sec":    "-1",
		"aof_last_bgrewrite_status":       "ok",
		"aof_last_write_status":           "ok",
		"aof_last_cow_size":               "0",
		"module_fork_in_progress":         "0",
		"module_fork_last_cow_size":       "0",
		"total_connections_received":      "5",
		"total_commands_processed":        "4",
		"instantaneous_ops_per_sec":       "0",
		"total_net_input_bytes":           "509",
		"total_net_output_bytes":          "19516",
		"instantaneous_input_kbps":        "0.01",
		"instantaneous_output_kbps":       "6.56",
		"rejected_connections":            "0",
		"sync_full":                       "0",
		"sync_partial_ok":                 "0",
		"sync_partial_err":                "0",
		"expired_keys":                    "0",
		"expired_stale_perc":              "0.00",
		"expired_time_cap_reached_count":  "0",
		"expire_cycle_cpu_milliseconds":   "756",
		"evicted_keys":                    "0",
		"keyspace_hits":                   "0",
		"keyspace_misses":                 "0",
		"pubsub_channels":                 "0",
		"pubsub_patterns":                 "0",
		"latest_fork_usec":                "0",
		"migrate_cached_sockets":          "0",
		"slave_expires_tracked_keys":      "0",
		"active_defrag_hits":              "0",
		"active_defrag_misses":            "0",
		"active_defrag_key_hits":          "0",
		"active_defrag_key_misses":        "0",
		"tracking_total_keys":             "0",
		"tracking_total_items":            "0",
		"tracking_total_prefixes":         "0",
		"unexpected_error_replies":        "0",
		"total_reads_processed":           "13",
		"total_writes_processed":          "8",
		"io_threaded_reads_processed":     "0",
		"io_threaded_writes_processed":    "0",
		"role":                            "master",
		"connected_slaves":                "0",
		"master_replid":                   "605b7aacc1c67bcf3620c3fc8f018b35ce4a31f8",
		"master_replid2":                  "0000000000000000000000000000000000000000",
		"master_repl_offset":              "0",
		"second_repl_offset":              "-1",
		"repl_backlog_active":             "0",
		"repl_backlog_size":               "1048576",
		"repl_backlog_first_byte_offset":  "0",
		"repl_backlog_histlen":            "0",
		"used_cpu_sys":                    "21.235879",
		"used_cpu_user":                   "14.505765",
		"used_cpu_sys_children":           "0.000000",
		"used_cpu_user_children":          "0.000000",
		"cluster_enabled":                 "0",
		"db0":                             "keys=8,expires=0,avg_ttl=0",
	}) {
		t.Fatal(resp)
	}
}
