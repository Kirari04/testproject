# Configuring a Dynamic Per-Stream Bandwidth Limit 
Souce: [https://www.haproxy.com/blog/introduction-to-traffic-shaping-using-haproxy#configuring-a-dynamic-per-stream-bandwidth-limit](https://www.haproxy.com/blog/introduction-to-traffic-shaping-using-haproxy#configuring-a-dynamic-per-stream-bandwidth-limit)
```
frontend primeflix
   filter bwlim-out video-streaming default-limit 320k default-period 1s # 720p by default
   
   # Detect the resolution by matching on the request path
   http-request set-var(txn.resolution) int(360) if { path_beg /360p }
   http-request set-var(txn.resolution) int(480) if { path_beg /480p }
   http-request set-var(txn.resolution) int(720) if { path_beg /720p }
   http-request set-var(txn.resolution) int(1080) if { path_beg /1080p }
   http-request set-var(txn.resolution) int(4000) if { path_beg /4k }
 
   acl is_mp4   res.hdr(content-type) -m beg video/mp4
   acl is_360p  var(txn.resolution) -m int 360
   acl is_480p  var(txn.resolution) -m int 480
   acl is_720p  var(txn.resolution) -m int 720
   acl is_1080p var(txn.resolution) -m int 1080
   acl is_4k	var(txn.resolution) -m int 4000
 
   http-response allow if !is_mp4                                      	  # Only set a bandwidth limit for mp4 video
   http-response set-bandwidth-limit video-streaming                    	  # Set the default limit   720p => 320KB/s
   http-response set-bandwidth-limit video-streaming limit  90K if is_360p  # Override default limit  320p => 90KB/s
   http-response set-bandwidth-limit video-streaming limit 140K if is_480p  # Override default limit  480p => 140K/s
   http-response set-bandwidth-limit video-streaming limit 625K if is_1080p # Override default limit 1080p => 625KB/s
   http-response set-bandwidth-limit video-streaming limit 2500K if is_4k	  # Override default limit   4K  => 2.5MB/s
```