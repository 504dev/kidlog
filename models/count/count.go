package count

import (
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/types"
	"time"
)

func Find(dashId int, logname string, hostname string) (types.Counts, error) {
	duration := Logger.Time("/logs:time", time.Millisecond)
	where := `dash_id = ? and logname = ? and timestamp > now() - interval 7 day`
	values := []interface{}{dashId, logname}
	if hostname != "" {
		where += ` and hostname = ?`
		values = append(values, hostname)
	}
	sql := `
      select
        toStartOfMinute(timestamp),
        hostname,
        keyname,
        sum(inc),
        max(max),
        min(min),
        sum(avg_sum),
        sum(avg_num),
        sum(per_tkn),
        sum(per_ttl)
      from counts
      where ` + where + `
      group by
        timestamp, hostname, keyname
      order by
        timestamp desc, hostname, keyname
    `
	fmt.Println(sql)
	rows, err := clickhouse.Conn().Query(sql, values...)
	if err != nil {
		return nil, err
	}

	counts := types.Counts{}
	for rows.Next() {
		var timestamp time.Time
		var hostname, keyname string
		var inc, max, min, avgSum, perTotal, perTaken *float64
		var avgNum *int
		err = rows.Scan(&timestamp, &hostname, &keyname, &inc, &max, &min, &avgSum, &avgNum, &perTaken, &perTotal)
		if err != nil {
			return nil, err
		}
		metrics := types.Metrics{}
		if inc != nil {
			metrics.Inc = &types.Inc{Val: *inc}
		}
		if max != nil {
			metrics.Max = &types.Max{Val: *max}
		}
		if min != nil {
			metrics.Min = &types.Min{Val: *min}
		}
		if avgNum != nil {
			metrics.Avg = &types.Avg{Sum: *avgSum, Num: *avgNum}
		}
		if perTotal != nil {
			metrics.Per = &types.Per{Total: *perTotal, Taken: *perTaken}
		}
		counts = append(counts, &types.Count{
			Timestamp: timestamp.Unix(),
			Hostname:  hostname,
			Keyname:   keyname,
			Metrics:   metrics,
		})
	}
	duration()
	Logger.Inc("/logs:cnt", 1)
	return counts, nil
}