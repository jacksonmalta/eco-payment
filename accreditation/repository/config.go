package repository

type Config struct {
	TableName string
}

func (c *Config) WithTableName(tableName string) *Config {
	c.TableName = tableName
	return c
}
