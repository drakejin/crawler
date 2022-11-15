package elastic

import (
	"time"

	"github.com/olivere/elastic/v7"
)

func New(dsn string, initTimeout, maxTimeout time.Duration) (*elastic.Client, error) {
	es, err := elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetGzip(true),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(initTimeout, maxTimeout))),
	)
	if err != nil {
		return nil, err
	}

	return es, nil
}
