package nfsbroker

import (
	"fmt"

	//"encoding/json"

	"database/sql"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type sqlStore struct {
	storeType string
	database  SqlConnection
}

func NewSqlStore(logger lager.Logger, dbDriver, username, password, host, port, dbName, caCert string) (Store, error) {

	var err error
	var toDatabase SqlVariant
	switch dbDriver {
	case "mysql":
		toDatabase = NewMySqlVariant(username, password, host, port, dbName, caCert)
	case "postgres":
		toDatabase = NewPostgresVariant(username, password, host, port, dbName, caCert)
	default:
		err = fmt.Errorf("Unrecognized Driver: %s", dbDriver)
		logger.Error("db-driver-unrecognized", err)
		return nil, err
	}
	return NewSqlStoreWithVariant(logger, toDatabase)
}

func NewSqlStoreWithVariant(logger lager.Logger, toDatabase SqlVariant) (Store, error) {
	database := NewSqlConnection(toDatabase)

	err := initialize(logger, database)

	if err != nil {
		logger.Error("sql-failed-to-initialize-database", err)
		return nil, err
	}

	return &sqlStore{
		database:  database,
	}, nil
}

func initialize(logger lager.Logger, db SqlConnection) error {
	logger = logger.Session("initialize-database")
	logger.Info("start")
	defer logger.Info("end")

	var err error
	err = db.Connect(logger)
	if err != nil {
		logger.Error("sql-failed-to-connect", err)
		return err
	}

	// TODO: uniquify table names?
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS service_instances(
				id VARCHAR(255) PRIMARY KEY,
				value VARCHAR(4096)
			)
		`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS service_bindings(
				id VARCHAR(255) PRIMARY KEY,
				value VARCHAR(4096)
			)
		`)
	return err
}

func (s *sqlStore) Restore(logger lager.Logger) error {
	//logger = logger.Session("restore-state")
	//logger.Info("start")
	//defer logger.Info("end")
	//
	//query := `SELECT id, value FROM service_instances`
	//rows, err := s.database.Query(query)
	//if err != nil {
	//	logger.Error("failed-query", err)
	//	return err
	//}
	//if rows != nil {
	//	for rows.Next() {
	//		var (
	//			id, value       string
	//			serviceInstance ServiceInstance
	//		)
	//
	//		err := rows.Scan(
	//			&id,
	//			&value,
	//		)
	//		if err != nil {
	//			logger.Error("failed-scanning", err)
	//			continue
	//		}
	//
	//		err = json.Unmarshal([]byte(value), &serviceInstance)
	//		if err != nil {
	//			logger.Error("failed-unmarshaling", err)
	//			continue
	//		}
	//		state.InstanceMap[id] = serviceInstance
	//	}
	//
	//	if rows.Err() != nil {
	//		logger.Error("failed-getting-next-row", rows.Err())
	//	}
	//}
	//
	//query = `SELECT id, value FROM service_bindings`
	//_, err = s.database.Query(query)
	//if err != nil {
	//	logger.Error("failed-query", err)
	//	return err
	//}
	//if rows != nil {
	//	for rows.Next() {
	//		var (
	//			id, value      string
	//			serviceBinding brokerapi.BindDetails
	//		)
	//
	//		err := rows.Scan(
	//			&id,
	//			&value,
	//		)
	//		if err != nil {
	//			logger.Error("failed-scanning", err)
	//			continue
	//		}
	//
	//		err = json.Unmarshal([]byte(value), &serviceBinding)
	//		if err != nil {
	//			logger.Error("failed-unmarshaling", err)
	//			continue
	//		}
	//		state.BindingMap[id] = serviceBinding
	//	}
	//
	//	if rows.Err() != nil {
	//		logger.Error("failed-getting-next-row", rows.Err())
	//	}
	//}

	return nil
}

func (s *sqlStore) Save(logger lager.Logger) error {
//	logger = logger.Session("save-state")
//	logger.Info("start", lager.Data{"instanceId": instanceId, "bindingId": bindingId})
//	defer logger.Info("end")
//
//	if instanceId != "" {
//		err, keyValueInTable := s.keyValueInTable(logger, "id", instanceId, "service_instances")
//		if err != nil {
//			return err
//		}
//		if keyValueInTable {
//			query := `DELETE FROM service_instances WHERE id=?`
//			_, err := s.database.Exec(query, instanceId)
//			if err != nil {
//				logger.Error("failed-exec", err)
//				return err
//			}
//			return nil
//		}
//		instance, _ := state.InstanceMap[instanceId]
//		logger.Info("instance-found", lager.Data{"instance": instance})
//		jsonValue, err := json.Marshal(&instance)
//		if err != nil {
//			logger.Error("failed-marshaling", err)
//			return err
//		}
//		query := `INSERT INTO service_instances (id, value) VALUES (?, ?)`
//
//		_, err = s.database.Exec(query, instanceId, jsonValue)
//		if err != nil {
//			logger.Error("failed-exec", err)
//			return err
//		}
//		state.InstanceMap = make(map[string]ServiceInstance)
//		logger.Info("Insert-Success", lager.Data{"Instance ID":instanceId, "JSON Value":jsonValue})
//		return nil
//
//	} else if bindingId != "" {
//		err, keyValueInTable := s.keyValueInTable(logger, "id", bindingId, "service_bindings")
//		if err != nil {
//			return err
//		}
//
//		if keyValueInTable {
//			query := `DELETE FROM service_bindings WHERE id=?`
//			_, err := s.database.Exec(query, bindingId)
//			if err != nil {
//				logger.Error("failed-exec", err)
//				return err
//			}
//			return nil
//
//		}
//		binding, _ := state.BindingMap[bindingId]
//		jsonValue, err := json.Marshal(&binding)
//		if err != nil {
//			logger.Error("failed-marshaling", err)
//			return err
//		}
//
//		query := `INSERT INTO service_bindings (id, value) VALUES (?, ?)`
//		_, err = s.database.Exec(query, bindingId, jsonValue)
//		if err != nil {
//			logger.Error("failed-exec", err)
//			return err
//		}
//		logger.Info("Insert-Success", lager.Data{"Binding ID":bindingId, "JSON Value":jsonValue})
//		return nil
//
//	}
//	err := fmt.Errorf("Both BindingID and InstanceID's were nil!")
//	logger.Error("failed-exec", err)
//	return err
return nil
}

func (s *sqlStore) Cleanup() error {
	return s.database.Close()
}

func (s *sqlStore) keyValueInTable(logger lager.Logger, key, value, table string) (error, bool) {
	var queriedServiceID string
	query := fmt.Sprintf(`SELECT %s.%s FROM %s WHERE %s.%s = ?`, table, key, table, table, key)
	row := s.database.QueryRow(query, value)
	if row == nil {
		err := fmt.Errorf("Row error!")
		logger.Error("failed-query", err)
		return err, true
	}
	err := row.Scan(&queriedServiceID)
	if err == nil {
		return nil, true
	} else if err == sql.ErrNoRows {
		return nil, false
	}

	logger.Debug("failed-query", lager.Data{"Query": query})
	logger.Error("failed-query", err)
	return err, true
}

func (s *sqlStore) RetrieveInstanceDetails(id string) (ServiceInstance, error) {
	return ServiceInstance{}, nil
}
func (s *sqlStore) RetrieveBindingDetails(id string) (brokerapi.BindDetails, error) {
	return brokerapi.BindDetails{}, nil
}
func (s *sqlStore) CreateInstanceDetails(id string, details ServiceInstance) error {
	return nil
}
func (s *sqlStore) CreateBindingDetails(id string, details brokerapi.BindDetails) error {
	return nil
}
func (s *sqlStore) DeleteInstanceDetails(id string) error {
	return nil
}
func (s *sqlStore) DeleteBindingDetails(id string) error {
	return nil
}