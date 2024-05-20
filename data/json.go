package data

import (
	"encoding/json"
	"io"
	"os"
	"slices"
	"time"

	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/types"
	"golang.org/x/exp/maps"
)

type Decoder func(io.Reader) (map[types.ID]types.Actionable, error)

type Encoder func(map[types.ID]types.Actionable, io.Writer) error

type ReadGetter func() (io.ReadCloser, error)

type WriterGetter func() (io.WriteCloser, error)

type Comparer func(a, b types.Actionable) int

type jsonActivity struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Done        bool          `json:"done"`
}

const FileName = "database.json"

type JSON struct {
	reader  ReadGetter
	writer  WriterGetter
	decode  Decoder
	encode  Encoder
	compare Comparer
}

func NewJSONRepository(
	rcg ReadGetter,
	wcg WriterGetter,
	d Decoder,
	e Encoder,
	c Comparer,
) Repository {
	return &JSON{reader: rcg, writer: wcg, decode: d, encode: e, compare: c}
}

func FileJSONRepository() Repository {
	return NewJSONRepository(FileReader, FileWriter, JSONDecode, JSONEncode, Compare)
}

func (j *JSON) Get(id types.ID) (types.Actionable, error) {
	am, err := j.mapFromFile()
	if err != nil {
		return nil, err
	}
	a, ok := am[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (j *JSON) List() ([]types.Actionable, error) {
	am, err := j.mapFromFile()
	if err != nil {
		return nil, err
	}
	as := maps.Values(am)
	slices.SortFunc(as, j.compare)
	return as, nil
}

func (j *JSON) Put(a types.Actionable) error {
	am, err := j.mapFromFile()
	if err != nil {
		return err
	}
	am[a.Identify()] = a
	return j.mapToFile(am)
}

func (j *JSON) Remove(id types.ID) error {
	am, err := j.mapFromFile()
	if err != nil {
		return err
	}
	delete(am, id)
	return j.mapToFile(am)
}

func (j JSON) mapFromFile() (map[types.ID]types.Actionable, error) {
	rc, err := j.reader()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	am, err := j.decode(rc)
	if err != nil {
		if err == io.EOF {
			return map[types.ID]types.Actionable{}, nil
		}
		return nil, err
	}
	return am, nil
}

func (j JSON) mapToFile(am map[types.ID]types.Actionable) error {
	wc, err := j.writer()
	if err != nil {
		return err
	}
	defer wc.Close()
	return j.encode(am, wc)
}

func JSONDecode(r io.Reader) (map[types.ID]types.Actionable, error) {
	var concreteMap map[types.ID]jsonActivity
	err := json.NewDecoder(r).Decode(&concreteMap)
	abstractMap := apply.ToValues(concreteMap, func(concrete jsonActivity) types.Actionable {
		abstract := types.NewActivity(types.ID(concrete.ID), concrete.Description)
		abstract.Work(concrete.Duration)
		if concrete.Done {
			abstract.Do()
		} else {
			abstract.Undo()
		}
		return abstract
	})
	return abstractMap, err
}

func JSONEncode(abstractMap map[types.ID]types.Actionable, w io.Writer) error {
	concreteMap := apply.ToValues(abstractMap, func(abstract types.Actionable) jsonActivity {
		return jsonActivity{
			ID:          string(abstract.Identify()),
			Description: abstract.Describe(),
			Duration:    abstract.Worked(),
			Done:        abstract.Done(),
		}
	})
	return json.NewEncoder(w).Encode(concreteMap)
}

func FileReader() (io.ReadCloser, error) {
	r, err := os.Open(FileName)
	if err != nil {
		return os.Create(FileName)
	}
	return r, nil
}

func FileWriter() (io.WriteCloser, error) {
	return os.Create(FileName)
}

func Compare(a, b types.Actionable) int {
	idA, idB := a.Identify(), b.Identify()
	if idA > idB {
		return 1
	}
	if idA < idB {
		return -1
	}
	return 0
}
