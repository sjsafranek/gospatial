package gospatial

import "github.com/sjsafranek/DiffDB/diff_store"
import "github.com/sjsafranek/DiffDB/diff_db"

var diffDb diff_db.DiffDb

func init() {
	diffDb = diff_db.NewDiffDb("skeleton.db")
}

func update_timeseries_datasource(datasource_id string, value []byte) {

	update_value := string(value)
	var ddata diff_store.DiffStore
	data, err := diffDb.Load(datasource_id)
	if nil != err {
		if err.Error() == "Not found" {
			// create new diffstore if key not found in database
			ddata = diff_store.NewDiffStore(datasource_id)
		} else {
			panic(err)
		}
	} else {
		ddata.Decode(data)
	}

	// update diffstore
	ddata.Update(update_value)

	// save to database
	enc, err := ddata.Encode()
	if nil != err {
		panic(err)
	}

	diffDb.Save(ddata.Name, enc)
}
