package export

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"freelanceledger/model"
	"io"
	"strconv"
)

func JSON(w io.Writer, records []model.Record) error { return json.NewEncoder(w).Encode(records) }
func CSV(w io.Writer, records []model.Record) error {
	c := csv.NewWriter(w)
	if e := c.Write([]string{"id", "profile", "kind", "amount", "currency", "status", "occurred_at"}); e != nil {
		return e
	}
	for _, r := range records {
		if e := c.Write([]string{r.ID, r.ProfileID, r.Kind, strconv.FormatInt(r.Amount, 10), r.Currency, r.Status, r.OccurredAt.Format("2006-01-02T15:04:05Z")}); e != nil {
			return e
		}
	}
	c.Flush()
	return c.Error()
}
func ReadCSV(r io.Reader) ([]model.Record, error) {
	cr := csv.NewReader(bufio.NewReader(r))
	if _, e := cr.Read(); e != nil {
		return nil, e
	}
	out := []model.Record{}
	for {
		row, e := cr.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		if len(row) < 7 {
			return nil, fmt.Errorf("short row")
		}
		amt, e := strconv.ParseInt(row[3], 10, 64)
		if e != nil {
			return nil, e
		}
		out = append(out, model.Record{ID: row[0], ProfileID: row[1], Kind: row[2], Amount: amt, Currency: row[4], Status: row[5]})
	}
	return out, nil
}
func JSONLines(w io.Writer, records []model.Record) error {
	enc := json.NewEncoder(w)
	for _, r := range records {
		if e := enc.Encode(r); e != nil {
			return e
		}
	}
	return nil
}
func MarshalRecord(r model.Record) string { b, _ := json.Marshal(r); return string(b) }
