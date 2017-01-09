package check

type error_log_format struct {
	Line int
	File string
}

func Err(err error) {
	if err != nil {
		//_, file, line, ok := runtime.Caller(1)
		panic(err)
		/*
			trace.Info("service_lua_doRequest", service_lua_log_format{
				Url:    url,
				Params: data,
				Res:    addStr,
			})
		*/
	}
}
