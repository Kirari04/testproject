local http = require("socket.http")
local ltn12 = require("ltn12")

core.register_action("verify_request", { "http-req" }, function(txn)
    local response_body = {}
    local res, code, headers, status = http.request {
        url = "http://127.0.0.1:8082/",
        sink = ltn12.sink.table(response_body)
    }
    if code == 200 then
        txn.set_var(txn, "txn.request_verified", true)
    else
        txn.set_var(txn, "txn.request_verified", false)
    end

    -- txn:set_status(200)
    -- txn:add_header("content-length", string.len(table.concat(response_body)))
    -- txn:add_header("content-type", "text/html")
    -- txn:start_response()
    -- txn:send(table.concat(response_body))
end)
