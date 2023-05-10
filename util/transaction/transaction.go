package transaction

func Query[T any](txManager TransactionManager, queryFn func() (T, error)) (T, error) {
	var ret T
	err := txManager.Run(func() error {
		retI, err := queryFn()
		if err == nil {
			ret = retI
		}
		return err
	})
	return ret, err
}

type TransactionManager interface {
	Run(func() error) error
}
