:: Run all tests - excluding vendor
for /f "" %%G in ('go list ./... ^| find /i /v "/vendor/"') do @go test %%G