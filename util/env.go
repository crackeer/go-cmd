package util

import "os"

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetHost() string {
	if value := GetEnv("GF_HOST"); len(value) > 0 {
		return value
	}
	return "127.0.0.1"
}

func GetMySQLHost() string {
	if value := GetEnv("MYSQL_HOST"); len(value) > 0 {
		return value
	}
	if value := GetHost(); len(value) > 0 {
		return value
	}
	return "127.0.0.1"
}

func GetMySQLPort() string {
	if value := GetEnv("MYSQL_PORT"); len(value) > 0 {
		return value
	}
	return "3306"
}

func GetMySQLPassword() string {
	if value := GetEnv("MYSQL_PASSWORD"); len(value) > 0 {
		return value
	}
	return "z_php_root"
}

func GetMySQLUser() string {
	if value := GetEnv("MYSQL_USER"); len(value) > 0 {
		return value
	}
	return "root"
}

func GetApolloUser() string {
	if value := GetEnv("GF_APOLLO_USER"); len(value) > 0 {
		return value
	}
	return "apollo"
}

func GetApolloPassword() string {
	if value := GetEnv("APOLLO_PASSWORD"); len(value) > 0 {
		return value
	}
	return "admin"
}

func GetApolloHost() string {
	if value := GetEnv("APOLLO_HOST"); len(value) > 0 {
		return value
	}

	if value := GetHost(); len(value) > 0 {
		return "http://" + value + ":8070"
	}
	return "http://127.0.0.1:8070"
}
