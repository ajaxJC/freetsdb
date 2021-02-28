package text

import (
	"bufio"
	"io"
	"strconv"

	"github.com/freetsdb/freetsdb/cmd/freets_tools/internal/format"
	"github.com/freetsdb/freetsdb/services/influxql"
	"github.com/freetsdb/freetsdb/models"
	"github.com/freetsdb/freetsdb/pkg/escape"
	"github.com/freetsdb/freetsdb/tsdb"
)

type Writer struct {
	w   *bufio.Writer
	key []byte
	err error
	m   Mode
}

type Mode bool

const (
	Series Mode = false
	Values Mode = true
)

func NewWriter(w io.Writer, mode Mode) *Writer {
	var wr *bufio.Writer
	if wr, _ = w.(*bufio.Writer); wr == nil {
		wr = bufio.NewWriter(w)
	}
	return &Writer{
		w:   wr,
		key: make([]byte, 1024),
		m:   mode,
	}
}

func (w *Writer) NewBucket(start, end int64) (format.BucketWriter, error) {
	return w, nil
}

func (w *Writer) Close() error { return w.w.Flush() }
func (w *Writer) Err() error   { return w.err }

func (w *Writer) BeginSeries(name, field []byte, typ influxql.DataType, tags models.Tags) {
	if w.err != nil {
		return
	}

	if w.m == Series {
		w.key = models.AppendMakeKey(w.key[:0], name, tags)
		w.key = append(w.key, ' ')
		w.key = append(w.key, escape.Bytes(field)...)
		w.w.Write(w.key)
		w.w.WriteByte('\n')
	}
}

func (w *Writer) EndSeries() {}

func (w *Writer) WriteIntegerCursor(cur tsdb.IntegerArrayCursor) {
	if w.err != nil || w.m == Series {
		return
	}

	buf := w.key
	for {
		a := cur.Next()
		if a.Len() == 0 {
			break
		}
		for i := range a.Timestamps {
			buf = buf[:0]

			buf = strconv.AppendInt(buf, a.Values[i], 10)
			buf = append(buf, 'i')
			buf = append(buf, ' ')
			buf = strconv.AppendInt(buf, a.Timestamps[i], 10)
			buf = append(buf, '\n')
			if _, w.err = w.w.Write(buf); w.err != nil {
				return
			}
		}
	}
}

func (w *Writer) WriteFloatCursor(cur tsdb.FloatArrayCursor) {
	if w.err != nil || w.m == Series {
		return
	}

	buf := w.key
	for {
		a := cur.Next()
		if a.Len() == 0 {
			break
		}
		for i := range a.Timestamps {
			buf = buf[:0]

			buf = strconv.AppendFloat(buf, a.Values[i], 'g', -1, 64)
			buf = append(buf, ' ')
			buf = strconv.AppendInt(buf, a.Timestamps[i], 10)
			buf = append(buf, '\n')
			if _, w.err = w.w.Write(buf); w.err != nil {
				return
			}
		}
	}
}

func (w *Writer) WriteUnsignedCursor(cur tsdb.UnsignedArrayCursor) {
	if w.err != nil || w.m == Series {
		return
	}

	buf := w.key
	for {
		a := cur.Next()
		if a.Len() == 0 {
			break
		}
		for i := range a.Timestamps {
			buf = buf[:0]

			buf = strconv.AppendUint(buf, a.Values[i], 10)
			buf = append(buf, 'u')
			buf = append(buf, ' ')
			buf = strconv.AppendInt(buf, a.Timestamps[i], 10)
			buf = append(buf, '\n')
			if _, w.err = w.w.Write(buf); w.err != nil {
				return
			}
		}
	}
}

func (w *Writer) WriteBooleanCursor(cur tsdb.BooleanArrayCursor) {
	if w.err != nil || w.m == Series {
		return
	}

	buf := w.key
	for {
		a := cur.Next()
		if a.Len() == 0 {
			break
		}
		for i := range a.Timestamps {
			buf = buf[:0]

			buf = strconv.AppendBool(buf, a.Values[i])
			buf = append(buf, ' ')
			buf = strconv.AppendInt(buf, a.Timestamps[i], 10)
			buf = append(buf, '\n')
			if _, w.err = w.w.Write(buf); w.err != nil {
				return
			}
		}
	}
}

func (w *Writer) WriteStringCursor(cur tsdb.StringArrayCursor) {
	if w.err != nil || w.m == Series {
		return
	}

	buf := w.key
	for {
		a := cur.Next()
		if a.Len() == 0 {
			break
		}
		for i := range a.Timestamps {
			buf = buf[:0]

			buf = append(buf, '"')
			buf = append(buf, models.EscapeStringField(a.Values[i])...)
			buf = append(buf, '"')
			buf = append(buf, ' ')
			buf = strconv.AppendInt(buf, a.Timestamps[i], 10)
			buf = append(buf, '\n')
			if _, w.err = w.w.Write(buf); w.err != nil {
				return
			}
		}
	}
}
