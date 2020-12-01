package goBase

import "errors"

// MysqlConfig mysql数据库的配置
type MysqlConfig struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DbName       string
	PreTable     string
	Charset      string
	MaxOpenCons  int
	MaxIdleConns int
}

// RedisConfig Redis的配置信息
type RedisConfig struct {
	Host string
	Port int
	Auth string
	DB   int
}

// Config 配置文件
type AppConfig struct {
	Port       int                //端口
	APIRouters map[string]*Action //
	MysqlConfs map[string]*MysqlConfig
	RedisConfs map[string]*RedisConfig
	KV         map[string]interface{}
}

// getDBCOnf 获取数据库默认配置
func (a *AppConfig) getDBCOnf() (*MysqlConfig, error) {
	return a.getDBConfByName("default")
}

// getDBConfByName   根据名字获取数据库的配置信息
func (a *AppConfig) getDBConfByName(name string) (*MysqlConfig, error) {
	val, ok := a.MysqlConfs[name]
	if ok {
		return val, nil
	}
	return &MysqlConfig{}, errors.New("配置不存在")
}

// getRedisConf 获取Redis的默认配置
func (a *AppConfig) getRedisConf() (*RedisConfig, error) {
	return a.getRedisConfByName("default")
}

// getRedisConfByName 根据名字查询redis的配置信息
func (a *AppConfig) getRedisConfByName(name string) (*RedisConfig, error) {
	val, ok := a.RedisConfs[name]
	if ok {
		return val, nil
	}
	return &RedisConfig{}, errors.New("配置不存在")
}
